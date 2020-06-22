package cmd

import (
	"fmt"
	"github.com/Samasource/jen/spec"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	root, err := loadSpec()
	if err != nil {
		return err
	}

	fmt.Printf("%v", root)
	return nil
}

func loadSpec() (spec.Root, error) {
	// Load file to buffer
	data, err := ioutil.ReadFile("examples/jen.yaml")
	if err != nil {
		return spec.Root{}, err
	}

	// Parse buffer as yaml into map
	doc := spec.Root{}
	err = yaml.Unmarshal(data, &doc)
	if err != nil {
		return spec.Root{}, err
	}

	return doc, nil
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
