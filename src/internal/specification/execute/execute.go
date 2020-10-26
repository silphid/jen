package execute

import (
	"fmt"
	"github.com/Samasource/jen/internal"
)

type Execute struct {
	Command string
}

func (p Execute) Execute(context internal.Context) error {
	return fmt.Errorf("not implemented")
}
