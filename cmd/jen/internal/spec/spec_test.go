package spec

import (
	"strings"
	"testing"

	"github.com/go-test/deep"
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/silphid/jen/cmd/jen/internal/exec"
	"github.com/silphid/jen/cmd/jen/internal/steps"
	"github.com/silphid/jen/cmd/jen/internal/steps/choice"
	"github.com/silphid/jen/cmd/jen/internal/steps/do"
	execstep "github.com/silphid/jen/cmd/jen/internal/steps/exec"
	"github.com/silphid/jen/cmd/jen/internal/steps/input"
	"github.com/silphid/jen/cmd/jen/internal/steps/option"
	"github.com/silphid/jen/cmd/jen/internal/steps/options"
	"github.com/silphid/jen/cmd/jen/internal/steps/render"
	"github.com/silphid/jen/cmd/jen/internal/steps/set"
	"github.com/stretchr/testify/assert"
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
				Then: exec.Executables{
					input.Prompt{
						Message: "Message",
						Var:     "Variable",
						Default: "Default",
					},
				},
			},
		},
		{
			Name: "confirm",
			Buffer: `
confirm: Message
then:
  - input:
      question: Message
      var: Variable
      default: Default`,
			Expected: steps.Confirm{
				Message: "Message",
				Then: exec.Executables{
					input.Prompt{
						Message: "Message",
						Var:     "Variable",
						Default: "Default",
					},
				},
			},
		},
		{
			Name: "set step",
			Buffer: `
set:
  VAR1: Value1
  VAR2: Value2`,
			Expected: set.Set{
				Variables: []set.Variable{
					{
						Name:  "VAR1",
						Value: "Value1",
					},
					{
						Name:  "VAR2",
						Value: "Value2",
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
    - value: Value 3`,
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
						Text:  "Value 3",
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
				InputDir: "Source",
			},
		},
		{
			Name: "render step short-hand",
			Buffer: `
render: Source`,
			Expected: render.Render{
				InputDir: "Source",
			},
		},
		{
			Name: "exec step long-hand",
			Buffer: `
exec:
  commands:
    - Command 1
    - Command 2`,
			Expected: execstep.Exec{
				Commands: []string{
					"Command 1",
					"Command 2",
				},
			},
		},
		{
			Name: "exec step multiple child strings",
			Buffer: `
exec:
  - Command 1
  - Command 2`,
			Expected: execstep.Exec{
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
			Expected: execstep.Exec{
				Commands: []string{"Command 1"},
			},
		},
		{
			Name: "do step long-hand",
			Buffer: `
do:
  actions:
    - Action 1
    - Action 2`,
			Expected: do.Do{
				Actions: []string{"Action 1", "Action 2"},
			},
		},
		{
			Name: "do step multiple child strings",
			Buffer: `
do:
  - Action 1
  - Action 2`,
			Expected: do.Do{
				Actions: []string{"Action 1", "Action 2"},
			},
		},
		{
			Name: "do step short-hand",
			Buffer: `
do: Action`,
			Expected: do.Do{
				Actions: []string{"Action"},
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
			Expected: ActionMap{
				"action1": Action{
					Name: "action1",
					Steps: exec.Executables{
						steps.If{
							Condition: "Condition 1",
							Then: exec.Executables{
								input.Prompt{
									Message: "Message 1",
									Var:     "Variable 1",
								},
							},
						},
					},
				},
				"action2": Action{
					Name: "action2",
					Steps: exec.Executables{
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
version: 2021.04
description: Description
placeholders:
  projekt: {{.PROJECT | lower}}
  Projekt: {{.PROJECT | title}}
  PROJEKT: {{.PROJECT | upper}}
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
			Expected: &Spec{
				Name:        "template_name",
				Version:     "2021.04",
				Description: "Description",
				Placeholders: map[string]string{
					"projekt": "{{.PROJECT | lower}}",
					"Projekt": "{{.PROJECT | title}}",
					"PROJEKT": "{{.PROJECT | upper}}",
				},
				Actions: ActionMap{
					"action1": Action{
						Name: "action1",
						Steps: exec.Executables{
							steps.If{
								Condition: "Condition 1",
								Then: exec.Executables{
									input.Prompt{
										Message: "Message 1",
										Var:     "Variable 1",
									},
								},
							},
						},
					},
					"action2": Action{
						Name: "action2",
						Steps: exec.Executables{
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
		return loadFromMap(m, "path/to/template_name")
	})
}
