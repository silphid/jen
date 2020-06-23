package model

import "fmt"

type Executable interface {
	Execute(context Context) error
}

func (root Spec) Execute(context Context) error {
	for _, step := range root.Steps {
		if err := step.Execute(context); err != nil {
			return err
		}
	}
	return nil
}

func (step StepUnion) Execute(context Context) error {
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
	default:
		fmt.Print("Ignoring unsupported command (for now!)")
	}
	return err
}
