package persist

import (
	"strings"

	. "github.com/Samasource/jen/internal/constant"
	"github.com/Samasource/jen/internal/model"
)

func LoadJenFile(config *model.Config) error {
	jenfile, err := LoadJenFileFromDir(config.ProjectDir)
	if err != nil {
		return err
	}

	if config.TemplateName == "" {
		config.TemplateName = jenfile.TemplateName
	}
	config.Values.Variables = jenfile.Variables

	InitDefaultPlaceholders(config)
	return nil
}

func SaveJenFile(config *model.Config) error {
	jenfile := model.JenFile{
		Version:      JenFileVersion,
		TemplateName: config.TemplateName,
		Variables:    config.Values.Variables,
	}

	err := SaveJenFileToDir(config.ProjectDir, jenfile)
	if err != nil {
		return err
	}

	InitDefaultPlaceholders(config)
	return nil
}

func InitDefaultPlaceholders(config *model.Config) {
	config.Values.Placeholders = make(model.VarMap, 2)

	project, ok := config.Values.Variables["PROJECT"]
	if ok {
		config.Values.Placeholders["projekt"] = strings.ToLower(project)
		config.Values.Placeholders["PROJEKT"] = strings.ToUpper(project)
	}
}
