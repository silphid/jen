package internal

import "fmt"

var (
	TemplatesDir string
	Verbose bool
)

func Logf(message string, a ...interface{}) {
	if Verbose {
		fmt.Printf(message, a...)
		fmt.Println()
	}
}

func Log(message string) {
	if Verbose {
		fmt.Print(message)
		fmt.Println()
	}
}