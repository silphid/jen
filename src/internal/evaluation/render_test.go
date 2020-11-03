package evaluation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvalName(t *testing.T) {
	values := Values{
		"VAR1":      "value1",
		"VAR2":      "value2",
		"TRUE_VAR":  "true",
		"EMPTY_VAR": "",
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
	}

	for _, f := range fixtures {
		t.Run(f.Name, func(t *testing.T) {
			actualName, actualInclude, err := evalName(values, f.Name)

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
