package choice

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
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

func (p Prompt) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
