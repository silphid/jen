package evaluation

import (
	"github.com/Samasource/jen/internal/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvalBoolExpression(t *testing.T) {
	values := model.Values{
		Variables: model.VarMap{
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
			actual, err := EvalBoolExpression(values, f.Condition)

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
	values := model.Values{
		Variables: model.VarMap{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  "true",
			"EMPTY_VAR": "",
		},
		Placeholders: model.VarMap{
			"projekt": "myproject",
			"PROJEKT": "MYPROJECT",
		},
	}

	fixtures := []struct {
		Name            string
		ExpectedInclude bool
		ExpectedName    string
		Error           string
	}{
		{
			Name:            `Name with true [[ .TRUE_VAR ]]conditional`,
			ExpectedInclude: true,
			ExpectedName:    `Name with true conditional`,
		},
		{
			Name:            `Name with false [[ .EMPTY_VAR ]]conditional`,
			ExpectedInclude: false,
			ExpectedName:    ``,
		},
		{
			Name:            `Name with variable {{ .VAR1 }}`,
			ExpectedInclude: true,
			ExpectedName:    `Name with variable value1`,
		},
		{
			Name:            `Plain name`,
			ExpectedInclude: true,
			ExpectedName:    `Plain name`,
		},
		{
			Name:            "abcprojektdef {{.VAR1}} ABC_PROJEKT_DEF",
			ExpectedInclude: true,
			ExpectedName:    "abcmyprojectdef value1 ABC_MYPROJECT_DEF",
		},
	}

	for _, f := range fixtures {
		t.Run(f.Name, func(t *testing.T) {
			actualName, actualInclude, err := evalFileName(values, f.Name)

			if f.Error != "" {
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), f.Error)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, f.ExpectedInclude, actualInclude)
				assert.Equal(t, f.ExpectedName, actualName)
			}
		})
	}
}

func TestEvalPromptValueTemplate(t *testing.T) {
	values := model.Values{
		Variables: model.VarMap{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  "true",
			"EMPTY_VAR": "",
		},
		Placeholders: model.VarMap{
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
	}

	for _, f := range fixtures {
		t.Run(f.Name, func(t *testing.T) {
			actual, err := EvalPromptValueTemplate(values, f.Value)
			assert.NoError(t, err)
			assert.Equal(t, f.Expected, actual)
		})
	}
}
