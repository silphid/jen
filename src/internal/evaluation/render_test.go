package evaluation

import (
	"github.com/Samasource/jen/internal/specification"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestRender(t *testing.T) {
	values := specification.Values{
		Variables: map[string]interface{}{
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
			err := Render(values, path.Join("testdata", f.DataDir, "input"), outputDir)
			assert.NoError(t, err)
			compareDirsRecursively(t, path.Join("testdata", f.DataDir, "output"), outputDir)
		})
	}
}
