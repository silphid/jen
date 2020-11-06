package exec

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
	"github.com/Samasource/jen/internal/specification/steps"
)

type Exec struct {
	If      string
	Command string
}

func (e Exec) String() string {
	return "exec"
}

func (e Exec) Execute(context specification.Context) error {
	ok, err := steps.ShouldExecute(e.String(), e.If, context.Values)
	if !ok || err != nil {
		return err
	}

	return fmt.Errorf("not implemented")
}
