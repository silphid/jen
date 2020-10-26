package specification

import (
	"github.com/Samasource/jen/internal/specification/executable"
	"github.com/Samasource/jen/internal/specification/prompts"
	"github.com/Samasource/jen/internal/specification/prompts/text"
	"github.com/go-test/deep"
	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestLoadActions(t *testing.T) {
	tests := []struct {
		name     string
		buffer   string
		expected executable.Executable
	}{
		{
			name: "text prompt with default type",
			buffer: `
prompt:
  question: Question 1
  var: Variable 1
  default: Default Value 1`,
			expected: text.Prompt{
				Prompt: prompts.Prompt{
					Question: "Question 1",
				},
				Var:     "Variable 1",
				Default: "Default Value 1",
			},
		},
		{
			name: "text prompt with explicit type",
			buffer: `
prompt:
  type: text
  question: Question 1
  var: Variable 1
  default: Default Value 1`,
			expected: text.Prompt{
				Prompt: prompts.Prompt{
					Question: "Question 1",
				},
				Var:     "Variable 1",
				Default: "Default Value 1",
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
			assert.Nil(t, err)

			// Compare with expected actual
			if diff := deep.Equal(actual, tt.expected); diff != nil {
				t.Error(diff)
			}
		})
	}
}
