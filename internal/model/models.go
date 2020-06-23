package model

type Values map[string]interface{}

type Context struct {
	Values Values
	Spec Spec
}

type Spec struct {
	Name        string
	Description string
	Actions     map[string][]StepUnion
	Steps       []*StepUnion
}

type StepUnion struct {
	Step
	Value  *ValueStep
	Option *OptionStep
	Multi  *MultiOptionStep
	Select *SelectStep
	Render string
	Do     string
	Exec   string
	If     string
}

type Step struct {
	Steps []*StepUnion
}

type ValueStep struct {
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
