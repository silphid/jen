package persist

import (
	"fmt"
	. "github.com/Samasource/jen/internal/constant"
	"github.com/Samasource/jen/internal/model"
	"github.com/kylelemons/go-gypsy/yaml"
	"path"
)

func LoadJenFileFromDir(projectDir string) (*model.JenFile, error) {
	specFilePath := path.Join(projectDir, JenFileName)
	yamlFile, err := yaml.ReadFile(specFilePath)
	if err != nil {
		return nil, err
	}

	_map, ok := yamlFile.Root.(yaml.Map)
	if !ok {
		return nil, fmt.Errorf("jen file root is expected to be an object")
	}
	return loadJenFileFromMap(_map)
}

func loadJenFileFromMap(_map yaml.Map) (*model.JenFile, error) {
	jenfile := new(model.JenFile)

	// Load metadata
	metadata, err := getRequiredMap(_map, "metadata")
	if err != nil {
		return nil, err
	}
	jenfile.TemplateName, err = getRequiredStringFromMap(metadata, "template")
	if err != nil {
		return nil, err
	}
	jenfile.Version, err = getRequiredStringFromMap(metadata, "version")
	if err != nil {
		return nil, err
	}
	if jenfile.Version != JenFileVersion {
		return nil, fmt.Errorf("unsupported jenfile version %s (expected %s)", jenfile.Version, JenFileVersion)
	}

	// Load variables
	variables, err := getRequiredMap(_map, "variables")
	if err != nil {
		return nil, err
	}
	jenfile.Variables, err = loadVariables(variables)
	if err != nil {
		return nil, err
	}

	return jenfile, nil
}

func loadVariables(_map yaml.Map) (model.VarMap, error) {
	variables := make(model.VarMap, len(_map))
	for name, value := range _map {
		scalar, ok := value.(yaml.Scalar)
		if !ok {
			return nil, fmt.Errorf("value of variable %q must be a raw string", name)
		}
		variables[name] = scalar.String()
	}
	return variables, nil
}
