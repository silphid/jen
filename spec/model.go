package spec

type Root struct {
	Name string
	Description string
	Actions map[string][]StepUnion
	Steps[] *StepUnion
}

type StepUnion struct {
	Step
	Value *ValueStep
	Option *OptionStep
	Multi *MultiStep
	Select *SelectStep
	Render string
	Do string
	Exec string
	If string
}

type Step struct {
	Steps[] *StepUnion
}

type ValueStep struct {
	Step
	Value *ValueStep
	Name string
	Title string
}

type OptionStep struct {
	Step
	Name string
	Title string

}

type MultiStep struct {
	Step
	Title string
	Items[] *MultiItem
}

type MultiItem struct {
	Name string
	Title string
	Steps[] *Step // Optional steps that are only executed if this item gets selected
}

type SelectStep struct {
	Step
	Value string
	Title string
	Items[] *SelectItem
}

type SelectItem struct {
	Value string
	Title string
	Steps[] *Step // Optional steps that are only executed if this item gets selected
}
