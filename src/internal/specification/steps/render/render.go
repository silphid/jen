package render

import (
	"github.com/Samasource/jen/internal/evaluation"
	"github.com/Samasource/jen/internal/specification"
)

type Render struct {
	Source string
}

func (r Render) String() string {
	return "render"
}

func (r Render) Execute(context specification.Context) error {
	return evaluation.Render(context.Values, context.InputDir, context.OutputDir)
}
