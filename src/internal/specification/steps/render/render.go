package render

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
)

type Render struct {
	If     string
	Source string
}

func (r Render) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
