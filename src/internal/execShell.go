package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func execShell(context Context, command string) error {
	outputDir, err := filepath.Abs(context.OutputDir)
	if err != nil {
		return err
	}
	if err := createOutputDir(outputDir); err != nil {
		return err
	}
	Logf("Executing command %q in dir %q", command, outputDir)
	cmd := exec.Command("bash", "-c", "set -e; " + command)
	cmd.Env = getEnvFromValues(context.Values)
	cmd.Dir = outputDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	Log("--")
	defer Log("--")
	return cmd.Run()
}

func getEnvFromValues(values Values) []string {
	// Pass current process env vars
	var env []string
	for _, entry := range os.Environ() {
		env = append(env, entry)
	}

	// Then values env vars
	Log("Environment variables:")
	for key, value := range values {
		entry := fmt.Sprintf("%s=%v", key, value)
		env = append(env, entry)
		Log(entry)
	}
	return env
}
