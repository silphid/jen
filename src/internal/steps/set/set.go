package set

import (
	"github.com/Samasource/jen/src/internal/evaluation"
	"github.com/Samasource/jen/src/internal/exec"
	"github.com/Samasource/jen/src/internal/helpers/variables"
)

// Variable represents a single variable to be set to a given value
type Variable struct {
	Name  string
	Value string
}

// Set represents a step that sets one or multiple variables to given values without user intervention
type Set struct {
	Variables []Variable
}

func (p Set) String() string {
	return "set"
}

// Execute prompts user for input value
func (p Set) Execute(context exec.Context) error {
	vars := context.GetVars()

	for _, variable := range p.Variables {
		// Do not set variables that have been overriden at command-line
		if context.IsVarOverriden(variable.Name) {
			continue
		}

		// Compute value
		value, ok := variables.TryGetString(vars, variable.Name)
		if !ok {
			var err error
			value, err = evaluation.EvalTemplate(context, variable.Value)
			if err != nil {
				return err
			}
		}

		vars[variable.Name] = value
	}

	return context.SetVars(vars)
}
