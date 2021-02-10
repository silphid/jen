package model

// Config contains all information required to execute jen commands
type Config struct {
	TemplateName    string
	TemplateDir     string
	Spec            *Spec
	Values          Values
	BinDirs         []string
	SkipConfirm     bool
	VarOverrides    map[string]string
	OnValuesChanged func() error
}
