package persist

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/Samasource/jen/internal/constant"
	"github.com/Samasource/jen/internal/model"
	"gopkg.in/yaml.v2"
)

// SaveJenFileToDir saves jen file into given project directory
func SaveJenFileToDir(projectDir string, jenfile model.JenFile) error {
	_map := map[interface{}]interface{}{
		"metadata": map[interface{}]interface{}{
			"version":  constant.JenFileVersion,
			"template": jenfile.TemplateName,
		},
		"variables": jenfile.Variables,
	}

	doc, err := yaml.Marshal(&_map)
	if err != nil {
		return err
	}

	filePath := path.Join(projectDir, constant.JenFileName)
	return ioutil.WriteFile(filePath, doc, os.ModePerm)
}
