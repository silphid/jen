package executable

type Executable interface {
	Execute(context Context) error
}

type Values map[string]interface{}

type Context struct {
	TemplateDir string
	OutputDir   string
	Values      Values
	Spec        Spec
}

type Spec struct {
	Name        string
	Description string
	Version     string
	Values      Values
	Actions     map[string][]*Executable
}
