package persist

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Samasource/jen/internal/constant"
	"github.com/Samasource/jen/internal/model"
)

// LoadConfig loads config object from jen file
func LoadConfig(config *model.Config) error {
	jenfile, err := LoadJenFileFromDir(config.ProjectDir)
	if err != nil {
		return err
	}

	if config.TemplateName == "" {
		config.TemplateName = jenfile.TemplateName
	}
	config.Values.Variables = jenfile.Variables

	parseSetVars(config)
	initDefaultPlaceholders(config)
	return nil
}

// SaveConfig saves config object to jen file
func SaveConfig(config *model.Config) error {
	jenfile := model.JenFile{
		Version:      constant.JenFileVersion,
		TemplateName: config.TemplateName,
		Variables:    config.Values.Variables,
	}

	err := SaveJenFileToDir(config.ProjectDir, jenfile)
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

var setVarRegexp = regexp.MustCompile(`^(\w+)=(.*)$`)

func parseSetVars(config *model.Config) error {
	config.SetVars = make(map[string]string, len(config.SetVarsRaw))
	for _, raw := range config.SetVarsRaw {
		submatch := setVarRegexp.FindStringSubmatch(raw)
		if submatch == nil {
			return fmt.Errorf("failed to parse set variable %q", raw)
		}
		config.SetVars[submatch[1]] = submatch[2]
		config.Values.Variables[submatch[1]] = submatch[2]
	}
	return nil
}
