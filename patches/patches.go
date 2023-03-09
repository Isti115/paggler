package patches

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"golang.org/x/term"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	cursor  int
	message string
	patches []string
}

func InitialModel() Model {
	makeDir()

	return Model{
		patches: getPatches(),
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
			if m.cursor < len(m.patches)-1 {
				m.cursor++
			}

		case ">":
			exec.Command(
				"git", "apply",
				filepath.Join("paggler", m.patches[m.cursor]),
			).Run()
			m.message = "applied: " + filepath.Join("paggler", m.patches[m.cursor])

		case "<":
			exec.Command(
				"git", "apply", "--reverse",
				filepath.Join("paggler", m.patches[m.cursor]),
			).Run()
			m.message = "reversed: " + filepath.Join("paggler", m.patches[m.cursor])
		}
	}

	return m, nil
}

func (m Model) View() string {
	s := "Select a patch to toggle!\n\n"

	for i, choice := range m.patches {

		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if true {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += fmt.Sprintf("\n(Status: %s)\n", m.message)

	s += "\nPress q to quit.\n"

	s2 := getPatch(filepath.Join("paggler", m.patches[m.cursor]))

	physicalWidth, physicalHeight, _ := term.GetSize(int(os.Stdout.Fd()))
	descStyle := lipgloss.NewStyle().Margin(2)
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		descStyle.MaxWidth(physicalWidth/2).MaxHeight(physicalHeight).Render(s),
		descStyle.MaxWidth(physicalWidth/2).MaxHeight(physicalHeight).Render(s2),
	)
}
