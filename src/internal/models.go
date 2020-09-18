package internal

type Values map[string]interface{}

type Context struct {
	TemplateDir string
	OutputDir string
	Values      Values
	Spec        Spec
}

type Spec struct {
	Name        string
	Description string
	Steps       []*StepUnion
	Actions     map[string][]*StepUnion
}

type StepUnion struct {
	Step
	String    *StringStep
	Secret    *SecretStep
	Option    *OptionStep
	Multi     *MultiOptionStep
	Select    *SelectStep
	SetOutput string `yaml:"setOutput"`
	Render    string
	Do        string
	Exec      string
	If        string
}

type Step struct {
	Steps []*StepUnion
}

type StringStep struct {
	Step
	Name  string
	Title string
}

type SecretStep struct {
	Step
	Name  string
	Title string
}

type OptionStep struct {
	Step
	Name  string
	Title string
}

type MultiOptionStep struct {
	Step
	Title string
	Items []*MultiOptionItem
}

type MultiOptionItem struct {
	Name string
	Title string
	Steps []*Step // Optional steps that are only executed if this item gets selected
}

type SelectStep struct {
	Step
	Name  string
	Title string
	Items []*SelectItem
}

type SelectItem struct {
	Value string
	Title string
	Steps []*Step // Optional steps that are only executed if this item gets selected
}
