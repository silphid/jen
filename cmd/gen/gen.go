package gen

import (
	"github.com/Samasource/jen/internal"
	"github.com/spf13/cobra"
	"path"
)

var (
	Cmd = &cobra.Command{
		Use:   "gen",
		Short: "Generates a new project by executing given template",
		RunE: run,
	}

	template  string
	outputDir string
)

func init() {
	Cmd.PersistentFlags().StringVarP(&template, "template", "t", "", "name of template to use")
	Cmd.PersistentFlags().StringVarP(&outputDir, "output", "o", ".", "output directory")
}

func run(_ *cobra.Command, _ []string) error {
	// Load spec
	templateDir := path.Join(internal.TemplatesDir, template)
	spec, err := internal.Load(templateDir)
	if err != nil {
		return err
	}

	// Execute all spec steps
	context := internal.Context{
		TemplateDir: templateDir,
		OutputDir:   outputDir,
		Spec:        spec,
		Values:      make(internal.Values),
	}
	if err := context.Spec.Execute(context); err != nil {
		return err
	}

	return nil
}
