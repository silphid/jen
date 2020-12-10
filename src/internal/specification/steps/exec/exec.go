package exec

import (
	"fmt"
	. "github.com/Samasource/jen/internal/logging"
	"github.com/Samasource/jen/internal/model"
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

func (e Exec) Execute(config *model.Config) error {
	dir, err := filepath.Abs(config.ProjectDir)
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
	cmd.Env = getEnvFromValues(config.Values)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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
