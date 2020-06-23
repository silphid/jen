package cmd

import (
	"fmt"
	"github.com/Samasource/jen/cmd/gen"
	"github.com/Samasource/jen/internal"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "jen",
		Short: "Jen is a code generator for scaffolding microservices from Go templates boasting best practices.",
		Long:  `Jen is a code generator for scaffolding microservices from Go templates boasting best practices.`,
		SilenceUsage: true,
	}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&internal.TemplatesDir, "template-dir", "", "location of templates directory")
	rootCmd.AddCommand(gen.Cmd)
}

func Execute() error {
	return rootCmd.Execute()
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".jen")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
