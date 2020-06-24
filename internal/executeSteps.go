package internal

import (
	"fmt"
)

type Executable interface {
	Execute(context Context) error
}

func (root Spec) Execute(context Context) error {
	return execute(context, root.Steps)
}

func execute(context Context, steps []*StepUnion) error {
	for i, step := range steps {
		if err := step.execute(i, context); err != nil {
			return err
		}
	}
	return nil
}

func (step StepUnion) execute(index int, context Context) error {
	if step.If != "" {
		result, err := EvalExpression(context, step.If)
		if err != nil {
			return fmt.Errorf("evaluate step #%d conditional expression: %w", index, err)
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
	case step.Secret != nil:
		err = step.Secret.Execute(context)
	case step.Option != nil:
		err = step.Option.Execute(context)
	case step.Multi != nil:
		err = step.Multi.Execute(context)
	case step.Select != nil:
		err = step.Select.Execute(context)
	case step.Render != "":
		err = render(context, step.Render)
	case step.Do != "":
		err = do(context, step.Do)
	case step.Exec != "":
		err = exec(context, step.Exec)
	default:
		return fmt.Errorf("unsupported step #%d", index)
	}
	return err
}
