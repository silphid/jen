package cmd

import (
	"fmt"
	"github.com/Samasource/jen/internal"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:          "jen",
		Short:        "Jen is a code generator for scaffolding microservices from Go templates boasting best practices.",
		Long:         `Jen is a code generator for scaffolding microservices from Go templates boasting best practices.`,
		SilenceUsage: true,
	}
	configFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&internal.Verbose, "verbose", "v", false, "display verbose messages")
	//rootCmd.AddCommand(create.Cmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Failed to determine home dir: %v", err)
		os.Exit(1)
	}
	value, ok := os.LookupEnv("JEN_TEMPLATES")
	if !ok {
		value = "~/.jen/templates"
	}
	internal.TemplatesDir = strings.ReplaceAll(value, "~", home)
	internal.Log("Using templates in: %s", internal.TemplatesDir)
}
