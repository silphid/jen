package option

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
	"github.com/Samasource/jen/internal/specification/steps"
)

type Prompt struct {
	If       string
	Question string
	Var      string
	Default  bool
}

func (p Prompt) String() string {
	return "option"
}

func (p Prompt) Execute(context specification.Context) error {
	ok, err := steps.ShouldExecute(p.String(), p.If, context.Values)
	if !ok || err != nil {
		return err
	}

	return fmt.Errorf("not implemented")
}
