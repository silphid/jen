package evaluation

import (
	"bytes"
	"fmt"
	"github.com/Masterminds/sprig"
	"github.com/Samasource/jen/internal/model"
	"github.com/Samasource/jen/internal/shell"
	"os/exec"
	"regexp"
	"strings"
	"text/template"
)

func EvalBoolExpression(values model.Values, expression string) (bool, error) {
	ifExpr := "{{if " + expression + "}}true{{end}}"
	result, err := EvalTemplate(values, ifExpr)
	if err != nil {
		return false, fmt.Errorf("evaluate expression %q: %w", expression, err)
	}
	return result == "true", nil
}

func EvalPromptValueTemplate(values model.Values, text string) (string, error) {
	// Interpolate go templating
	str, err := EvalTemplate(values, text)
	if err != nil {
		return "", err
	}

	if !strings.Contains(str, "$") {
		return str, nil
	}

	// Interpolates env vars
	buf := &bytes.Buffer{}
	cmd := &exec.Cmd{
		Path:   "/bin/bash",
		Args:   []string{"/bin/bash", "-c", `echo -n "` + str + `"`},
		Env:    shell.GetEnvFromValues(values.Variables),
		Stdout: buf,
	}
	if err = cmd.Run(); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func EvalTemplate(values model.Values, text string) (string, error) {
	// Perform replacement of placeholders
	for search, replace := range values.Placeholders {
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
var gotmplExtensionRegexp = regexp.MustCompile(`\.tmpl$`)

// evalFileName interpolates the double-brace expressions, evaluates and removes the conditionals in double-bracket
// expressions and returns the final file/dir name and whether it should be included in output and whether it should be
// rendered.
func evalFileName(values model.Values, name string) (string, bool, bool, error) {
	// Double-bracket expressions (ie: "[[.option]]") in names are evaluated to determine whether the file/folder should be
	// included in output and that expression then gets stripped from the name
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
			return "", false, false, fmt.Errorf("failed to eval double-bracket expression in name %q: %w", name, err)
		}

		// Should we exclude file/folder?
		if !value {
			return "", false, false, nil
		}

		// Remove expression from name
		name = name[:loc[0]] + name[loc[1]:]
	}

	// Double-brace expressions (ie: "{{.name}}") in names get interpolated as expected
	result, err := EvalTemplate(values, name)
	if err != nil {
		return "", false, false, fmt.Errorf("failed to evaluate double-brace expression in name %q: %w", name, err)
	}

	// Find .tmpl extension
	render := false
	if HasTmplExtension(result) {
		// Remove .tmpl extension
		result = result[:len(result)-len(".tmpl")]
		render = true
	}

	return result, true, render, nil
}

func HasTmplExtension(name string) bool {
	return gotmplExtensionRegexp.MatchString(name)
}
