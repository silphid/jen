package evaluation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type varMap = map[string]interface{}
type strMap = map[string]string

type context struct {
	vars         varMap
	placeholders strMap
}

func (c context) GetEvalVars() map[string]interface{} {
	return c.vars
}

func (c context) GetPlaceholders() map[string]string {
	return c.placeholders
}

func (c context) GetShellVars(includeProcessVars bool) []string {
	vars := make([]string, len(c.vars))
	for key, value := range c.vars {
		entry := fmt.Sprintf("%s=%v", key, value)
		vars = append(vars, entry)
	}
	return vars
}

func TestEvalBoolExpression(t *testing.T) {
	context := context{
		vars: varMap{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  true,
			"FALSE_VAR": false,
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
			Condition: `eq .TRUE_VAR true`,
			Expected:  true,
		},
		{
			Condition: `eq .FALSE_VAR false`,
			Expected:  true,
		},
		{
			Condition: `.VAR1`,
			Expected:  true,
		},
		{
			Condition: `.TRUE_VAR`,
			Expected:  true,
		},
		{
			Condition: `.FALSE_VAR`,
			Expected:  false,
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
	context := context{
		vars: varMap{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  true,
			"FALSE_VAR": false,
			"EMPTY_VAR": "",
		},
		placeholders: strMap{
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
			Name:            `Name with false [[ .FALSE_VAR ]]conditional`,
			ExpectedName:    ``,
			ExpectedInclude: false,
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
