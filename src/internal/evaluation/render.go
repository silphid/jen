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

	entries, err := getEntries(context, inputDir, outputDir, renderMode)
	if err != nil {
		return fmt.Errorf("failed to determine entries to render: %w", err)
	}

	for _, entry := range entries {
		err = renderFile(context, entry.input, entry.output, entry.mode)
		if err != nil {
			return err
		}
	}

	return nil
}

type entry struct {
	input  string
	output string
	mode   RenderMode
}

func getEntries(context Context, inputDir, outputDir string, parentMode RenderMode) ([]entry, error) {
	var entries []entry
	infos, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		// Determine input/output names and render mode
		inputName := info.Name()
		inputPath := filepath.Join(inputDir, inputName)
		outputName, included, mode, err := evalFileName(context, inputName)
		if err != nil {
			return nil, err
		}
		outputPath := filepath.Join(outputDir, outputName)

		// Skip item?
		if !included {
			continue
		}

		// Mode defaults to parent's mode
		if mode == DefaultMode {
			mode = parentMode
		}

		// Directory?
		if info.IsDir() {
			if mode == InsertMode {
				return nil, fmt.Errorf("the .insert extension is not supported for directories: %q", inputName)
			}
			children, err := getEntries(context, inputPath, outputPath, mode)
			if err != nil {
				return nil, err
			}
			entries = append(entries, children...)
		} else {
			entries = append(entries, entry{
				input:  inputPath,
				output: outputPath,
				mode:   mode,
			})
		}
	}
	return entries, nil
}

func renderFile(context Context, inputPath, outputPath string, renderMode RenderMode) error {
	logging.Log("Rendering file %q -> %q", inputPath, outputPath)

	// Read input file
	inputText, err := ioutil.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read template file %q: %w", inputPath, err)
	}

	// Render input as template or copy as-is
	var outputText string
	if renderMode == TemplateMode {
		// Render file as template
		outputText, err = EvalTemplate(context, string(inputText))
		if err != nil {
			return fmt.Errorf("failed to render template %q: %w", inputPath, err)
		}
	} else if renderMode == InsertMode {
		// Parse insertion template
		insert, err := NewInsert(string(inputText))
		if err != nil {
			return fmt.Errorf("failed to parse insertion template %q: %w", inputPath, err)
		}
		// Read target file
		targetText, err := ioutil.ReadFile(outputPath)
		if err != nil {
			return fmt.Errorf("failed to read insertion target file %q: %w", outputPath, err)
		}
		// Perform insertion
		outputText, err = insert.Eval(context, string(targetText))
		if err != nil {
			return fmt.Errorf("failed to insert template %q into target file %q: %w", inputPath, outputPath, err)
		}
	} else {
		// Copy file as-is
		outputText = string(inputText)
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
