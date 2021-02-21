package exec

// Context encapsulates everything required by implementors
// of the Executable interface to perform their work
type Context interface {
	// GetVars returns a dictionary of the project's variable names mapped to
	// their corresponding values. It does not include the process' env var.
	// Whenever you alter this map, you are responsible for later calling
	// SetVars() to save your changes back to the project file.
	GetVars() map[string]interface{}

	// SetVars saves given variables in project file.
	SetVars(vars map[string]interface{}) error

	// IsVarOverriden returns whether given variable has been overriden via command
	// line. This is used to skip prompting for those variables.
	IsVarOverriden(name string) bool

	// GetPlaceholders returns a map of special placeholders that can be used instead
	// of go template expression, for lighter weight templating, especially for the
	// project's name, which appears everywhere.
	GetPlaceholders() map[string]string

	// GetEvalVars returns a dictionary of the project's variable names mapped to
	// their corresponding values for evaluation purposes. It does not include the
	// process' env var.
	GetEvalVars() map[string]interface{}

	// GetShellVars returns all env vars to be used when invoking shell commands,
	// including the current process' env vars, the project's vars and an augmented
	// PATH var including extra bin dirs.
	GetShellVars() []string

	// GetAction returns action with given name within same spec file or nil if not
	// found.
	GetAction(name string) Executable

	// GetProjectDir returns the current project's dir
	GetProjectDir() string
}

// Executable represents an entity that can perform some work
type Executable interface {
	Execute(context Context) error
}

// Executables represents a slice of multiple executables
type Executables []Executable

// Execute delegates the invocation to multiple child executables
func (executables Executables) Execute(context Context) error {
	for _, e := range executables {
		if err := e.Execute(context); err != nil {
			return err
		}
	}
	return nil
}
