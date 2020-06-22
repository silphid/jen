package cmd

import (
"fmt"
"os"

"github.com/mitchellh/go-homedir"
"github.com/spf13/cobra"
"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command {
	Use:   "jen",
	Short: "Jen is a code generator for scaffolding microservices from Go templates boasting best practices.",
	Long: `Jen is a code generator for scaffolding microservices from Go templates boasting best practices.`,
	}
)

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