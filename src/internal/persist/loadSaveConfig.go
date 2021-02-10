package persist

import (
	"strings"

	"github.com/Samasource/jen/src/internal/constant"
	"github.com/Samasource/jen/src/internal/model"
)

// LoadConfig loads config object from jen file
func LoadConfig(config *model.Config, projectDir string) error {
	jenfile, err := LoadJenFileFromDir(projectDir)
	if err != nil {
		return err
	}

	if config.TemplateName == "" {
		config.TemplateName = jenfile.TemplateName
	}
	config.Values.Variables = jenfile.Variables

	initDefaultPlaceholders(config)
	return nil
}

// SaveConfig saves config object to jen file
func SaveConfig(config *model.Config, projectDir string) error {
	jenfile := model.JenFile{
		Version:      constant.JenFileVersion,
		TemplateName: config.TemplateName,
		Variables:    config.Values.Variables,
	}

	err := SaveJenFileToDir(projectDir, jenfile)
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
