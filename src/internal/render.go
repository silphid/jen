/*package internal

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
)

var doubleBracketRegexp = regexp.MustCompile(`\[\[.*]]`)

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
	Logf("Rendering dir %q -> %q", inputPath, outputPath)
	infos, err := ioutil.ReadDir(inputPath)
	if err != nil {
		return err
	}
	if err := createOutputDir(outputPath); err != nil {
		return err
	}
	for _, info := range infos {
		outputName, include, err := resolveName(context, info.Name())
		if err != nil {
			return err
		}
		fullInput := path.Join(inputPath, info.Name())
		fullOutput := path.Join(outputPath, outputName)
		if !include {
			Logf("Skipping %q because conditional evaluates to false", fullInput)
			continue
		}
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
	Logf("Rendering file %q -> %q", inputPath, outputPath)
	tmpl, err := template.New(path.Base(inputPath)).Funcs(sprig.TxtFuncMap()).ParseFiles(inputPath)
	if err != nil {
		return err
	}
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file for template %v: %w", inputPath, err)
	}
	err = tmpl.Execute(f, context.Values)
	if err != nil {
		return fmt.Errorf("render template %v: %w", inputPath, err)
	}
	return f.Close()
}

func resolveName(context Context, name string) (string, bool, error) {
	// Double-bracket expression (ie: "[[.option]]") in names are evaluated to determine
	// whether the file/folder should be rendered and that expression then gets stripped
	// from the name
	for {
		// Find expression
		loc := doubleBracketRegexp.FindStringIndex(name)
		if loc == nil {
			break
		}
		exp := name[loc[0]+2 : loc[1]-2]

		// Evaluate expression
		value, err := EvalBoolExpression(context, exp)
		if err != nil {
			return "", false, fmt.Errorf("eval double-bracket expression in name %q: %w", name, err)
		}

		// Should we exclude file/folder?
		if !value {
			return "", false, nil
		}

		// Remove expression from name
		name = name[:loc[0]] + name[loc[1]:]
	}

	// Double-brace expression (ie: "{{.name}}") in names get interpolated as expected
	if strings.Index(name, "{{") != -1 {
		tmpl, err := template.New("base").Parse(name)
		if err != nil {
			return "", false, fmt.Errorf("parse double-brace expression in name %q: %w", name, err)
		}
		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, context.Values)
		if err != nil {
			return "", false, fmt.Errorf("render double-brace expression in name %q: %w", name, err)
		}
		return buffer.String(), true, nil
	}

	return name, true, nil
}
*/
