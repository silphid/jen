package internal

import (
	"bytes"
	"fmt"
	"text/template"
)

func EvalExpression(context Context, expression string) (bool, error) {
	ifExpr := "{{if " + expression + "}}true{{end}}"
	tmpl, err := template.New("base").Parse(ifExpr)
	if err != nil {
		return false, fmt.Errorf("parse expression %q: %w", expression, err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, context.Values)
	if err != nil {
		return false, fmt.Errorf("evaluate expression %q: %w", expression, err)
	}
	return buffer.String() == "true", nil
}
