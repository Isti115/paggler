package utils

import (
	"fmt"
	"strings"
)

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

func HighlightDiff(diff string) string {
	lines := strings.Split(diff, "\n")
	colored := make([]string, len(lines))

	for i, line := range lines {
		colored[i] = highlightLine(line)
	}

	return strings.Join(colored[:], "\n")
}
