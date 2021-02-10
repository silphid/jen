package persist

import (
	"io/ioutil"
	"testing"

	"github.com/Samasource/jen/src/internal/model"
	"github.com/Samasource/jen/src/internal/project"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

func TestSaveAndLoad(t *testing.T) {
	// Save
	jenFile := project.JenFile{Variables: model.VarMap{
		"VAR1": "true",
		"VAR2": "abc",
	}}
	dir := getTempDir()
	err := jenFile.Save(dir)
	assert.NoError(t, err)

	// Load
	actualJenFile, err := project.Load(dir)
	assert.NoError(t, err)

	// Compare
	if diff := deep.Equal(jenFile.Variables, actualJenFile.Variables); diff != nil {
		t.Error(diff)
	}
}

func getTempDir() string {
	dir, err := ioutil.TempDir("/tmp", "jen_test_")
	if err != nil {
		panic(err)
	}
	return dir
}
