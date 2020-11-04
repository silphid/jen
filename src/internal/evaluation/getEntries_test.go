package evaluation

import (
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

func TestGetEntries(t *testing.T) {
	values := Values{
		Variables: map[string]interface{}{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  "true",
			"EMPTY_VAR": "",
		},
	}

	fixtures := []struct {
		Name     string
		Files    []string
		Expected []entry
		Error    string
	}{
		{
			Name: "plain names",
			Files: []string{
				"dir1/file2.txt",
				"dir2/file3.txt",
				"file1.txt",
			},
			Expected: []entry{
				{input: "dir1/file2.txt", output: "dir1/file2.txt"},
				{input: "dir2/file3.txt", output: "dir2/file3.txt"},
				{input: "file1.txt", output: "file1.txt"},
			},
		},
		{
			Name: "conditional files",
			Files: []string{
				"dir1/file1[[.TRUE_VAR]].txt",
				"dir1/file2[[.UNDEFINED_VAR]].txt",
			},
			Expected: []entry{
				{input: "dir1/file1[[.TRUE_VAR]].txt", output: "dir1/file1.txt"},
			},
		},
		{
			Name: "conditional dirs",
			Files: []string{
				"dir1[[.TRUE_VAR]]/file1.txt",
				"dir2[[.UNDEFINED_VAR]]/file2.txt",
			},
			Expected: []entry{
				{input: "dir1[[.TRUE_VAR]]/file1.txt", output: "dir1/file1.txt"},
			},
		},
		{
			Name: "variables",
			Files: []string{
				"dir1{{.VAR1}}/file1{{.VAR2}}.txt",
			},
			Expected: []entry{
				{input: "dir1{{.VAR1}}/file1{{.VAR2}}.txt", output: "dir1value1/file1value2.txt"},
			},
		},
		{
			Name: "mixed variables and conditionals",
			Files: []string{
				"dir1{{.VAR1}}[[.TRUE_VAR]]/file1{{.VAR2}}[[.TRUE_VAR]].txt",
			},
			Expected: []entry{
				{input: "dir1{{.VAR1}}[[.TRUE_VAR]]/file1{{.VAR2}}[[.TRUE_VAR]].txt", output: "dir1value1/file1value2.txt"},
			},
		},
		{
			Name: "invalid double-brace expression",
			Files: []string{
				"file1{{..}}.txt",
			},
			Error: `failed to parse double-brace expression in name "file1{{..}}.txt": template: base:1: unexpected <.> in operand`,
		},
	}

	getExpected := func(entries []entry, inputDir string) []entry {
		var results []entry
		for _, ent := range entries {
			results = append(results, entry{
				input:  path.Join(inputDir, ent.input),
				output: path.Join("/output", ent.output),
			})
		}
		return results
	}

	for _, f := range fixtures {
		t.Run(f.Name, func(t *testing.T) {

			inputDir := getTempDir()
			outputDir := "/output"
			defer removeAll(inputDir)

			for _, file := range f.Files {
				inputFile := path.Join(inputDir, file)
				createEmptyFile(inputFile)
			}

			actual, err := getEntries(values, inputDir, outputDir)
			expected := getExpected(f.Expected, inputDir)

			if f.Error != "" {
				assert.NotNil(t, err)
				assert.Equal(t, f.Error, err.Error())
			} else {
				assert.Nil(t, err)
				assert.Equal(t, expected, actual)
			}
		})
	}
}
