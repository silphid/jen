package evaluation

import (
	"fmt"
	. "github.com/Samasource/jen/internal/logging"
	"github.com/Samasource/jen/internal/model"
	"io/ioutil"
	"os"
	"path"
)

func Render(values model.Values, inputDir, outputDir string) error {
	renderRecursively := HasTmplExtension(inputDir)
	entries, err := getEntries(values, inputDir, outputDir, renderRecursively)
	if err != nil {
		return fmt.Errorf("failed to determine entries to render: %w", err)
	}

	for _, entry := range entries {
		err = renderFile(values, entry.input, entry.output, entry.render)
		if err != nil {
			return err
		}
	}

	return nil
}

type entry struct {
	input  string
	output string
	render bool
}

func getEntries(values model.Values, inputDir, outputDir string, renderRecursively bool) ([]entry, error) {
	var entries []entry
	infos, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		inputName := info.Name()
		inputPath := path.Join(inputDir, inputName)
		outputName, included, render, err := evalFileName(values, inputName)
		if err != nil {
			return nil, err
		}
		render = render || renderRecursively
		outputPath := path.Join(outputDir, outputName)

		if !included {
			continue
		}
		if info.IsDir() {
			children, err := getEntries(values, inputPath, outputPath, render)
			if err != nil {
				return nil, err
			}
			entries = append(entries, children...)
		} else {
			entries = append(entries, entry{
				input:  inputPath,
				output: outputPath,
				render: render,
			})
		}
	}
	return entries, nil
}

func renderFile(values model.Values, inputPath, outputPath string, render bool) error {
	Log("Rendering file %q -> %q", inputPath, outputPath)

	// Read input file
	inputText, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Render input as template or copy as-is
	outputText := string(inputText)
	if render {
		outputText, err = EvalTemplate(values, outputText)
		if err != nil {
			return fmt.Errorf("failed to render template %v: %w", inputPath, err)
		}
	}

	// Create output dir
	outputDir := path.Dir(outputPath)
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory %q: %w", outputDir, err)
	}

	// Write file
	err = ioutil.WriteFile(outputPath, []byte(outputText), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write rendered output file for template %v: %w", inputPath, err)
	}
	return nil
}
