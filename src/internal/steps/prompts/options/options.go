package options

import (
	"github.com/Samasource/jen/internal"
	"github.com/Samasource/jen/internal/steps/prompts"
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

func (p *Prompt) Execute(context internal.Context) {
}
