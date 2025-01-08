package cmd

import (
	"github.com/go-test/deep"
	"testing"
)

type fixture struct {
	name                    string
	inputArgs, expectedArgs []string
}

func TestSplitArgs(t *testing.T) {
	fixtures := []fixture{
		{
			name: "no args",
		},
		{
			name:      "no flags",
			inputArgs: []string{"arg1", "arg2"},
		},
		{
			name:         "only flags",
			inputArgs:    []string{"--flag-a", "--flag-b", "--flag-c"},
			expectedArgs: []string{"--flag-a", "--flag-b", "--flag-c"},
		},
		{
			name:         "flags before args",
			inputArgs:    []string{"--flag-a", "--flag-b", "--flag-c", "arg1", "arg2"},
			expectedArgs: []string{"--flag-a", "--flag-b", "--flag-c"},
		},
		{
			name:      "flags after args",
			inputArgs: []string{"arg1", "arg2", "--flag-a", "--flag-b", "--flag-c"},
		},
		{
			name:         "flags and args mixed",
			inputArgs:    []string{"--flag-a", "arg1", "--flag-b", "arg2", "--flag-c"},
			expectedArgs: []string{"--flag-a"},
		},
	}

	for _, f := range fixtures {
		t.Run(f.name, func(t *testing.T) {
			flagArgs := getFlags(f.inputArgs)

			if diff := deep.Equal(flagArgs, f.expectedArgs); diff != nil {
				t.Error("Different expected flags", diff)
			}
		})
	}
}
