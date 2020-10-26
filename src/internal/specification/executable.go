package specification

import "github.com/Samasource/jen/internal"

type Executable interface {
	Execute(context internal.Context) error
}
