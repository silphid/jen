package render

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
)

type Render struct {
	Directory string
}

func (r Render) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
