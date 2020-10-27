package options

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
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

func (p Prompt) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
