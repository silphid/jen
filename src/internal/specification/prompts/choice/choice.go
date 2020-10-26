package choice

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
	"github.com/Samasource/jen/internal/specification/prompts"
)

type Choice struct {
	Display string
	Value   string
}

type Prompt struct {
	prompts.Prompt
	Var     string
	Default string
	Choices []Choice
}

func (p Prompt) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
