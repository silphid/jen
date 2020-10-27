package input

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
)

type Prompt struct {
	If       string
	Question string
	Var      string
	Default  string
}

func (p Prompt) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
