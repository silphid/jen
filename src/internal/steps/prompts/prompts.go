package prompts

import "github.com/Samasource/jen/internal/steps"

type Type string

const (
	Text    Type = "text"
	Option  Type = "option"
	Options Type = "options"
	Choice  Type = "choice"
)

type Prompt struct {
	steps.Step
	Question string
	Type     Type
}
