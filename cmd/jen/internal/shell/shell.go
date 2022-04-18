package shell

import (
	"os"
	"os/exec"
	"strings"

	"github.com/silphid/jen/cmd/jen/internal/logging"
)

// Execute executes one or multiple shell commands, given specific env vars (defaults to process' env vars if nil)
// and working directory (defaults to process's current work dir if empty string passed)
func Execute(vars []string, dir string, commands ...string) error {
	// Env vars default to current process' env vars
	if vars == nil {
		vars = os.Environ()
	}

	// Configure command struct
	cmd := &exec.Cmd{
		Path:   "/bin/bash",
		Args:   []string{"/bin/bash", "-c", "set -e; " + strings.Join(commands, "; ")},
		Dir:    dir,
		Env:    vars,
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
