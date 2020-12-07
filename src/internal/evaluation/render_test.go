package evaluation

//import (
//	"github.com/Samasource/jen/internal/specification"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestRender(t *testing.T) {
//	values := specification.Values{
//		Variables: map[string]interface{}{
//			"VAR1":      "value1",
//			"VAR2":      "value2",
//			"TRUE_VAR":  "true",
//			"EMPTY_VAR": "",
//		},
//		Replacements: map[string]string{
//			"projekt": "myproject",
//			"PROJEKT": "MYPROJECT",
//		},
//	}
//
//	inputFile := writeTempFile(f.Input)
//	outputFile := getTempFile()
//	defer deleteFile(inputFile)
//	defer deleteFile(outputFile)
//	err := Render(values, inputFile, outputFile)
//	actual := readFile(outputFile)
//
//	if f.Error != "" {
//		assert.NotNil(t, err)
//		assert.Equal(t, f.Error, err.Error())
//	} else {
//		assert.Nil(t, err)
//		assert.Equal(t, f.Expected, actual)
//	}
//}
