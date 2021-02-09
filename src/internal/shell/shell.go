package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Samasource/jen/src/internal/helpers"
	"github.com/Samasource/jen/src/internal/logging"
	"github.com/Samasource/jen/src/internal/model"
)

// Execute executes one or multiple shell commands, while injecting project's environment variables and bin directories path
func Execute(vars model.VarMap, dir string, binDirs []string, commands ...string) error {
	// Configure command struct
	cmd := &exec.Cmd{
		Path:   "/bin/bash",
		Args:   []string{"/bin/bash", "-c", "set -e; " + strings.Join(commands, "; ")},
		Dir:    dir,
		Env:    GetEnvFromProcessAndProjectVariables(vars, binDirs),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// Execute
	logging.Log("Executing command(s) %q in directory %q", commands, dir)
	logging.Log("--")
	defer logging.Log("--")
	return cmd.Run()
}

// GetEnvFromProcessAndProjectVariables computes environment variables based on current process environment and project variables
func GetEnvFromProcessAndProjectVariables(vars model.VarMap, binDirs []string) []string {
	// Add bin dirs to PATH env var
	pathVar := os.Getenv("PATH")
	for _, dir := range binDirs {
		if helpers.PathExists(dir) {
			pathVar = dir + ":" + pathVar
		}
	}

	// Collect all current process env vars, except PATH
	var env []string
	for _, entry := range os.Environ() {
		if !strings.HasPrefix(entry, "PATH=") {
			env = append(env, entry)
		}
	}

	// Override PATH env var
	entry := fmt.Sprintf("PATH=%v", pathVar)
	env = append(env, entry)
	logging.Log(entry)

	// Then values env vars
	logging.Log("Environment variables:")
	for key, value := range vars {
		entry := fmt.Sprintf("%s=%v", key, value)
		env = append(env, entry)
		logging.Log(entry)
	}
	return env
}
