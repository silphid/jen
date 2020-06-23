package cmd

import (
	"fmt"
	"github.com/Samasource/jen/internal"
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "jen",
		Short: "Jen is a code generator for scaffolding microservices from Go templates boasting best practices.",
		Long:  `Jen is a code generator for scaffolding microservices from Go templates boasting best practices.`,
		RunE: run,
	}

	templatesDir string
	template string
)

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&templatesDir, "template-dir", "", "location of templates directory")
	rootCmd.PersistentFlags().StringVarP(&template, "template", "t", "", "name of template to use")
}

func run(_ *cobra.Command, _ []string) error {
	// Load spec
	templateDir := path.Join(templatesDir, template)
	spec, err := internal.Load(templateDir)
	if err != nil {
		return err
	}

	// Execute all spec steps
	context := internal.Context{
		Spec: spec,
		Values: make(internal.Values),
	}
	if err := context.Spec.Execute(context); err != nil {
		return err
	}

	fmt.Printf("Context:\n%v", context)
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		er(err)
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".jen")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
