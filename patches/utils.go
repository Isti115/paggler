package patches

import (
	"log"
	"os"
	"github.com/isti115/paggler/utils"
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

func getPatch(path string) string {
	patch, err := os.ReadFile(path)

	if err != nil {
		log.Fatal(err)
	}

	return utils.HighlightDiff(string(patch))
}
