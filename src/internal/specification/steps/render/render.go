package render

import (
	"github.com/Samasource/jen/internal/evaluation"
	"github.com/Samasource/jen/internal/model"
	"path"
)

type Render struct {
	Source string
}

func (r Render) String() string {
	return "render"
}

func (r Render) Execute(config *model.Config) error {
	inputDir := path.Join(config.TemplateDir, r.Source)
	return evaluation.Render(config.Values, inputDir, config.ProjectDir)
}
