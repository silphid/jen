package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Samasource/jen/src/internal/exec"
	"github.com/Samasource/jen/src/internal/helpers"
	"github.com/Samasource/jen/src/internal/home"
	"github.com/Samasource/jen/src/internal/logging"
	"github.com/Samasource/jen/src/internal/project"
	"github.com/Samasource/jen/src/internal/spec"
)

// Options represents all command line configurations
type Options struct {
	TemplateName string
	SkipConfirm  bool
	VarOverrides []string
}

// NewContext creates a context to be used for executing executables
func (o Options) NewContext() (exec.Context, error) {
	_, err := home.GetOrCloneJenRepo()
	if err != nil {
		return nil, err
	}

	proj, err := project.LoadOrCreate(o.TemplateName, o.SkipConfirm, o.VarOverrides)
	if err != nil {
		return nil, err
	}

	templateDir, err := proj.GetTemplateDir()
	if err != nil {
		return nil, err
	}

	specification, err := spec.Load(templateDir)
	if err != nil {
		return nil, err
	}

	jenHomeDir, err := home.GetJenHomeDir()
	if err != nil {
		return nil, err
	}

	return context{
		jenHomeDir:  jenHomeDir,
		templateDir: templateDir,
		project:     *proj,
		spec:        *specification,
	}, nil
}

// context contains all the information for implementing both the
// exec.context and evaluation.context interfaces
type context struct {
	jenHomeDir  string
	templateDir string
	project     project.Project
	spec        spec.Spec
}

// GetVars returns a dictionary of the project's variable names mapped to
// their corresponding values. It does not include the process' env var.
// Whenever you alter this map, you are responsible for later calling SaveProject().
func (c context) GetVars() map[string]interface{} {
	return c.project.Vars
}

// SaveProject saves all of the project's variables to project file.
func (c context) SaveProject() error {
	return c.project.Save()
}

// IsVarOverriden returns whether given variable has been overriden via command
// line. This is used to skip prompting for those variables.
func (c context) IsVarOverriden(name string) bool {
	for _, x := range c.project.OverridenVars {
		if x == name {
			return true
		}
	}
	return false
}

// GetPlaceholders returns a map of special placeholders that can be used instead
// of go template expression, for lighter weight templating, especially for the
// project's name, which appears everywhere. For now, the only supported placeholder
// is PROJECT, but we will eventually make placeholders configurable in spec file.
func (c context) GetPlaceholders() map[string]string {
	value, _ := c.project.Vars["PROJECT"]
	str, ok := value.(string)
	if !ok {
		return nil
	}
	return map[string]string{
		"projekt": strings.ToLower(str),
		"PROJEKT": strings.ToUpper(str),
	}
}

// GetShellVars returns all env vars to be used when invoking shell commands,
// including the current process' env vars, the project's vars and an augmented
// PATH var including extra bin dirs.
func (c context) GetShellVars() []string {
	binDirs := []string{
		filepath.Join(c.jenHomeDir, "bin"),
		filepath.Join(c.project.Dir, "bin"),
	}

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
	for key, value := range c.project.Vars {
		entry := fmt.Sprintf("%s=%v", key, value)
		env = append(env, entry)
		logging.Log(entry)
	}
	return env
}

// GetAction returns action with given name within same
// spec file or nil if not found.
func (c context) GetAction(name string) exec.Executable {
	action, ok := c.spec.Actions[name]
	if !ok {
		return nil
	}
	return action
}

// GetProjectDir returns the current project's dir
func (c context) GetProjectDir() string {
	return c.project.Dir
}
