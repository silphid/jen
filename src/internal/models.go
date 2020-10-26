package internal

import "github.com/Samasource/jen/internal/specification"

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
	Actions     map[string][]*specification.Executable
}
