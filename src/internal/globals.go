package internal

import "fmt"

var (
	TemplatesDir string
	Verbose      bool
)

func Log(message string, a ...interface{}) {
	if Verbose {
		fmt.Printf(message, a...)
		fmt.Println()
	}
}
