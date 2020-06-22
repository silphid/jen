package spec

type Root struct {
	Step
	Name string
	Description string
}

type Step struct {
	Value *ValueStep
	Option *OptionStep
	Multi *MultiStep
	Select *SelectStep
	Render string
	Do string
	Exec string
	If string
	Steps[] *Step
}

type ValueStep struct {
	Step
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
