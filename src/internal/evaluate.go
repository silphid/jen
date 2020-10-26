/*package internal

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"text/template"
)

func EvalBoolExpression(context Context, expression string) (bool, error) {
	ifExpr := "{{if " + expression + "}}true{{end}}"
	result, err := EvalTemplate(context, ifExpr)
	if err != nil {
		return false, fmt.Errorf("evaluate expression %q: %w", expression, err)
	}
	return result == "true", nil
}

func EvalTemplate(context Context, text string) (string, error) {
	tmpl, err := template.New("base").Funcs(sprig.TxtFuncMap()).Parse(text)
	if err != nil {
		return "", fmt.Errorf("parse template %q: %w", text, err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, context.Values)
	if err != nil {
		return "", fmt.Errorf("evaluate template %q: %w", text, err)
	}
	return buffer.String(), nil
}
*/
