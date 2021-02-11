package render

import (
	"github.com/Samasource/jen/src/internal/evaluation"
	"github.com/Samasource/jen/src/internal/exec"
)

// Render represents an executable that renders a given source sub-folder
// of the current template's dir into the project's dir.
type Render struct {
	InputDir string
}

func (r Render) String() string {
	return "render"
}

// Execute renders a given source sub-folder of the current template's dir
// into the project's dir.
func (r Render) Execute(context exec.Context) error {
	return evaluation.Render(context, r.InputDir, context.GetProjectDir())
}
