package internal

import (
	"fmt"
)

type Executable interface {
	Execute(context Context) error
}

func (root Spec) Execute(context Context) error {
	for i, step := range root.Steps {
		if err := step.Execute(i, context); err != nil {
			return err
		}
	}
	return nil
}

func (step StepUnion) Execute(index int, context Context) error {
	if step.If != "" {
		result, err := EvalExpression(context, step.If)
		if err != nil {
			return fmt.Errorf("evaluate if expression for step #%d: %v", index, err)
		}
		if !result {
			Logf("Skipping step #%d because conditional %q evaluates to false", index, step.If)
			return nil
		}
	}

	var err error
	switch {
	case step.Value != nil:
		err = step.Value.Execute(context)
	case step.Option != nil:
		err = step.Option.Execute(context)
	case step.Multi != nil:
		err = step.Multi.Execute(context)
	case step.Select != nil:
		err = step.Select.Execute(context)
	case step.Render != "":
		err = render(context, step.Render)
	default:
		fmt.Print("Ignoring unsupported command (for now!)")
	}
	return err
}
