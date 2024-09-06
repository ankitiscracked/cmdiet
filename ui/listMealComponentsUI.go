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
	componentToUpdate selectedComponent
)

type selectedComponent struct {
	Type        meals.MealComponentType
	Name        string
	Calories    string
	Protein     string
	Carbs       string
	Fat         string
	ComponentID int
}

type componentListItem struct {
	title, desc   string
	componentID   int
	componentType meals.MealComponentType
}

func (i componentListItem) Title() string       { return i.title }
func (i componentListItem) Description() string { return i.desc }
func (i componentListItem) FilterValue() string { return i.title }

type componentsModel struct {
	list       list.Model
	editing    bool
	detailForm *huh.Form
	quitting   bool
}

type refreshComponentListMsg struct{}

func (m componentsModel) Init() tea.Cmd {
	return nil
}

func (m componentsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := listStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		if m.editing {
			switch msg.String() {
			case "esc":
				m.detailForm = nil
				m.editing = false
				return m, nil
			}
		} else {
			switch msg.String() {
			case "ctrl+c", "q":
				m.quitting = true
				return m, tea.Quit
			case "enter":
				selectedItem := m.list.SelectedItem().(componentListItem)
				m.detailForm = createMealComponentDetailForm(selectedItem)
				m.editing = true
				return m, m.detailForm.Init()
			}
		}

	case refreshComponentListMsg:
		var cmd tea.Cmd
		items, err := componentListItems()
		if err != nil {
			log.Fatal(err)
		}
		m.list.SetItems(items)
		m.list, cmd = m.list.Update(msg)
		return m, cmd
	}

	if m.editing {
		m.detailForm.UpdateFieldPositions()
		var cmd tea.Cmd
		form, cmd := m.detailForm.Update(msg)
		updatedForm, ok := form.(*huh.Form)
		if !ok {
			// Handle the error case if the type assertion fails
			return m, cmd
		}
		m.detailForm = updatedForm
		if m.detailForm.State == huh.StateCompleted {
			m.detailForm = nil
			m.editing = false
			return m, tea.Batch(cmd, func() tea.Msg {
				updateComponent(componentToUpdate)
				return refreshComponentListMsg{}
			})
		}
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m componentsModel) View() string {
	if m.quitting {
		return ""
	}
	if m.editing {
		return lipgloss.JoinHorizontal(lipgloss.Center, listStyle.Render(m.list.View()), m.detailForm.View())
	}
	return listStyle.Render(m.list.View())
}

func createMealComponentDetailForm(selectedItem componentListItem) *huh.Form {
	var err error
	if selectedItem.componentType == meals.FixedMealComponentType {
		var fixedComponent meals.FixedMealComponent
		fixedComponent, err = meals.MCS.GetFixedMealComponent(selectedItem.componentID)
		componentToUpdate = selectedComponent{
			ComponentID: selectedItem.componentID,
			Type:        selectedItem.componentType,
			Name:        fixedComponent.Name,
			Protein:     strconv.Itoa(fixedComponent.Protein),
			Carbs:       strconv.Itoa(fixedComponent.Carbs),
			Fat:         strconv.Itoa(fixedComponent.Fat),
		}
	} else {
		var variableComponent meals.VariableMealComponent
		variableComponent, err = meals.MCS.GetVariableMealComponent(selectedItem.componentID)
		componentToUpdate = selectedComponent{
			ComponentID: selectedItem.componentID,
			Type:        selectedItem.componentType,
			Name:        variableComponent.Name,
			Protein:     strconv.Itoa(variableComponent.ProteinPerHundredGram),
			Carbs:       strconv.Itoa(variableComponent.CarbsPerHundredGram),
			Fat:         strconv.Itoa(variableComponent.FatPerHundredGram),
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	return huh.NewForm(
		huh.NewGroup(
			huh.NewInput().Title("Name").Value(&componentToUpdate.Name),
			huh.NewInput().Title("Protein").Value(&componentToUpdate.Protein),
			huh.NewInput().Title("Carbs").Value(&componentToUpdate.Carbs),
			huh.NewInput().Title("Fat").Value(&componentToUpdate.Fat),
		).Title("Update meal component"),
	)
}

func updateComponent(component selectedComponent) {
	payload := meals.UpdateMealComponentPayload{
		Type:    component.Type,
		Name:    component.Name,
		Protein: atoi(component.Protein),
		Carbs:   atoi(component.Carbs),
		Fat:     atoi(component.Fat),
	}

	err := meals.MCS.UpdateMealComponent(component.ComponentID, payload)
	if err != nil {
		log.Fatal(err)
	}
}

func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func createComponentList() list.Model {
	items, err := componentListItems()
	if err != nil {
		log.Fatal(err)
	}
	componentList := list.New(items, list.NewDefaultDelegate(), 0, 0)
	componentList.Title = "Your meal components"
	return componentList
}

func ListMealsComponentsModel() componentsModel {
	componentList := createComponentList()
	return componentsModel{list: componentList}
}

func componentListItems() ([]list.Item, error) {
	allComponents, err := meals.MCS.GetAllMealComponents()
	if err != nil {
		return nil, err
	}

	var items []list.Item
	for _, component := range allComponents {
		items = append(items, componentListItem{
			title: component.Name,
			desc: fmt.Sprintf("Calories: %d | Protein: %d | Carbs: %d | Fat: %d",
				component.TotalCalories, component.TotalProtein, component.TotalCarbs, component.TotalFat),
			componentID:   int(component.Id),
			componentType: component.Type,
		})
	}
	return items, nil
}
