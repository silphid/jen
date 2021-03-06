package evaluation

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Samasource/jen/src/internal/logging"
)

// Render copies all files from inputDir into outputDir, rendering as templates those for which rendering is enabled
// interpolating folder and file names appropriately and skipping folders and files for which bracket expressions
// evaluate to false
func Render(context Context, inputDir, outputDir string) error {
	// Determine if rendering should be turned on from the start
	renderMode, _ := getRenderModeAndRemoveExtension(inputDir)
	renderRecursively := renderMode == EnableRendering

	entries, err := getEntries(context, inputDir, outputDir, renderRecursively)
	if err != nil {
		return fmt.Errorf("failed to determine entries to render: %w", err)
	}

	for _, entry := range entries {
		err = renderFile(context, entry.input, entry.output, entry.render)
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

func getEntries(context Context, inputDir, outputDir string, renderParent bool) ([]entry, error) {
	var entries []entry
	infos, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		// Determine input/output names and render mode
		inputName := info.Name()
		inputPath := filepath.Join(inputDir, inputName)
		outputName, included, renderMode, err := evalFileName(context, inputName)
		if err != nil {
			return nil, err
		}
		outputPath := filepath.Join(outputDir, outputName)

		// Determine render enabled/disabled for this item
		renderItem := renderParent
		if renderMode == EnableRendering {
			renderItem = true
		} else if renderMode == DisableRendering {
			renderItem = false
		}

		// Skip item?
		if !included {
			continue
		}

		// Directory?
		if info.IsDir() {
			children, err := getEntries(context, inputPath, outputPath, renderItem)
			if err != nil {
				return nil, err
			}
			entries = append(entries, children...)
		} else {
			entries = append(entries, entry{
				input:  inputPath,
				output: outputPath,
				render: renderItem,
			})
		}
	}
	return entries, nil
}

func renderFile(context Context, inputPath, outputPath string, render bool) error {
	logging.Log("Rendering file %q -> %q", inputPath, outputPath)

	// Read input file
	inputText, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	// Render input as template or copy as-is
	outputText := string(inputText)
	if render {
		outputText, err = EvalTemplate(context, outputText)
		if err != nil {
			return fmt.Errorf("failed to render template %v: %w", inputPath, err)
		}
	}

	// Create output dir
	outputDir := filepath.Dir(outputPath)
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
