package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/Samasource/jen/src/internal/constant"
	"gopkg.in/yaml.v2"
)

// JenFile represents the configuration file in a project's root dir
type JenFile struct {
	Version      string
	TemplateName string
	Variables    map[string]string
}

// Save saves jen file into given project directory
func (jenFile JenFile) Save(dir string) error {
	jenFile.Version = constant.JenFileVersion
	doc, err := yaml.Marshal(jenFile)
	if err != nil {
		return err
	}

	filePath := path.Join(dir, constant.JenFileName)
	return ioutil.WriteFile(filePath, doc, os.ModePerm)
}

// Load loads the jen file from given project directory
func Load(dir string) (*JenFile, error) {
	specFilePath := filepath.Join(dir, constant.JenFileName)
	buf, err := ioutil.ReadFile(specFilePath)
	if err != nil {
		return nil, fmt.Errorf("loading jen file: %w", err)
	}
	var jenFile JenFile
	err = yaml.Unmarshal(buf, &jenFile)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling jen file yaml: %w", err)
	}

	if jenFile.Version != constant.JenFileVersion {
		return nil, fmt.Errorf("unsupported jenfile version %s (expected %s)", jenFile.Version, constant.JenFileVersion)
	}
	return &jenFile, nil
}
