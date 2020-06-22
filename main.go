package main

import (
	"fmt"
	"github.com/Samasource/jen/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Printf("Error: %e", err)
		os.Exit(-1)
	}
}