package evaluation

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"regexp"
	"strings"
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

var doubleBracketRegexp = regexp.MustCompile(`\[\[.*]]`)

func evalFileName(values Values, name string) (string, bool, error) {
	// Double-bracket expression (ie: "[[.option]]") in names are evaluated to determine
	// whether the file/folder should be rendered and that expression then gets stripped
	// from the name
	for {
		// Find expression
		loc := doubleBracketRegexp.FindStringIndex(name)
		if loc == nil {
			break
		}
		exp := name[loc[0]+2 : loc[1]-2]

		// Evaluate expression
		value, err := EvalBoolExpression(values, exp)
		if err != nil {
			return "", false, fmt.Errorf("failed to eval double-bracket expression in name %q: %w", name, err)
		}

		// Should we exclude file/folder?
		if !value {
			return "", false, nil
		}

		// Remove expression from name
		name = name[:loc[0]] + name[loc[1]:]
	}

	// Double-brace expression (ie: "{{.name}}") in names get interpolated as expected
	if strings.Index(name, "{{") != -1 {
		tmpl, err := template.New("base").Parse(name)
		if err != nil {
			return "", false, fmt.Errorf("failed to parse double-brace expression in name %q: %w", name, err)
		}
		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, values)
		if err != nil {
			return "", false, fmt.Errorf("failed to render double-brace expression in name %q: %w", name, err)
		}
		return buffer.String(), true, nil
	}

	return name, true, nil
}
