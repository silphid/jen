package text

import (
	"fmt"
	"github.com/Samasource/jen/internal"
	"github.com/Samasource/jen/internal/specification/prompts"
)

type Prompt struct {
	prompts.Prompt
	Var     string
	Default string
}

func (p Prompt) Execute(context internal.Context) error {
	return fmt.Errorf("not implemented")
}
