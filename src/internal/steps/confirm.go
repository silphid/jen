package steps

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/silphid/jen/src/internal/evaluation"
	"github.com/silphid/jen/src/internal/exec"
	logging "github.com/silphid/jen/src/internal/logging"
)

// Confirm represents a conditional step that executes its child executable only if
// user answers Yes when prompted for given message
type Confirm struct {
	Message string
	Then    exec.Executables
}

func (c Confirm) String() string {
	return "confirm"
}

// Execute executes a child action only if user answers Yes when prompted for given message
func (c Confirm) Execute(context exec.Context) error {
	message, err := evaluation.EvalTemplate(context, c.Message)
	if err != nil {
		return err
	}
	prompt := &survey.Confirm{
		Message: message,
		Default: false,
	}
	value := false
	if err := survey.AskOne(prompt, &value); err != nil {
		return err
	}
	if !value {
		logging.Log("Skipping sub-steps because user cancelled")
		return nil
	}
	logging.Log("Executing sub-steps because user confirmed")
	return c.Then.Execute(context)
}
