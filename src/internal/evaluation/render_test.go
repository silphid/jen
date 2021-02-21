package evaluation

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	context := context{
		vars: varMap{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  "true",
			"EMPTY_VAR": "",
		},
	}

	names := []string{
		"conditionals",
		"escaped-braces",
		"non-templated",
		"templated",
	}

	for _, name := range names {
		t.Run(name, func(t *testing.T) {
			outputDir := getTempDir()
			defer removeAll(outputDir)
			err := Render(context, path.Join("testdata", name, "input"), outputDir)
			assert.NoError(t, err)
			compareDirsRecursively(t, path.Join("testdata", name, "output"), outputDir)
		})
	}
}
