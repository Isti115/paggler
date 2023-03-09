package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"

	"github.com/isti115/paggler/patches"
	"github.com/isti115/paggler/stashes"
)

type Mode string

const (
	Patches Mode = "patches"
	Stashes Mode = "stashes"
)

type model struct {
	mode    Mode
	patches patches.Model
	stashes stashes.Model
}

func initialModel() model {
	return model{
		mode:    Patches,
		patches: patches.InitialModel(),
		stashes: stashes.InitialModel(),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.patches.Init(),
		m.stashes.Init(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "p", "P":
			m.mode = Patches
		case "s", "S":
			m.mode = Stashes
		default:
			switch m.mode {

			case Patches:
				pm, pc := m.patches.Update(msg)
				m.patches = pm
				return m, pc

			case Stashes:
				sm, sc := m.stashes.Update(msg)
				m.stashes = sm
				return m, sc

			}
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.mode {

	case Patches:
		return m.patches.View()

	case Stashes:
		return m.stashes.View()

	}

	return "Select a mode using `p` or `s`!"
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
