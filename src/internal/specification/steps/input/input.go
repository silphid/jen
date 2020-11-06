package input

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
	"github.com/Samasource/jen/internal/specification/steps"
)

type Prompt struct {
	If       string
	Question string
	Var      string
	Default  string
}

func (p Prompt) Execute(context specification.Context) error {
	ok, err := steps.ShouldExecute("input", p.If, context.Values)
	if !ok || err != nil {
		return err
	}

	return fmt.Errorf("not implemented")
}
