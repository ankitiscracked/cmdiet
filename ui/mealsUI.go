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
	listStyle        = lipgloss.NewStyle().Margin(1, 2)
	mealEditingStyle = lipgloss.NewStyle()
)

var (
	mealToUpdate   mealFormValues
	windowSizeMsg  tea.WindowSizeMsg
	formInputIndex = 0
)

type listItem struct {
	title, desc string
	mealId      int
}

func (i listItem) Title() string {
	return i.title
}

func (i listItem) Description() string {
	return i.desc
}

func (i listItem) FilterValue() string {
	return i.title
}

type mealsModel struct {
	list       list.Model
	editing    bool
	detailForm *huh.Form
	quitting   bool
}

type mealFormValues struct {
	name     string
	calories string
	protein  string
	carbs    string
	fats     string
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
				selectedItem := m.list.SelectedItem().(listItem)
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
		selectedItem := m.list.SelectedItem().(listItem)
		meal := meals.Meal{
			Id:       selectedItem.mealId,
			Name:     mealToUpdate.name,
			Calories: atoiIgnoreError(mealToUpdate.calories),
			Protein:  atoiIgnoreError(mealToUpdate.protein),
			Carbs:    atoiIgnoreError(mealToUpdate.carbs),
			Fat:      atoiIgnoreError(mealToUpdate.fats),
		}
		if err := meals.MS.UpdateMeal(meal); err != nil {
			return errorMsg{err}
		}
		return refreshMealListMsg{}
	}
}

func atoiIgnoreError(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func createDetailForm(selectedMeal listItem) *huh.Form {
	meal, err := meals.MS.GetMeal(selectedMeal.mealId)
	if err != nil {
		log.Fatal(err)
	}
	mealToUpdate = mealFormValues{name: meal.Name, calories: strconv.Itoa(meal.Calories), protein: strconv.Itoa(meal.Protein), carbs: strconv.Itoa(meal.Carbs), fats: strconv.Itoa(meal.Fat)}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Meal name").Value(&mealToUpdate.name),
			huh.NewInput().Title("Calories").Value(&mealToUpdate.calories),
			huh.NewInput().Title("Protein").Value(&mealToUpdate.protein),
			huh.NewInput().Title("Carbs").Value(&mealToUpdate.carbs),
			huh.NewInput().Title("Fats").Value(&mealToUpdate.fats),
		).Title("Let's update your meal"),
	)
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
		items = append(items, listItem{
			title:  meal.Name,
			desc:   fmt.Sprintf("Calories: %s | Protein: %s | Carbs: %s | Fats: %s", strconv.Itoa(meal.Calories), strconv.Itoa(meal.Protein), strconv.Itoa(meal.Carbs), strconv.Itoa(meal.Fat)),
			mealId: meal.Id,
		})
	}
	return items, nil
}

func ViewMealsModel() mealsModel {
	items, err := mealListItems()
	if err != nil {
		log.Fatal(err)
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "Your meals"

	model := mealsModel{list: list}
	return model
}
