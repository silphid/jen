package persist

import (
	"github.com/Samasource/jen/internal/model"
	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	// Save
	jenFile := model.JenFile{Variables: model.VarMap{
		"VAR1": "true",
		"VAR2": "abc",
	}}
	dir := getTempDir()
	err := SaveJenFileToDir(dir, jenFile)
	assert.NoError(t, err)

	// Load
	actualJenFile, err := LoadJenFileFromDir(dir)
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
