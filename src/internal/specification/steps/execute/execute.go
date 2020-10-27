package execute

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
)

type Execute struct {
	If      string
	Command string
}

func (p Execute) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
