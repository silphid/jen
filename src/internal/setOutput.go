/*package internal

func setOutput(context *Context, dir string) error {
	if context.OutputDir != "" {
		Logf("Skipping setOutput step because output dir was overriden to %q", context.OutputDir)
		return nil
	}

	dir, err := EvalTemplate(*context, dir)
	if err != nil {
		return err
	}

	Logf("Setting output dir to %q", dir)
	context.OutputDir = dir
	return nil
}
*/
