package render

import (
	"path/filepath"

	"github.com/Samasource/jen/src/internal/evaluation"
	"github.com/Samasource/jen/src/internal/exec"
)

// Render represents an executable that renders a given source sub-folder
// of the current template's dir into the project's dir.
type Render struct {
	InputDir  string
	OutputDir string
}

func (r Render) String() string {
	return "render"
}

// Execute renders a given source sub-folder of the current template's dir
// into the project's dir.
func (r Render) Execute(context exec.Context) error {
	inputDir := filepath.Join(context.GetTemplateDir(), r.InputDir)
	outputDir := filepath.Join(context.GetProjectDir(), r.OutputDir)
	return evaluation.Render(context, inputDir, outputDir)
}
