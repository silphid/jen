package shell

import (
	"fmt"
	"os"
	"strings"

	"github.com/silphid/jen/cmd/jen/cmd/internal"
	"github.com/silphid/jen/cmd/jen/internal/shell"
	"github.com/spf13/cobra"
)

// New creates a cobra command
func New(options *internal.Options) *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Starts a sub-shell with project's environment variables and scripts in PATH",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, args []string) error {
			return run(options, args)
		},
	}
}

func run(options *internal.Options, args []string) error {
	execContext, err := options.NewContext()
	if err != nil {
		return err
	}

	shellEnv := os.Getenv("SHELL")
	if shellEnv == "" {
		shellEnv = "/bin/bash"
	}

	flag := "--norc"
	if strings.HasSuffix(shellEnv, "/zsh") {
		flag = "--norcs"
	}

	cmd := fmt.Sprintf("%s %s", shellEnv, flag)
	return shell.Execute(execContext.GetShellVars(true), "", cmd)
}
