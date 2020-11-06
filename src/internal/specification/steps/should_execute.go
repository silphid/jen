package steps

import (
	"fmt"
	"github.com/Samasource/jen/internal"
	"github.com/Samasource/jen/internal/evaluation"
)

func ShouldExecute(name, condition string, values evaluation.Values) (bool, error) {
	if condition != "" {
		result, err := evaluation.EvalBoolExpression(values, condition)
		if err != nil {
			return false, fmt.Errorf("evaluate step %q conditional expression: %w", name, err)
		}
		if !result {
			internal.Log("Skipping step %q because condition %q evaluates to false", name, condition)
			return false, nil
		}
	}
	internal.Log("Executing step %q", name)
	return true, nil
}
