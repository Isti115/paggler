package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func makePatch(name, content string) {
	f, _ := os.Create(fmt.Sprintf("./%s.patch", name))
	f.WriteString(content)
	f.Close()
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

func (m model) Init() bubbletea.Cmd {
	return nil
}

func (m model) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case bubbletea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, bubbletea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "p":
			choice := m.choices[m.cursor]
			makePatch(
				strings.Trim(strings.Split(choice, ":")[2], " "),
				getStash(m.cursor),
			)
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

	descStyle := lipgloss.NewStyle().Margin(2)
	return lipgloss.JoinHorizontal(lipgloss.Top, descStyle.Render(s), descStyle.Render(s2))
}

func main() {
	p := bubbletea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
