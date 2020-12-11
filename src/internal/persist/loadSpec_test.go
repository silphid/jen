package persist

import (
	"github.com/Samasource/jen/internal/model"
	"github.com/Samasource/jen/internal/steps"
	"github.com/Samasource/jen/internal/steps/choice"
	"github.com/Samasource/jen/internal/steps/do"
	"github.com/Samasource/jen/internal/steps/exec"
	"github.com/Samasource/jen/internal/steps/input"
	"github.com/Samasource/jen/internal/steps/option"
	"github.com/Samasource/jen/internal/steps/options"
	"github.com/Samasource/jen/internal/steps/render"
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
				assert.NotNil(t, err)
				assert.Equal(t, f.Error, err.Error())
			} else {
				assert.NoError(t, err)
				if diff := deep.Equal(f.Expected, actual); diff != nil {
					t.Error(diff)
				}
			}
		})
	}
}

func TestLoadStep(t *testing.T) {
	fixtures := []fixture{
		{
			Name: "if",
			Buffer: `
if: Condition
then:
  - input:
      question: Message
      var: Variable
      default: Default`,
			Expected: steps.If{
				Condition: "Condition",
				Then: model.Executables{
					input.Prompt{
						Message: "Message",
						Var:     "Variable",
						Default: "Default",
					},
				},
			},
		},
		{
			Name: "input prompt",
			Buffer: `
input:
  question: Message
  var: Variable`,
			Expected: input.Prompt{
				Message: "Message",
				Var:     "Variable",
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
  question: Message`,
			Error: `missing required property "var"`,
		},
		{
			Name: "option prompt",
			Buffer: `
option:
  question: Message
  var: Variable
  default: true`,
			Expected: option.Prompt{
				Message: "Message",
				Var:     "Variable",
				Default: true,
			},
		},
		{
			Name: "option prompt with default default value",
			Buffer: `
option:
  question: Message
  var: Variable`,
			Expected: option.Prompt{
				Message: "Message",
				Var:     "Variable",
				Default: false,
			},
		},
		{
			Name: "option prompt with invalid default value",
			Buffer: `
option:
  question: Message
  var: Variable
  default: Whatever`,
			Error: `invalid bool value: "Whatever"`,
		},
		{
			Name: "options prompt",
			Buffer: `
options:
  question: Message
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
				Message: "Message",
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
  question: Message
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
				Message: "Message",
				Var:     "Variable",
				Default: "Default",
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
			Name: "render step long-hand",
			Buffer: `
render:
  source: Source`,
			Expected: render.Render{
				Source: "Source",
			},
		},
		{
			Name: "render step short-hand",
			Buffer: `
render: Source`,
			Expected: render.Render{
				Source: "Source",
			},
		},
		{
			Name: "exec step long-hand",
			Buffer: `
exec:
  commands:
    - Command 1
    - Command 2`,
			Expected: exec.Exec{
				Commands: []string{
					"Command 1",
					"Command 2",
				},
			},
		},
		{
			Name: "exec step short-hand",
			Buffer: `
exec: Command 1`,
			Expected: exec.Exec{
				Commands: []string{"Command 1"},
			},
		},
		{
			Name: "do step long-hand",
			Buffer: `
do:
  action: Action`,
			Expected: do.Do{
				Action: "Action",
			},
		},
		{
			Name: "do step short-hand",
			Buffer: `
do: Action`,
			Expected: do.Do{
				Action: "Action",
			},
		},
	}

	run(t, fixtures, func(m yaml.Map) (interface{}, error) {
		return loadExecutable(m)
	})
}

func TestLoadActions(t *testing.T) {
	fixtures := []fixture{
		{
			Name: "two actions with one step each",
			Buffer: `
action1:
  - if: Condition 1
    then:
      - input:
          question: Message 1
          var: Variable 1
action2:
  - input:
      question: Message 2
      var: Variable 2`,
			Expected: model.ActionMap{
				"action1": {
					Name: "action1",
					Steps: model.Executables{
						steps.If{
							Condition: "Condition 1",
							Then: model.Executables{
								input.Prompt{
									Message: "Message 1",
									Var:     "Variable 1",
								},
							},
						},
					},
				},
				"action2": {
					Name: "action2",
					Steps: model.Executables{
						input.Prompt{
							Message: "Message 2",
							Var:     "Variable 2",
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
  name: Name
  description: Description
  version: 0.2.0
import:
  common: common
  go: go/common
actions:
  action1:
    - if: Condition 1
      then:
        - input:
            question: Message 1
            var: Variable 1
  action2:
    - input:
        question: Message 2
        var: Variable 2`,
			Expected: &model.Spec{
				Name:        "Name",
				Description: "Description",
				Version:     "0.2.0",
				Actions: model.ActionMap{
					"action1": {
						Name: "action1",
						Steps: model.Executables{
							steps.If{
								Condition: "Condition 1",
								Then: model.Executables{
									input.Prompt{
										Message: "Message 1",
										Var:     "Variable 1",
									},
								},
							},
						},
					},
					"action2": {
						Name: "action2",
						Steps: model.Executables{
							input.Prompt{
								Message: "Message 2",
								Var:     "Variable 2",
							},
						},
					},
				},
			},
		},
	}

	run(t, fixtures, func(m yaml.Map) (interface{}, error) {
		return loadSpecFromMap(m)
	})
}
