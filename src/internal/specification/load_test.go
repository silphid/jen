package specification

import (
	"github.com/Samasource/jen/internal/specification/executable"
	"github.com/Samasource/jen/internal/specification/prompts/option"
	"github.com/Samasource/jen/internal/specification/prompts/text"
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
			name: "text prompt with default type",
			buffer: `
if: Condition
prompt:
  question: Question
  var: Variable
  default: Default Value`,
			expected: text.Prompt{
				If:       "Condition",
				Question: "Question",
				Var:      "Variable",
				Default:  "Default Value",
			},
		},
		{
			name: "text prompt with explicit type",
			buffer: `
if: Condition
prompt:
  type: text
  question: Question
  var: Variable
  default: Default Value`,
			expected: text.Prompt{
				If:       "Condition",
				Question: "Question",
				Var:      "Variable",
				Default:  "Default Value",
			},
		},
		{
			name: "text prompt without if or default",
			buffer: `
prompt:
  question: Question
  var: Variable`,
			expected: text.Prompt{
				Question: "Question",
				Var:      "Variable",
			},
		},
		{
			name: "missing required question property",
			buffer: `
prompt:
  var: Variable`,
			error: `missing required property "question"`,
		},
		{
			name: "missing required var property",
			buffer: `
prompt:
  question: Question`,
			error: `missing required property "var"`,
		},
		{
			name: "option prompt",
			buffer: `
prompt:
  type: option
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
prompt:
  type: option
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
prompt:
  type: option
  question: Question
  var: Variable
  default: Whatever`,
			error: `invalid bool value: "Whatever"`,
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
