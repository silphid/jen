package do

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
	"github.com/Samasource/jen/internal/specification/steps"
)

type Do struct {
	If     string
	Action string
}

func (d Do) String() string {
	return "do"
}

func (d Do) Execute(context specification.Context) error {
	ok, err := steps.ShouldExecute(d.String(), d.If, context.Values)
	if !ok || err != nil {
		return err
	}

	return fmt.Errorf("not implemented")
}
