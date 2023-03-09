package stashes

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	cursor  int
	stashes []string
}

func InitialModel() Model {
	return Model{
		stashes: getStashes(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.stashes)-1 {
				m.cursor++
			}
		case "enter", " ":
			choice := m.stashes[m.cursor]
			makePatch(
				strings.Trim(strings.Split(choice, ":")[2], " "),
				getStash(m.cursor),
			)
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := "Select a stash to extract into a patch!\n\n"

	for i, choice := range m.stashes {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	s += "\nPress q to quit.\n\n"

	s2 := getStash(m.cursor)

	physicalWidth, physicalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	descStyle := lipgloss.NewStyle().Margin(2)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		descStyle.MaxWidth(physicalWidth/2).MaxHeight(physicalHeight).Render(s),
		descStyle.MaxWidth(physicalWidth/2).MaxHeight(physicalHeight).Render(s2),
	)
}
