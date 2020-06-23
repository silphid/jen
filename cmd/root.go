package cmd

import (
	"fmt"
	"github.com/Samasource/jen/internal/model"
	"os"

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
)

func run(_ *cobra.Command, _ []string) error {
	// Load spec
	spec, err := model.Load()
	if err != nil {
		return err
	}

	// Execute all spec steps
	context := model.Context{
		Spec: spec,
		Values: make(model.Values),
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

func init() {
	cobra.OnInitialize(initConfig)

	//	rootCmd.AddCommand(addCmd)
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
