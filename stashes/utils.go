package stashes

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func getOutput(name string, arg ...string) string {
	out, err := exec.Command(name, arg...).Output()

	if err != nil {
		log.Fatal(err)
	}

	return string(out)
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
		"stash", "show", "-p", fmt.Sprintf("stash@{%d}", i),
	)
}

func makePatch(name, content string) {
	f, _ := os.Create(fmt.Sprintf("paggler/[_]-%s.patch", name))
	f.WriteString(content)
	f.Close()
}
