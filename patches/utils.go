package patches

import (
	"log"
	"os"
	"sort"

	"github.com/isti115/paggler/utils"
)

type byName []string

func (s byName) Len() int {
	return len(s)
}
func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byName) Less(i, j int) bool {
	return s[i][4:] < s[j][4:]
}

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

	sort.Sort(byName(patches))

	return patches
}

func getPatch(path string) string {
	patch, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	return utils.HighlightDiff(string(patch))
}
