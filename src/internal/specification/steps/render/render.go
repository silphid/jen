package render

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
	"github.com/Samasource/jen/internal/specification/steps"
)

type Render struct {
	If     string
	Source string
}

func (r Render) String() string {
	return "render"
}

func (r Render) Execute(context specification.Context) error {
	ok, err := steps.ShouldExecute(r.String(), r.If, context.Values)
	if !ok || err != nil {
		return err
	}

	return fmt.Errorf("not implemented")
}
