package exec

import (
	"fmt"
	"github.com/Samasource/jen/internal"
	"github.com/Samasource/jen/internal/specification"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Exec struct {
	Commands []string
}

func (e Exec) String() string {
	return "exec"
}

func (e Exec) Execute(context specification.Context) error {
	dir, err := filepath.Abs(context.OutputDir)
	if err != nil {
		return err
	}

	// Concatenate commands
	builder := strings.Builder{}
	for _, command := range e.Commands {
		builder.WriteString(command)
		builder.WriteString("; ")
	}
	commands := builder.String()

	// Configure command struct
	cmd := exec.Command("bash", "-c", "set -e; "+commands)
	cmd.Env = getEnvFromValues(context.Values)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute
	internal.Log("Executing commands %q in directory %q", commands, dir)
	internal.Log("--")
	defer internal.Log("--")
	return cmd.Run()
}

func getEnvFromValues(values specification.Values) []string {
	// Pass current process env vars
	var env []string
	for _, entry := range os.Environ() {
		env = append(env, entry)
	}

	// Then values env vars
	internal.Log("Environment variables:")
	for key, value := range values.Variables {
		entry := fmt.Sprintf("%s=%v", key, value)
		env = append(env, entry)
		internal.Log(entry)
	}
	return env
}
