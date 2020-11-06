package choice

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
	"github.com/Samasource/jen/internal/specification/steps"
)

type Item struct {
	Text  string
	Value string
}

type Prompt struct {
	If       string
	Question string
	Var      string
	Default  string
	Items    []Item
}

func (p Prompt) String() string {
	return "choice"
}

func (p Prompt) Execute(context specification.Context) error {
	ok, err := steps.ShouldExecute(p.String(), p.If, context.Values)
	if !ok || err != nil {
		return err
	}

	return fmt.Errorf("not implemented")
}
