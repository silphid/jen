package options

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
	"github.com/Samasource/jen/internal/specification/steps"
)

type Item struct {
	Text    string
	Var     string
	Default bool
}

type Prompt struct {
	If       string
	Question string
	Items    []Item
}

func (p Prompt) String() string {
	return "options"
}

func (p Prompt) Execute(context specification.Context) error {
	ok, err := steps.ShouldExecute(p.String(), p.If, context.Values)
	if !ok || err != nil {
		return err
	}

	return fmt.Errorf("not implemented")
}
