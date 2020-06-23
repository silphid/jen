package model

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Load() (Spec, error) {
	// Load file to buffer
	data, err := ioutil.ReadFile("examples/jen.yaml")
	if err != nil {
		return Spec{}, err
	}

	// Parse buffer as yaml into map
	doc := Spec{}
	err = yaml.Unmarshal(data, &doc)
	if err != nil {
		return Spec{}, err
	}

	return doc, nil
}
