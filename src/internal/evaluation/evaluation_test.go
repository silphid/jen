package evaluation

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalBoolExpression(t *testing.T) {
	context := Context{
		Variables: VarMap{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  "true",
			"FALSE_VAR": "false",
			"EMPTY_VAR": "",
		},
	}

	fixtures := []struct {
		Condition string
		Expected  bool
		Error     string
	}{
		{
			Condition: `eq "abc" "abc"`,
			Expected:  true,
		},
		{
			Condition: `eq "abc" "def"`,
			Expected:  false,
		},
		{
			Condition: `not (eq "abc" "def")`,
			Expected:  true,
		},
		{
			Condition: `eq .VAR1 "value1"`,
			Expected:  true,
		},
		{
			Condition: `not (eq .VAR1 "value2")`,
			Expected:  true,
		},
		{
			Condition: `not (eq .VAR1 .VAR2)`,
			Expected:  true,
		},
		{
			Condition: `true`,
			Expected:  true,
		},
		{
			Condition: `eq .TRUE_VAR "true"`,
			Expected:  true,
		},
		{
			Condition: `.TRUE_VAR`,
			Expected:  true,
		},
		{
			Condition: `.FALSE_VAR`,
			Expected:  true,
		},
		{
			Condition: `.EMPTY_VAR`,
			Expected:  false,
		},
		{
			Condition: `.UNDEFINED_VAR`,
			Expected:  false,
		},
		{
			Condition: `.VAR1`,
			Expected:  true,
		},
	}

	for _, f := range fixtures {
		t.Run(f.Condition, func(t *testing.T) {
			actual, err := EvalBoolExpression(context, f.Condition)

			if f.Error != "" {
				assert.NotNil(t, err)
				assert.Equal(t, f.Error, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, f.Expected, actual)
			}
		})
	}
}

func TestEvalFileName(t *testing.T) {
	context := Context{
		Variables: VarMap{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  "true",
			"EMPTY_VAR": "",
		},
		Placeholders: VarMap{
			"projekt": "myproject",
			"PROJEKT": "MYPROJECT",
		},
	}

	fixtures := []struct {
		Name            string
		ExpectedName    string
		ExpectedInclude bool
		ExpectedRender  RenderMode
		Error           string
	}{
		{
			Name:            `Name with true [[ .TRUE_VAR ]]conditional`,
			ExpectedName:    `Name with true conditional`,
			ExpectedInclude: true,
			ExpectedRender:  UnchangedRendering,
		},
		{
			Name:            `Name with false [[ .EMPTY_VAR ]]conditional`,
			ExpectedName:    ``,
			ExpectedInclude: false,
			ExpectedRender:  UnchangedRendering,
		},
		{
			Name:            `Name with variable {{ .VAR1 }}`,
			ExpectedName:    `Name with variable value1`,
			ExpectedInclude: true,
			ExpectedRender:  UnchangedRendering,
		},
		{
			Name:            `Plain name`,
			ExpectedName:    `Plain name`,
			ExpectedInclude: true,
			ExpectedRender:  UnchangedRendering,
		},
		{
			Name:            "abcprojektdef {{.VAR1}} ABC_PROJEKT_DEF",
			ExpectedName:    "abcmyprojectdef value1 ABC_MYPROJECT_DEF",
			ExpectedInclude: true,
			ExpectedRender:  UnchangedRendering,
		},
		{
			Name:            "Name of file to render.tmpl",
			ExpectedName:    "Name of file to render",
			ExpectedInclude: true,
			ExpectedRender:  EnableRendering,
		},
		{
			Name:            "Name of file NOT to render.notmpl",
			ExpectedName:    "Name of file NOT to render",
			ExpectedInclude: true,
			ExpectedRender:  DisableRendering,
		},
	}

	for _, f := range fixtures {
		t.Run(f.Name, func(t *testing.T) {
			actualName, actualInclude, actualRender, err := evalFileName(context, f.Name)

			if f.Error != "" {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), f.Error)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, f.ExpectedName, actualName)
				assert.Equal(t, f.ExpectedInclude, actualInclude)
				assert.Equal(t, f.ExpectedRender, actualRender)
			}
		})
	}
}

func TestEvalPromptValueTemplate(t *testing.T) {
	context := Context{
		Variables: VarMap{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  "true",
			"EMPTY_VAR": "",
		},
		Placeholders: VarMap{
			"projekt": "myproject",
			"PROJEKT": "MYPROJECT",
		},
	}

	fixtures := []struct {
		Name     string
		Value    string
		Expected string
	}{
		{
			Name:     "Plain text",
			Value:    `Hello World`,
			Expected: `Hello World`,
		},
		{
			Name:     "Go var",
			Value:    `Hello{{.VAR1}}World`,
			Expected: `Hellovalue1World`,
		},
		{
			Name:     "Braced shell var",
			Value:    `Hello${VAR1}World`,
			Expected: `Hellovalue1World`,
		},
		{
			Name:     "Unbraced shell var",
			Value:    `Hello $VAR1 World`,
			Expected: `Hello value1 World`,
		},
		{
			Name:     "Mixed go and shell vars",
			Value:    `Hello{{.VAR1}}World$VAR1`,
			Expected: `Hellovalue1Worldvalue1`,
		},
		{
			Name:     "Shell expression",
			Value:    `Hello$(echo -n Nice)World`,
			Expected: `HelloNiceWorld`,
		},
		{
			Name:     "Go var within shell expression",
			Value:    `Hello$(echo -n {{.VAR1}})World`,
			Expected: `Hellovalue1World`,
		},
		{
			Name:     "Braced shell var within go expression",
			Value:    `Hello{{"${VAR1}"}}World`,
			Expected: `Hellovalue1World`,
		},
		{
			Name:     "Unbraced shell var within go expression",
			Value:    `Hello {{"$VAR1"}} World`,
			Expected: `Hello value1 World`,
		},
		{
			Name:     "Ignore escaped triple braces",
			Value:    `Hello {{{ .VAR }}} World`,
			Expected: `Hello {{ .VAR }} World`,
		},
	}

	for _, f := range fixtures {
		t.Run(f.Name, func(t *testing.T) {
			actual, err := EvalPromptValueTemplate(context, f.Value)
			assert.NoError(t, err)
			assert.Equal(t, f.Expected, actual)
		})
	}
}
