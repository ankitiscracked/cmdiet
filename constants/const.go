package constants

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/go-playground/validator/v10"
)

var Program *tea.Program
var Validator *validator.Validate

var TableViewKeyMap = []key.Binding{
	key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("->/l", "next batch of diets"),
	),
	key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("<-/h", "previous batch of diets"),
	),
}

var DayViewKeyMap = []key.Binding{
	key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("->/l", "diet summary for the next day"),
	),
	key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("<-/h", "diet summary for the previous day"),
	),
}
