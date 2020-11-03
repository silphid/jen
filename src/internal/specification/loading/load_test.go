package loading

import (
	"github.com/Samasource/jen/internal/specification"
	"github.com/Samasource/jen/internal/specification/executable"
	"github.com/Samasource/jen/internal/specification/steps/choice"
	"github.com/Samasource/jen/internal/specification/steps/do"
	"github.com/Samasource/jen/internal/specification/steps/execute"
	"github.com/Samasource/jen/internal/specification/steps/input"
	"github.com/Samasource/jen/internal/specification/steps/option"
	"github.com/Samasource/jen/internal/specification/steps/options"
	"github.com/Samasource/jen/internal/specification/steps/render"
	"github.com/go-test/deep"
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type fixture struct {
	Name     string
	Buffer   string
	Expected interface{}
	Error    string
}

func run(t *testing.T, fixtures []fixture, load func(yaml.Map) (interface{}, error)) {
	for _, f := range fixtures {
		t.Run(f.Name, func(t *testing.T) {
			// Load Buffer into yaml dom
			reader := strings.NewReader(f.Buffer)
			file := new(yaml.File)
			var err error
			file.Root, err = yaml.Parse(reader)
			assert.Nil(t, err)

			// Load actual from yaml dom
			actual, err := load(file.Root.(yaml.Map))

			if f.Error != "" {
				// Ensure proper Error was returned
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), f.Error)
			} else {
				// Compare with Expected actual
				assert.Nil(t, err)
				if diff := deep.Equal(actual, f.Expected); diff != nil {
					t.Error(diff)
				}
			}
		})
	}
}

func TestLoadStep(t *testing.T) {
	fixtures := []fixture{
		{
			Name: "input prompt",
			Buffer: `
if: Condition
input:
  question: Question
  var: Variable
  default: Default`,
			Expected: input.Prompt{
				If:       "Condition",
				Question: "Question",
				Var:      "Variable",
				Default:  "Default",
			},
		},
		{
			Name: "input prompt without if or default",
			Buffer: `
input:
  question: Question
  var: Variable`,
			Expected: input.Prompt{
				Question: "Question",
				Var:      "Variable",
			},
		},
		{
			Name: "missing required question property",
			Buffer: `
input:
  var: Variable`,
			Error: `missing required property "question"`,
		},
		{
			Name: "missing required var property",
			Buffer: `
input:
  question: Question`,
			Error: `missing required property "var"`,
		},
		{
			Name: "option prompt",
			Buffer: `
option:
  question: Question
  var: Variable
  default: true`,
			Expected: option.Prompt{
				Question: "Question",
				Var:      "Variable",
				Default:  true,
			},
		},
		{
			Name: "option prompt with default default value",
			Buffer: `
option:
  question: Question
  var: Variable`,
			Expected: option.Prompt{
				Question: "Question",
				Var:      "Variable",
				Default:  false,
			},
		},
		{
			Name: "option prompt with invalid default value",
			Buffer: `
option:
  question: Question
  var: Variable
  default: Whatever`,
			Error: `invalid bool value: "Whatever"`,
		},
		{
			Name: "options prompt",
			Buffer: `
options:
  question: Question
  items:
    - text: Text 1
      var: Variable 1
      default: true
    - text: Text 2
      var: Variable 2
      default: false
    - text: Text 3
      var: Variable 3`,
			Expected: options.Prompt{
				Question: "Question",
				Items: []options.Item{
					{
						Text:    "Text 1",
						Var:     "Variable 1",
						Default: true,
					},
					{
						Text:    "Text 2",
						Var:     "Variable 2",
						Default: false,
					},
					{
						Text:    "Text 3",
						Var:     "Variable 3",
						Default: false,
					},
				},
			},
		},
		{
			Name: "choice prompt",
			Buffer: `
choice:
  question: Question
  var: Variable
  default: Default
  items:
    - text: Text 1
      value: Value 1
    - text: Text 2
      value: Value 2
    - text: Text 3
      value: Value 3`,
			Expected: choice.Prompt{
				Question: "Question",
				Var:      "Variable",
				Default:  "Default",
				Items: []choice.Item{
					{
						Text:  "Text 1",
						Value: "Value 1",
					},
					{
						Text:  "Text 2",
						Value: "Value 2",
					},
					{
						Text:  "Text 3",
						Value: "Value 3",
					},
				},
			},
		},
		{
			Name: "render step",
			Buffer: `
if: Condition
render:
  source: Source`,
			Expected: render.Render{
				If:     "Condition",
				Source: "Source",
			},
		},
		{
			Name: "exec step",
			Buffer: `
if: Condition
exec:
  command: Command`,
			Expected: execute.Execute{
				If:      "Condition",
				Command: "Command",
			},
		},
		{
			Name: "do step",
			Buffer: `
if: Condition
do:
  action: Action`,
			Expected: do.Do{
				If:     "Condition",
				Action: "Action",
			},
		},
	}

	run(t, fixtures, func(m yaml.Map) (interface{}, error) {
		return loadStep(m)
	})
}

func TestLoadActions(t *testing.T) {
	fixtures := []fixture{
		{
			Name: "two actions with one step each",
			Buffer: `
action1:
  - if: Condition 1
    input:
      question: Question 1
      var: Variable 1
action2:
  - input:
      question: Question 2
      var: Variable 2`,
			Expected: []specification.Action{
				{
					Name: "action1",
					Steps: []executable.Executable{
						input.Prompt{
							If:       "Condition 1",
							Question: "Question 1",
							Var:      "Variable 1",
						},
					},
				},
				{
					Name: "action2",
					Steps: []executable.Executable{
						input.Prompt{
							Question: "Question 2",
							Var:      "Variable 2",
						},
					},
				},
			},
		},
	}

	run(t, fixtures, func(m yaml.Map) (interface{}, error) {
		return loadActions(m)
	})
}

func TestLoadSpec(t *testing.T) {
	fixtures := []fixture{
		{
			Name: "",
			Buffer: `
metadata:
  Name: Name
  description: Description
  version: 0.0.1
import:
  common: common
  go: go/common
actions:
  action1:
    - if: Condition 1
      input:
        question: Question 1
        var: Variable 1
  action2:
    - input:
        question: Question 2
        var: Variable 2`,
			Expected: &specification.Spec{
				Name:        "Name",
				Description: "Description",
				Version:     "0.0.1",
				Actions: []specification.Action{
					{
						Name: "action1",
						Steps: []executable.Executable{
							input.Prompt{
								If:       "Condition 1",
								Question: "Question 1",
								Var:      "Variable 1",
							},
						},
					},
					{
						Name: "action2",
						Steps: []executable.Executable{
							input.Prompt{
								Question: "Question 2",
								Var:      "Variable 2",
							},
						},
					},
				},
			},
		},
	}

	run(t, fixtures, func(m yaml.Map) (interface{}, error) {
		return LoadSpec(m)
	})
}
