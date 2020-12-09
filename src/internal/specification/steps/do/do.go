package do

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
)

type Do struct {
	Action string
}

func (d Do) String() string {
	return "do"
}

func (d Do) Execute(context specification.Context) error {
	return fmt.Errorf("not implemented")
}
