package internal

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func executeShellCommand(context Context, command string) error {
	Logf("Executing command %q", command)
	cmd := exec.Command("bash", "-c", command)
	cmd.Env = getEnvFromValues(context.Values)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getEnvFromValues(values Values) []string {
	var env []string
	Log("Environment variables:")
	for key, value := range values {
		entry := fmt.Sprintf("%s=%v", toSnakeCase(key), value)
		env = append(env, entry)
		Log(entry)
	}
	return env
}

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake  = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}