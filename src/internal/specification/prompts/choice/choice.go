package choice

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
)

type Choice struct {
	Display string
	Value   string
}

type Prompt struct {
	If       string
	Question string
	Var      string
	Default  string
	Choices  []Choice
}

func (p Prompt) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
