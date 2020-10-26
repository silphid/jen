package specification

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
)

type Action struct {
	Name  string
	Steps []executable.Executable
}

func (p Action) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
