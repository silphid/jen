package choice

import (
	"github.com/Samasource/jen/internal"
	"github.com/Samasource/jen/internal/steps/prompts"
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

func (p *Prompt) Execute(context internal.Context) {
}
