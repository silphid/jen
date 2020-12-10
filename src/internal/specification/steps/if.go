package steps

import (
	"fmt"
	"github.com/Samasource/jen/internal/evaluation"
	. "github.com/Samasource/jen/internal/logging"
	"github.com/Samasource/jen/internal/model"
)

type If struct {
	Condition string
	Then      model.Executables
}

func (i If) String() string {
	return "do"
}

func (i If) Execute(config *model.Config) error {
	result, err := evaluation.EvalBoolExpression(config.Values, i.Condition)
	if err != nil {
		return fmt.Errorf("evaluating if conditional: %w", err)
	}
	if !result {
		Log("Skipping sub-steps because condition %q evaluates to false", i.Condition)
		return nil
	}
	Log("Executing sub-steps because condition %q evaluates to true", i.Condition)
	return i.Then.Execute(config)
}
