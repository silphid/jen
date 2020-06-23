package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
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
		fullInput := path.Join(inputPath, info.Name())
		fullOutput := path.Join(outputPath, info.Name())
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
	t, err := template.ParseFiles(inputPath)
	if err != nil {
		return err
	}
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file for template %v: %v", inputPath, err)
	}
	err = t.Execute(f, context.Values)
	if err != nil {
		return fmt.Errorf("render t %v: %v", inputPath, err)
	}
	return f.Close()
}
