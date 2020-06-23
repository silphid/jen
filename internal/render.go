package internal

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

func render(context Context, relativeInputDir string) error {
	inputDir, err := filepath.Abs(path.Join(context.TemplateDir, relativeInputDir))
	if err != nil {
		return err
	}
	outputDir, err := filepath.Abs(context.OutputDir)
	if err != nil {
		return err
	}
	return renderDir(context, inputDir, outputDir)
}

func renderDir(context Context, inputPath, outputPath string) error {
	infos, err := ioutil.ReadDir(inputPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		return err
	}
	for _, info := range infos {
		outputName, err := resolveName(context, info.Name())
		if err != nil {
			return err
		}
		fullInput := path.Join(inputPath, info.Name())
		fullOutput := path.Join(outputPath, outputName)
		if info.IsDir() {
			err = renderDir(context, fullInput, fullOutput)
		} else {
			err = renderFile(context, fullInput, fullOutput)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func renderFile(context Context, inputPath, outputPath string) error {
	tmpl, err := template.New(path.Base(inputPath)).Funcs(sprig.TxtFuncMap()).ParseFiles(inputPath)
	if err != nil {
		return err
	}
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file for template %v: %v", inputPath, err)
	}
	err = tmpl.Execute(f, context.Values)
	if err != nil {
		return fmt.Errorf("render template %v: %v", inputPath, err)
	}
	return f.Close()
}

func resolveName(context Context, name string) (string, error) {
	if strings.Index(name, "{{") == -1 {
		return name, nil
	}
	tmpl, err := template.New("base").Parse(name)
	if err != nil {
		return "", err
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, context.Values)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}
