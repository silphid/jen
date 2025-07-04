package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/silphid/jen/cmd/jen/internal/shell"

	"github.com/silphid/jen/cmd/jen/internal/exec"
	"github.com/silphid/jen/cmd/jen/internal/helpers"
	"github.com/silphid/jen/cmd/jen/internal/home"
	"github.com/silphid/jen/cmd/jen/internal/logging"
	"github.com/silphid/jen/cmd/jen/internal/project"
	"github.com/silphid/jen/cmd/jen/internal/spec"
)

// Options represents all command line configurations
type Options struct {
	TemplateName string
	SkipConfirm  bool
	SkipPull     bool
	VarOverrides []string
}

// NewContext creates a context to be used for executing executables
func (o Options) NewContext() (exec.Context, error) {
	_, err := home.GetOrCloneRepo()
	if err != nil {
		return nil, err
	}

	if !o.SkipPull {
		if err = maybePull(); err != nil {
			return nil, err
		}
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

	cloneSubDir, err := home.GetCloneSubDir()
	if err != nil {
		return nil, err
	}

	return context{
		cloneSubDir: cloneSubDir,
		templateDir: templateDir,
		project:     proj,
		spec:        *specification,
	}, nil
}

func maybePull() error {
	jenHome, err := home.GetOrCloneRepo()
	if err != nil {
		return err
	}

	sentinelFile := filepath.Join(jenHome, ".git/last-pull")

	if info, err := os.Stat(sentinelFile); err == nil {
		if time.Since(info.ModTime()).Hours() < 24 {
			logging.Log("Skipping git pull as it was run less than 24 hours ago. Use `jen pull` to force a pull.")
			return nil
		}
	}

	logging.Log("Running git pull in %s ...\n", jenHome)
	if err := shell.ExecuteOutputOnlyErrors(nil, jenHome, "git pull"); err != nil {
		return fmt.Errorf("pulling latest templates: %w\nuse --skip-pull flag to bypass pulling template git repo", err)
	}

	if err := os.WriteFile(sentinelFile, []byte{}, 0644); err != nil {
		return fmt.Errorf("Failed to update last pull sentinel file: %v", err)
	}

	return nil
}

// context contains all the information for implementing both the
// exec.context and evaluation.context interfaces
type context struct {
	cloneSubDir string
	templateDir string
	project     *project.Project
	spec        spec.Spec
}

// GetVars returns a dictionary of the project's variable names mapped to
// their corresponding values. It does not include the process' env var.
// Whenever you alter this map, you are responsible for later calling
// SetVars() to save your changes back to the project file.
func (c context) GetVars() map[string]interface{} {
	clone := make(map[string]interface{})
	for k, v := range c.project.Vars {
		clone[k] = v
	}
	return clone
}

// SetVars saves given variables in project file.
func (c context) SetVars(vars map[string]interface{}) error {
	c.project.Vars = vars
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
// of go template expressions, for more lightweight templating, especially for the
// project's name, which appears everywhere.
func (c context) GetPlaceholders() map[string]string {
	return c.spec.Placeholders
}

// GetEvalVars returns a dictionary of the project's variable names mapped to
// their corresponding values for evaluation purposes. It does not include the
// process' env var.
func (c context) GetEvalVars() map[string]interface{} {
	// Combine persistent and transient vars
	vars := make(map[string]interface{})
	for k, v := range c.project.Vars {
		if !strings.HasPrefix(k, "~") {
			vars[k] = v
		}
	}
	for k, v := range c.project.Vars {
		if strings.HasPrefix(k, "~") {
			k = strings.TrimPrefix(k, "~")
			vars[k] = v
		}
	}

	// Compute built-in vars
	absProjectDir, err := filepath.Abs(c.project.Dir)
	if err != nil {
		panic(fmt.Errorf("failed to determine project's absolute dir: %w", err))
	}
	vars["PROJECT_DIR"] = absProjectDir
	vars["PROJECT_DIR_NAME"] = filepath.Base(absProjectDir)
	return vars
}

// getBinDirs returns the list of bin dirs that actually exist
func (c context) getBinDirs() []string {
	binDirs := []string{
		filepath.Join(c.project.Dir, "bin"),
		filepath.Join(c.templateDir, "bin"),
		filepath.Join(c.cloneSubDir, "bin"),
	}

	// Add bin dirs to PATH env var
	var validBinDirs []string
	for _, dir := range binDirs {
		if helpers.PathExists(dir) {
			validBinDirs = append(validBinDirs, dir)
		}
	}
	return validBinDirs
}

// GetScripts returns the list of executable scripts in bin dirs
func (c context) GetScripts() ([]string, error) {
	var scripts []string
	binDirs := c.getBinDirs()
	for _, dir := range binDirs {
		infos, err := ioutil.ReadDir(dir)
		if err != nil {
			return nil, err
		}
		for _, info := range infos {
			// Is it executable by owner, group or any?
			if info.Mode()&0111 != 0 {
				scripts = append(scripts, info.Name())
			}
		}
	}
	return scripts, nil
}

// GetShellVars returns all env vars to be used when invoking shell commands,
// including the current process' env vars, the project's vars and an augmented
// PATH var including extra bin dirs.
func (c context) GetShellVars(includeProcessVars bool) []string {
	// Combine bin dirs with PATH env var
	pathVar := strings.Join(append(c.getBinDirs(), os.Getenv("PATH")), ":")

	// Collect all current process env vars, except PATH
	var env []string
	if includeProcessVars {
		for _, entry := range os.Environ() {
			if !strings.HasPrefix(entry, "PATH=") {
				env = append(env, entry)
			}
		}
	}

	// Override PATH env var
	entry := fmt.Sprintf("PATH=%v", pathVar)
	env = append(env, entry)
	logging.Log(entry)

	// Then values env vars
	logging.Log("Environment variables:")
	vars := c.GetEvalVars()
	for key, value := range vars {
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

// GetActionNames returns the names of all actions available in template.
func (c context) GetActionNames() []string {
	names := make([]string, 0, len(c.spec.Actions))
	for name := range c.spec.Actions {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetTemplateDir returns the current template's dir
func (c context) GetTemplateDir() string {
	return c.templateDir
}

// GetProjectDir returns the current project's dir
func (c context) GetProjectDir() string {
	return c.project.Dir
}
