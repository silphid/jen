package shell

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	. "github.com/Samasource/jen/internal/logging"
	"github.com/Samasource/jen/internal/model"
)

func Execute(vars model.VarMap, dir, pathEnvVar string, commands ...string) error {
	// Configure command struct
	cmd := &exec.Cmd{
		Path:   "/bin/bash",
		Args:   []string{"/bin/bash", "-c", "set -e; " + strings.Join(commands, "; ")},
		Dir:    dir,
		Env:    GetEnvFromValues(vars, pathEnvVar),
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

func GetEnvFromValues(vars model.VarMap, pathEnvVar string) []string {
	// Pass current process env vars
	var env []string
	for _, entry := range os.Environ() {
		if pathEnvVar == "" || !strings.HasPrefix(entry, "PATH=") {
			env = append(env, entry)
		}
	}

	// Overriden PATH env var
	if pathEnvVar != "" {
		env = append(env, "PATH="+pathEnvVar)
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
