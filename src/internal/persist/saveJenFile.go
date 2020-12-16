package persist

import (
	. "github.com/Samasource/jen/internal/constant"
	"github.com/Samasource/jen/internal/model"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

func SaveJenFileToDir(projectDir string, jenfile model.JenFile) error {
	_map := map[interface{}]interface{}{
		"metadata": map[interface{}]interface{}{
			"version":  JenFileVersion,
			"template": jenfile.TemplateName,
		},
		"variables": jenfile.Variables,
	}

	doc, err := yaml.Marshal(&_map)
	if err != nil {
		return err
	}

	filePath := path.Join(projectDir, JenFileName)
	return ioutil.WriteFile(filePath, doc, os.ModePerm)
}
