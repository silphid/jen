package render

import (
	"path"

	"github.com/Samasource/jen/src/internal/evaluation"
	"github.com/Samasource/jen/src/internal/model"
	"github.com/Samasource/jen/src/internal/project"
)

type Render struct {
	Source string
}

func (r Render) String() string {
	return "render"
}

func (r Render) Execute(config *model.Config) error {
	projectDir, err := project.GetDir()
	if err != nil {
		return err
	}

	inputDir := path.Join(config.TemplateDir, r.Source)
	return evaluation.Render(config.Values, inputDir, projectDir)
}
