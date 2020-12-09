package exec

import (
	"fmt"
	"github.com/Samasource/jen/internal/specification"
)

type Exec struct {
	Command string
}

func (e Exec) String() string {
	return "exec"
}

func (e Exec) Execute(context specification.Context) error {
	return fmt.Errorf("not implemented")
}
