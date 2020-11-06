package specification

import "github.com/Samasource/jen/internal/evaluation"

type Context struct {
	Spec      Spec
	OutputDir string
	Values    evaluation.Values
}
