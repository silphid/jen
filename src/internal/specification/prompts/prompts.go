package prompts

import (
	"github.com/Samasource/jen/internal/specification/step"
)

type Type string

const (
	Text    Type = "text"
	Option  Type = "option"
	Options Type = "options"
	Choice  Type = "choice"
)

type Prompt struct {
	step.Step
	Question string
}
