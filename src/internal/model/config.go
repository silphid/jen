package model

type Config struct {
	JenDir       string
	TemplatesDir string
	TemplateName string
	TemplateDir  string
	ProjectDir   string
	Spec         *Spec
	Values       Values
}
