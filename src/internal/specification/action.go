package specification

import (
	"fmt"
	"github.com/Samasource/jen/internal"
)

type Action struct {
	Name  string
	Steps []Executable
}

func (p Action) Execute(context internal.Context) error {
	return fmt.Errorf("not implemented")
}
