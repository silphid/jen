package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/Samasource/jen/src/internal/constant"
	"github.com/Samasource/jen/src/internal/helpers"
	"gopkg.in/yaml.v2"
)

// GetDir returns the project's root dir. It finds it by looking for the jen.yaml file
// in current working dir and then walking up the directory structure until it reaches the
// volume's root dir. If it doesn't find it, it returns an empty string.
func GetDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("finding project's root dir: %w", err)
	}

	for {
		path := filepath.Join(dir, constant.ProjectFileName)
		if helpers.PathExists(path) {
			return dir, nil
		}
		if dir == "/" {
			return "", nil
		}
		dir = filepath.Dir(dir)
	}
}

// Project represents the configuration file in a project's root dir
type Project struct {
	Version      string
	TemplateName string
	Variables    map[string]string
}

// Save saves project file into given project directory
func (project Project) Save(dir string) error {
	project.Version = constant.ProjectFileVersion
	doc, err := yaml.Marshal(project)
	if err != nil {
		return err
	}

	filePath := path.Join(dir, constant.ProjectFileName)
	return ioutil.WriteFile(filePath, doc, os.ModePerm)
}

// Load loads the project file from given project directory
func Load(dir string) (*Project, error) {
	specFilePath := filepath.Join(dir, constant.ProjectFileName)
	buf, err := ioutil.ReadFile(specFilePath)
	if err != nil {
		return nil, fmt.Errorf("loading project file: %w", err)
	}
	var project Project
	err = yaml.Unmarshal(buf, &project)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling project file yaml: %w", err)
	}

	if project.Version != constant.ProjectFileVersion {
		return nil, fmt.Errorf("unsupported jen project file version %s (expected %s)", project.Version, constant.ProjectFileVersion)
	}
	return &project, nil
}
