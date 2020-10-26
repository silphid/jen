package options

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
	"github.com/Samasource/jen/internal/specification/prompts"
)

type Option struct {
	Display string
	Var     string
	Default string
}

type Prompt struct {
	prompts.Prompt
	Options []Option
}

func (p Prompt) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
