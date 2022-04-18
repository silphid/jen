package templates

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/silphid/jen/cmd/jen/cmd/internal"
	"github.com/silphid/jen/cmd/jen/internal/home"
	"github.com/silphid/jen/cmd/jen/internal/spec"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:     "templates",
		Aliases: []string{"template"},
		Short:   "Lists templates available in git clone",
		Args:    cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return run(options, args)
		},
	}
}

func run(options *internal.Options, args []string) error {
	_, err := home.GetOrCloneRepo()
	if err != nil {
		return err
	}

	templatesDir, err := home.GetTemplatesDir()
	if err != nil {
		return err
	}

	// Read templates dir
	infos, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		return fmt.Errorf("reading templates directory %q: %w", templatesDir, err)
	}

	// Print templates with descriptions
	for _, info := range infos {
		template := info.Name()
		if strings.HasPrefix(template, ".") {
			continue
		}
		templateDir := filepath.Join(templatesDir, template)
		spec, err := spec.Load(templateDir)
		if err != nil {
			return err
		}
		fmt.Printf("%s - %s\n", template, spec.Description)
	}
	return nil
}
