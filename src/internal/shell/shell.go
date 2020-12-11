package shell

import (
	"fmt"
	. "github.com/Samasource/jen/internal/logging"
	"github.com/Samasource/jen/internal/model"
	"os"
	"os/exec"
	"strings"
)

func Execute(vars model.VarMap, dir string, commands ...string) error {
	// Configure command struct
	cmd := &exec.Cmd{
		Path:   "/bin/bash",
		Args:   []string{"/bin/bash", "-c", "set -e; " + strings.Join(commands, "; ")},
		Dir:    dir,
		Env:    getEnvFromValues(vars),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// Execute
	Log("Executing command(s) %q in directory %q", commands, dir)
	Log("--")
	defer Log("--")
	return cmd.Run()
}

func getEnvFromValues(vars model.VarMap) []string {
	// Pass current process env vars
	var env []string
	for _, entry := range os.Environ() {
		env = append(env, entry)
	}

	// Then values env vars
	Log("Environment variables:")
	for key, value := range vars {
		entry := fmt.Sprintf("%s=%v", key, value)
		env = append(env, entry)
		Log(entry)
	}
	return env
}
