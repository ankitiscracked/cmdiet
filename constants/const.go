package constants

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

var Program *tea.Program

type keyMap struct {
	NextBatch key.Binding
	PrevBatch key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.NextBatch,
		k.PrevBatch,
	}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextBatch, k.PrevBatch},
	}
}

var ViewTableKeyMap = keyMap{
	NextBatch: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("->/l", "next batch of diets"),
	),
	PrevBatch: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("<-/h", "previous batch of diets"),
	),
}
