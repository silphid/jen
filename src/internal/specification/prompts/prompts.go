package prompts

import (
	"github.com/Samasource/jen/internal/specification"
)

type Type string

const (
	Text    Type = "text"
	Option  Type = "option"
	Options Type = "options"
	Choice  Type = "choice"
)

type Prompt struct {
	specification.Step
	Question string
}
