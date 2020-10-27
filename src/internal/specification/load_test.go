package specification

import (
	"github.com/Samasource/jen/internal/specification/executable"
	"github.com/Samasource/jen/internal/specification/prompts/choice"
	"github.com/Samasource/jen/internal/specification/prompts/input"
	"github.com/Samasource/jen/internal/specification/prompts/option"
	"github.com/Samasource/jen/internal/specification/prompts/options"
	"github.com/go-test/deep"
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLoadStep(t *testing.T) {
	tests := []struct {
		name     string
		buffer   string
		expected executable.Executable
		error    string
	}{
		{
			name: "input prompt",
			buffer: `
if: Condition
input:
  question: Question
  var: Variable
  default: Default`,
			expected: input.Prompt{
				If:       "Condition",
				Question: "Question",
				Var:      "Variable",
				Default:  "Default",
			},
		},
		{
			name: "input prompt without if or default",
			buffer: `
input:
  question: Question
  var: Variable`,
			expected: input.Prompt{
				Question: "Question",
				Var:      "Variable",
			},
		},
		{
			name: "missing required question property",
			buffer: `
input:
  var: Variable`,
			error: `missing required property "question"`,
		},
		{
			name: "missing required var property",
			buffer: `
input:
  question: Question`,
			error: `missing required property "var"`,
		},
		{
			name: "option prompt",
			buffer: `
option:
  question: Question
  var: Variable
  default: true`,
			expected: option.Prompt{
				Question: "Question",
				Var:      "Variable",
				Default:  true,
			},
		},
		{
			name: "option prompt with default default value",
			buffer: `
option:
  question: Question
  var: Variable`,
			expected: option.Prompt{
				Question: "Question",
				Var:      "Variable",
				Default:  false,
			},
		},
		{
			name: "option prompt with invalid default value",
			buffer: `
option:
  question: Question
  var: Variable
  default: Whatever`,
			error: `invalid bool value: "Whatever"`,
		},
		{
			name: "options prompt",
			buffer: `
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
			expected: options.Prompt{
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
			name: "choice prompt",
			buffer: `
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
			expected: choice.Prompt{
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Load buffer into yaml dom
			reader := strings.NewReader(tt.buffer)
			f := new(yaml.File)
			var err error
			f.Root, err = yaml.Parse(reader)
			assert.Nil(t, err)

			// Load actual from yaml dom
			actual, err := loadStep(f.Root.(yaml.Map))

			if tt.error != "" {
				// Ensure proper error was returned
				assert.NotNil(t, err)
				assert.Equal(t, err.Error(), tt.error)
			} else {
				// Compare with expected actual
				assert.Nil(t, err)
				if diff := deep.Equal(actual, tt.expected); diff != nil {
					t.Error(diff)
				}
			}
		})
	}
}
