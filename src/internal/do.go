/*package internal

import "fmt"

func do(context *Context, action string) error {
	steps, ok := context.Spec.Actions[action]
	if !ok {
		return fmt.Errorf("action not found %q", action)
	}
	Logf("Doing action %q", action)
	return execute(context, steps)
}
*/
