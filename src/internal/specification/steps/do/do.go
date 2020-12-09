package do

import (
	"fmt"
	"github.com/Samasource/jen/internal/model"
)

type Do struct {
	Action string
}

func (d Do) String() string {
	return "do"
}

func (d Do) Execute(config model.Config) error {
	action, ok := config.Spec.Actions[d.Action]
	if !ok {
		return fmt.Errorf("action %q not found for do step", d.Action)
	}
	return action.Execute(config)
}
