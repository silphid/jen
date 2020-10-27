package do

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification/executable"
)

type Do struct {
	If     string
	Action string
}

func (p Do) Execute(context executable.Context) error {
	return fmt.Errorf("not implemented")
}
