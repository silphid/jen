package text

import (
	"github.com/Samasource/jen/internal"
	"github.com/Samasource/jen/internal/steps/prompts"
)

type Prompt struct {
	prompts.Prompt
	Var     string
	Default string
}

func (p *Prompt) Execute(context internal.Context) {
}
