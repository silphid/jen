package persist

import (
	"strings"

	"github.com/Samasource/jen/src/internal/model"
	"github.com/Samasource/jen/src/internal/project"
)

// LoadConfig loads config object from project file
func LoadConfig(config *model.Config, projectDir string) error {
	proj, err := project.Load(projectDir)
	if err != nil {
		return err
	}

	if config.TemplateName == "" {
		config.TemplateName = proj.TemplateName
	}
	config.Values.Variables = proj.Variables

	initDefaultPlaceholders(config)
	return nil
}

// SaveConfig saves config object to project file
func SaveConfig(config *model.Config, projectDir string) error {
	proj := project.Project{
		TemplateName: config.TemplateName,
		Variables:    config.Values.Variables,
	}

	err := proj.Save(projectDir)
	if err != nil {
		return err
	}

	initDefaultPlaceholders(config)
	return nil
}

func initDefaultPlaceholders(config *model.Config) {
	config.Values.Placeholders = make(model.VarMap, 2)

	project, ok := config.Values.Variables["PROJECT"]
	if ok {
		config.Values.Placeholders["projekt"] = strings.ToLower(project)
		config.Values.Placeholders["PROJEKT"] = strings.ToUpper(project)
	}
}
