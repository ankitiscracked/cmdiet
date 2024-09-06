package ui

import (
	"cmdiet/meals"
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	mealToUpdate     selectedMeal
	listStyle        = lipgloss.NewStyle().Margin(1, 2)
	mealEditingStyle = lipgloss.NewStyle()
)

var (
	windowSizeMsg  tea.WindowSizeMsg
	formInputIndex = 0
)

type selectedMeal struct {
	MealName       string
	MealComponents []mealComponent
}

type mealListItem struct {
	title, desc string
	mealId      int
}

func (i mealListItem) Title() string {
	return i.title
}

func (i mealListItem) Description() string {
	return i.desc
}

func (i mealListItem) FilterValue() string {
	return i.title
}

type mealsModel struct {
	list       list.Model
	editing    bool
	detailForm *huh.Form
	quitting   bool
}

type mealComponent struct {
	ComponentId   int
	ComponentType meals.MealComponentType
	Amount        string
}

type errorMsg struct {
	err error
}
type refreshMealListMsg struct{}

func (m mealsModel) Init() tea.Cmd {
	return nil
}

func (m mealsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowSizeMsg = msg
		x, y := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-x, msg.Height-y)
	case refreshMealListMsg:
		items, err := mealListItems()
		if err != nil {
			log.Fatal(err)
		}
		m.list.SetItems(items)
	case tea.KeyMsg:
		if !m.editing {
			switch msg.String() {
			case "ctrl+c", "esc", "q":
				m.quitting = true
				return m, tea.Quit
			case "enter":
				selectedItem := m.list.SelectedItem().(mealListItem)
				m.detailForm = createDetailForm(selectedItem)
				m.editing = true
				return m, m.detailForm.Init()
			}
		} else {
			switch msg.Type {
			case tea.KeyEsc:
				m.detailForm = nil
				m.editing = false
				return m, nil
			case tea.KeyDown:
				if formInputIndex != 4 {
					formInputIndex++
					return m, m.detailForm.NextField()
				}
			case tea.KeyUp:
				if formInputIndex != 0 {
					formInputIndex--
					return m, m.detailForm.PrevField()
				}
			}
		}
	}

	if !m.editing {
		m.list, cmd = m.list.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.detailForm.UpdateFieldPositions()
		form, cmd := m.detailForm.Update(msg)
		if f, ok := form.(*huh.Form); ok {
			m.detailForm = f
			cmds = append(cmds, cmd)
		}

		if m.detailForm.State == huh.StateCompleted {
			cmds = append(cmds, updateMealCmd(m))
			m.detailForm = nil
			m.editing = false
		}
	}

	return m, tea.Batch(cmds...)
}

func updateMealCmd(m mealsModel) tea.Cmd {
	return func() tea.Msg {
		selectedItem := m.list.SelectedItem().(mealListItem)
		payload := meals.MealPayload{
			Name:       mealToUpdate.MealName,
			Components: make([]meals.MealComponent, 0),
		}

		for _, component := range mealToUpdate.MealComponents {
			componentPayload := meals.MealComponent{
				Type: component.ComponentType,
				Id:   component.ComponentId,
			}

			if component.ComponentType == meals.FixedMealComponentType {
				componentPayload.Count = atoiIgnoreError(component.Amount)
			} else {
				componentPayload.AmountInGrams = atoiIgnoreError(component.Amount)
			}

			payload.Components = append(payload.Components, componentPayload)
		}

		if _, err := meals.MS.UpdateMeal(selectedItem.mealId, meals.UpdateMealPayload(payload)); err != nil {
			return errorMsg{err}
		}
		return refreshMealListMsg{}
	}
}

func atoiIgnoreError(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func createDetailForm(selectedMeal mealListItem) *huh.Form {
	meal, err := meals.MS.GetMeal(selectedMeal.mealId)
	if err != nil {
		log.Fatal(err)
	}

	var formGruops []*huh.Group
	formGruops = append(
		formGruops, huh.NewGroup(
			huh.NewInput().Title("Meal name").Value(&mealToUpdate.MealName),
		),
	)
	for _, fixedComponent := range meal.FixedComponentData {
		mealToUpdate.MealComponents = append(mealToUpdate.MealComponents, mealComponent{ComponentId: int(fixedComponent.ID), Amount: strconv.Itoa(fixedComponent.Amount)})
		formGruops = append(formGruops, huh.NewGroup(
			huh.NewInput().
				Title("Count").
				Value(&mealToUpdate.MealComponents[len(mealToUpdate.MealComponents)-1].Amount),
		).Title(fixedComponent.Component.Name))
	}

	for _, variableComponent := range meal.VariableComponentData {
		mealToUpdate.MealComponents = append(mealToUpdate.MealComponents, mealComponent{ComponentId: int(variableComponent.ID), Amount: strconv.Itoa(variableComponent.AmountInGrams)})
		formGruops = append(formGruops, huh.NewGroup(
			huh.NewInput().
				Title("Amount in grams").
				Value(&mealToUpdate.MealComponents[len(mealToUpdate.MealComponents)-1].Amount),
		).Title(variableComponent.Component.Name))
	}

	form := huh.NewForm(formGruops...)
	return form
}

func (m mealsModel) View() string {
	if m.quitting {
		return ""
	}
	if m.editing {
		return lipgloss.JoinHorizontal(lipgloss.Center, listStyle.Render(m.list.View()), m.detailForm.View())
	} else {
		return listStyle.Render(m.list.View())
	}
}

func mealListItems() ([]list.Item, error) {
	allMeals, err := meals.MS.GetAllMeals()
	if err != nil {
		return nil, err
	}

	var items []list.Item
	for _, meal := range allMeals {
		items = append(items, mealListItem{
			title:  meal.Name,
			desc:   fmt.Sprintf("Calories: %s | Protein: %s | Carbs: %s | Fats: %s", strconv.Itoa(meal.GetTotalCalories()), strconv.Itoa(meal.GetMealMacro(meals.Protein)), strconv.Itoa(meal.GetMealMacro(meals.Carbs)), strconv.Itoa(meal.GetMealMacro(meals.Fat))),
			mealId: int(meal.ID),
		})
	}
	return items, nil
}

func ListMealsModel() mealsModel {
	items, err := mealListItems()
	if err != nil {
		log.Fatal(err)
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Your meals"

	model := mealsModel{list: list}
	return model
}
