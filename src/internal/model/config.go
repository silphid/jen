package model

// Config contains all information required to execute jen commands
type Config struct {
	JenDir          string
	TemplatesDir    string
	TemplateName    string
	TemplateDir     string
	ProjectDir      string
	Spec            *Spec
	Values          Values
	BinDirs         []string
	SkipConfirm     bool
	RawVarOverrides []string
	VarOverrides    map[string]string
	OnValuesChanged func() error
}
