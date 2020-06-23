package main

import (
	"github.com/Samasource/jen/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(-1)
	}
}