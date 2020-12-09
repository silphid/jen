package do

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
)

type Do struct {
	Action string
}

func (d Do) String() string {
	return "do"
}

func (d Do) Execute(context specification.Context) error {
	action, ok := context.Spec.Actions[d.Action]
	if !ok {
		return fmt.Errorf("action %q not found for do step", d.Action)
	}
	return action.Execute(context)
}
