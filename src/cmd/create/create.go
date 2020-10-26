/*package create

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Samasource/jen/internal"
	"github.com/spf13/cobra"
	"io/ioutil"
	"path"
)

var (
	Cmd = &cobra.Command{
		Use:   "create",
		Short: "Generates a new project by executing given template",
		RunE: run,
	}

	template  string
	outputDir string
)

func init() {
	Cmd.PersistentFlags().StringVarP(&template, "template", "t", "", "name of template to use")
	Cmd.PersistentFlags().StringVarP(&outputDir, "output", "o", "", "output directory")
}

func run(_ *cobra.Command, _ []string) error {
	// Prompt for template
	if template == "" {
		var err error
		template, err = promptTemplate(internal.TemplatesDir)
		if err != nil {
			return err
		}
	}

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
	if err := context.Spec.Execute(&context); err != nil {
		return err
	}

	return nil
}

func promptTemplate(templatesDir string) (string, error) {
	// Read templates dir
	infos, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		return "", err
	}

	// Build list of choices
	var templates []string
	var titles []string
	for _, info := range infos {
		template := info.Name()
		templateDir := path.Join(internal.TemplatesDir, template)
		spec, err := internal.Load(templateDir)
		if err == nil {
			templates = append(templates, template)
			titles = append(titles, fmt.Sprintf("%s - %s", template, spec.Description))
		}
	}

	// Any templates found?
	if len(templates) == 0 {
		return "", fmt.Errorf("no templates found in %q", templatesDir)
	}

	// Prompt
	prompt := &survey.Select{
		Message: "Select template",
		Options: titles,
	}
	var index int
	if err := survey.AskOne(prompt, &index); err != nil {
		return "", err
	}

	return templates[index], nil
}*/
