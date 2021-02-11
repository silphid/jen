package evaluation

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	context := context{
		vars: strMap{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  "true",
			"EMPTY_VAR": "",
		},
	}

	fixtures := []struct {
		Name    string
		DataDir string
	}{
		{
			Name:    "render1",
			DataDir: "render1",
		},
	}

	for _, f := range fixtures {
		t.Run(f.Name, func(t *testing.T) {
			outputDir := getTempDir()
			defer removeAll(outputDir)
			err := Render(context, path.Join("testdata", f.DataDir, "input"), outputDir)
			assert.NoError(t, err)
			compareDirsRecursively(t, path.Join("testdata", f.DataDir, "output"), outputDir)
		})
	}
}
