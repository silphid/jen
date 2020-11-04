package evaluation

import (
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/Samasource/jen/internal"
	"os"
	"path"
	"text/template"
)

//func render(context Context, relativeInputDir string) error {
//	inputDir, err := filepath.Abs(path.Join(context.TemplateDir, relativeInputDir))
//	if err != nil {
//		return err
//	}
//	outputDir, err := filepath.Abs(context.OutputDir)
//	if err != nil {
//		return err
//	}
//	return renderDir(context, inputDir, outputDir)
//}
//
//func renderDir(context Context, inputPath, outputPath string) error {
//	Logf("Rendering dir %q -> %q", inputPath, outputPath)
//	infos, err := ioutil.ReadDir(inputPath)
//	if err != nil {
//		return err
//	}
//	if err := createOutputDir(outputPath); err != nil {
//		return err
//	}
//	for _, info := range infos {
//		outputName, include, err := resolveName(context, info.Name())
//		if err != nil {
//			return err
//		}
//		fullInput := path.Join(inputPath, info.Name())
//		fullOutput := path.Join(outputPath, outputName)
//		if !include {
//			Logf("Skipping %q because conditional evaluates to false", fullInput)
//			continue
//		}
//		if info.IsDir() {
//			err = renderDir(context, fullInput, fullOutput)
//		} else {
//			err = renderFile(context, fullInput, fullOutput)
//		}
//		if err != nil {
//			return err
//		}
//	}
//	return nil
//}

func renderFile(values Values, inputPath, outputPath string) error {
	internal.Logf("Rendering file %q -> %q", inputPath, outputPath)
	tmpl, err := template.New(path.Base(inputPath)).Funcs(sprig.TxtFuncMap()).ParseFiles(inputPath)
	if err != nil {
		return err
	}
	f, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file for template %v: %w", inputPath, err)
	}
	err = tmpl.Execute(f, values)
	if err != nil {
		return fmt.Errorf("render template %v: %w", inputPath, err)
	}
	return f.Close()
}
