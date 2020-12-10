package exec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	. "github.com/Samasource/jen/internal/logging"
	"github.com/Samasource/jen/internal/model"
)

type Exec struct {
	Commands []string
}

func (e Exec) String() string {
	return "exec"
}

func (e Exec) Execute(config *model.Config) error {
	dir, err := filepath.Abs(config.ProjectDir)
	if err != nil {
		return err
	}

	// Concatenate commands
	builder := strings.Builder{}
	for i, command := range e.Commands {
		if i > 0 {
			builder.WriteString("; ")
		}
		builder.WriteString(command)
	}
	commands := builder.String()

	// Configure command struct
	cmd := &exec.Cmd{
		Path:   "/bin/bash",
		Args:   []string{"/bin/bash", "-c", "set -e; " + commands},
		Dir:    dir,
		Env:    getEnvFromValues(config.Values),
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// Execute
	Log("Executing commands %q in directory %q", commands, dir)
	Log("--")
	defer Log("--")
	return cmd.Run()
}

func getEnvFromValues(values model.Values) []string {
	// Pass current process env vars
	var env []string
	for _, entry := range os.Environ() {
		env = append(env, entry)
	}

	// Then values env vars
	Log("Environment variables:")
	for key, value := range values.Variables {
		entry := fmt.Sprintf("%s=%v", key, value)
		env = append(env, entry)
		Log(entry)
	}
	return env
}
