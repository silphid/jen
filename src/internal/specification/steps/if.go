package steps

import (
	"fmt"
	"github.com/Samasource/jen/internal"
	"github.com/Samasource/jen/internal/evaluation"
	"github.com/Samasource/jen/internal/specification"
)

type If struct {
	Condition string
	Then      specification.Executables
}

func (i If) String() string {
	return "do"
}

func (i If) Execute(context specification.Context) error {
	result, err := evaluation.EvalBoolExpression(context.Values, i.Condition)
	if err != nil {
		return fmt.Errorf("evaluating if conditional: %w", err)
	}
	if !result {
		internal.Log("Skipping sub-steps because condition %q evaluates to false", i.Condition)
		return nil
	}
	internal.Log("Executing sub-steps because condition %q evaluates to true", i.Condition)
	return i.Then.Execute(context)
}
