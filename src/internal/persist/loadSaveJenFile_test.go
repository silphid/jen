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
	proj := project.Project{Variables: model.VarMap{
		"VAR1": "true",
		"VAR2": "abc",
	}}
	dir := getTempDir()
	err := proj.Save(dir)
	assert.NoError(t, err)

	// Load
	actualProj, err := project.Load(dir)
	assert.NoError(t, err)

	// Compare
	if diff := deep.Equal(proj.Variables, actualProj.Variables); diff != nil {
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
