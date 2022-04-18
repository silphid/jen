package logging

import "fmt"

var (
	Verbose bool
)

func Log(message string, a ...interface{}) {
	if Verbose {
		fmt.Printf(message, a...)
		fmt.Println()
	}
}
