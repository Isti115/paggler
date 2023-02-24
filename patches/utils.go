package patches

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func makeDir() {
	os.MkdirAll("paggler", 0755)
}

func getPatches() []string {
	dir, err := os.Open("paggler")

	if err != nil {
		log.Fatal(err)
	}

	patches, err := dir.Readdirnames(0)

	if err != nil {
		log.Fatal(err)
	}

	return patches
}

var highlightRules = []struct {
	prefix string
	format string
}{
	{"diff", "[1m%s[0m"},
	{"index", "[1m%s[0m"},
	{"@@", "[36m%s[0m"},
	{"---", "[1m[33m%s[0m"},
	{"+++", "[1m[33m%s[0m"},
	{"-", "[31m%s[0m"},
	{"+", "[32m%s[0m"},
}

func highlightLine(line string) string {
	for _, rule := range highlightRules {
		if strings.HasPrefix(line, rule.prefix) {
			return fmt.Sprintf(rule.format, line)
		}
	}

	return line
}

func highlightDiff(diff string) string {
	lines := strings.Split(diff, "\n")
	colored := make([]string, len(lines))

	for i, line := range lines {
		colored[i] = highlightLine(line)
	}

	return strings.Join(colored[:], "\n")
}

func getPatch(path string) string {
	patch, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	return highlightDiff(string(patch))
}
