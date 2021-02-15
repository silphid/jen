package project

import (
	"io/ioutil"
	"testing"

	"github.com/go-test/deep"
	"github.com/stretchr/testify/assert"
)

type varMap = map[string]interface{}
type strMap = map[string]string

func TestSaveAndLoad(t *testing.T) {
	// Save
	proj := Project{Vars: varMap{
		"BOOL_VAR": true,
		"STR_VAR":  "abc",
	}}
	proj.Dir = getTempDir()
	err := proj.Save()
	assert.NoError(t, err)

	// Load
	actualProj, err := Load(proj.Dir)
	assert.NoError(t, err)

	// Compare
	if diff := deep.Equal(proj.Vars, actualProj.Vars); diff != nil {
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
