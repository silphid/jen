package evaluation

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"text/template"
)

type Values map[string]interface{}

func EvalBoolExpression(values Values, expression string) (bool, error) {
	ifExpr := "{{if " + expression + "}}true{{end}}"
	result, err := EvalTemplate(values, ifExpr)
	if err != nil {
		return false, fmt.Errorf("evaluate expression %q: %w", expression, err)
	}
	return result == "true", nil
}

func EvalTemplate(values Values, text string) (string, error) {
	tmpl, err := template.New("base").Funcs(sprig.TxtFuncMap()).Parse(text)
	if err != nil {
		return "", fmt.Errorf("parse template %q: %w", text, err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, values)
	if err != nil {
		return "", fmt.Errorf("evaluate template %q: %w", text, err)
	}
	return buffer.String(), nil
}
