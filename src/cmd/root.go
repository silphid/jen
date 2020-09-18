package cmd

import (
	"fmt"
	"github.com/Samasource/jen/cmd/create"
	"github.com/Samasource/jen/internal"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"strings"
)

var (
	rootCmd = &cobra.Command{
		Use:   "jen",
		Short: "Jen is a code generator for scaffolding microservices from Go templates boasting best practices.",
		Long:  `Jen is a code generator for scaffolding microservices from Go templates boasting best practices.`,
		SilenceUsage: true,
	}
	configFile string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&internal.TemplatesDir, "templates", "", "location of templates directory (default is ~/.jen/templates)")
	rootCmd.PersistentFlags().BoolVarP(&internal.Verbose, "verbose", "v", false, "display verbose messages")
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is ~/.jen/config.yaml)")
	rootCmd.AddCommand(create.Cmd)
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
	jenDir := path.Join(home, internal.JenDirName)
	viper.AddConfigPath(jenDir)
	viper.SetConfigName(internal.ConfigFileName)
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("jen")
	viper.AutomaticEnv()

	viper.SetDefault("templates", path.Join(jenDir, "templates"))

	if err := viper.ReadInConfig(); err == nil {
		internal.Logf("Using config file:", viper.ConfigFileUsed())
	} else {
		internal.Log("Config file not found")
	}

	if internal.TemplatesDir == "" {
		internal.TemplatesDir = viper.GetString("templates")
	}

	internal.TemplatesDir = strings.Replace(internal.TemplatesDir, "~", home, -1)
}
