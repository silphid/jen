package evaluation

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/Samasource/jen/internal/specification"
	"regexp"
	"strings"
	"text/template"
)

func EvalBoolExpression(values specification.Values, expression string) (bool, error) {
	ifExpr := "{{if " + expression + "}}true{{end}}"
	result, err := EvalTemplate(values, ifExpr)
	if err != nil {
		return false, fmt.Errorf("evaluate expression %q: %w", expression, err)
	}
	return result == "true", nil
}

func EvalTemplate(values specification.Values, text string) (string, error) {
	// Perform replacements
	for search, replace := range values.Replacements {
		text = strings.ReplaceAll(text, search, replace)
	}

	// Render go template
	tmpl, err := template.New("base").Funcs(sprig.TxtFuncMap()).Parse(text)
	if err != nil {
		return "", fmt.Errorf("parse template %q: %w", text, err)
	}
	var buffer bytes.Buffer
	err = tmpl.Execute(&buffer, values.Variables)
	if err != nil {
		return "", fmt.Errorf("evaluate template %q: %w", text, err)
	}
	return buffer.String(), nil
}

var doubleBracketRegexp = regexp.MustCompile(`\[\[.*]]`)

// evalFileName interpolates the double-brace expressions, evaluates and removes the conditionals in double-bracket
// expressions and returns the final file/dir name and whether it should be included in template rendering.
func evalFileName(values specification.Values, name string) (string, bool, error) {
	// Double-bracket expressions (ie: "[[.option]]") in names are evaluated to determine
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

	// Double-brace expressions (ie: "{{.name}}") in names get interpolated as expected
	result, err := EvalTemplate(values, name)
	if err != nil {
		return "", false, fmt.Errorf("failed to evaluate double-brace expression in name %q: %w", name, err)
	}
	return result, true, nil
}
