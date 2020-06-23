package internal

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
)

func Load(templateDir string) (Spec, error) {
	// Load file to buffer
	data, err := ioutil.ReadFile(path.Join(templateDir, "jen.yaml"))
	if err != nil {
		return Spec{}, err
	}

	// Parse buffer as yaml into map
	doc := Spec{
		TemplateDir: templateDir,
	}
	err = yaml.Unmarshal(data, &doc)
	if err != nil {
		return Spec{}, err
	}

	return doc, nil
}
