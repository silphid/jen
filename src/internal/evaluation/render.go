package evaluation

import (
	"fmt"
	"github.com/Samasource/jen/internal"
	"github.com/Samasource/jen/internal/specification"
	"io/ioutil"
	"os"
	"path"
)

func Render(values specification.Values, inputDir, outputDir string) error {
	entries, err := getEntries(values, inputDir, outputDir)
	if err != nil {
		return fmt.Errorf("failed to determine entries to render: %w", err)
	}

	for _, entry := range entries {
		err = renderFile(values, entry.input, entry.output)
		if err != nil {
			return err
		}
	}

	return nil
}

type entry struct {
	input  string
	output string
}

func getEntries(values specification.Values, inputDir, outputDir string) ([]entry, error) {
	var entries []entry
	infos, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		inputName := info.Name()
		inputPath := path.Join(inputDir, inputName)
		outputName, included, err := evalFileName(values, inputName)
		if err != nil {
			return nil, err
		}
		outputPath := path.Join(outputDir, outputName)

		if !included {
			continue
		}
		if info.IsDir() {
			children, err := getEntries(values, inputPath, outputPath)
			if err != nil {
				return nil, err
			}
			entries = append(entries, children...)
		} else {
			entries = append(entries, entry{
				input:  inputPath,
				output: outputPath,
			})
		}
	}
	return entries, nil
}

func renderFile(values specification.Values, inputPath, outputPath string) error {
	internal.Log("Rendering file %q -> %q", inputPath, outputPath)

	// Read input file
	inputText, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Render template
	outputText, err := EvalTemplate(values, string(inputText))
	if err != nil {
		return fmt.Errorf("failed to render template %v: %w", inputPath, err)
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
