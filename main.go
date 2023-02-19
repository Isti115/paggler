package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	lip "github.com/charmbracelet/lipgloss"
)

type model struct {
	cursor   int
	choices  []string
	selected map[int]struct{}
}

func getOutput(name string, arg ...string) string {
	fout, ferr := exec.Command(name, arg...).Output()

	if ferr != nil {
		log.Fatal(ferr)
	}

	return string(fout)
}

func getLines(text string) []string {
	return strings.Split(strings.Trim(text, "\n"), "\n")
}

func getStashes() []string {
	return getLines(getOutput("git", "stash", "list"))
}

func getStash(i int) string {
	return getOutput(
		"git",
		"-c", "color.ui=always",
		"stash", "show", "-p", fmt.Sprintf("stash@{%d}", i),
	)
}

func initialModel() model {
	return model{
		// choices: []string{"Buy carrots", "Buy celery", "Buy kohlrabi"},
		choices: getStashes(),

		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	s := "What should we buy at the market?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	s += "\nPress q to quit.\n\n"

	s2 := getStash(m.cursor)

	descStyle := lip.NewStyle().Margin(2)
	return lip.JoinHorizontal(lip.Top, descStyle.Render(s), descStyle.Render(s2))
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
