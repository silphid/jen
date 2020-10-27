package option

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
)

type Prompt struct {
	If       string
	Question string
	Var      string
	Default  bool
}

func (p Prompt) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
