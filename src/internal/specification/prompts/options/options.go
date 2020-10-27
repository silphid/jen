package options

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
)

type Option struct {
	Question string
	Var      string
	Default  bool
}

type Prompt struct {
	If       string
	Question string
	Options  []Option
}

func (p Prompt) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
