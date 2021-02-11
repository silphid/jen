package evaluation

import (
	"path"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEntries(t *testing.T) {
	context := Context{
		Variables: VarMap{
			"VAR1":      "value1",
			"VAR2":      "value2",
			"TRUE_VAR":  "true",
			"EMPTY_VAR": "",
		},
		Placeholders: VarMap{
			"projekt": "myproject",
			"PROJEKT": "MYPROJECT",
		},
	}

	fixtures := []struct {
		Name     string
		Files    []string
		Expected []entry
		Error    string
	}{
		{
			Name: "plain names without rendering",
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
			Name: "plain names with rendering",
			Files: []string{
				"dir1/file2.txt.tmpl",
				"dir2/file3.txt",
				"file1.txt.tmpl",
			},
			Expected: []entry{
				{input: "dir1/file2.txt.tmpl", output: "dir1/file2.txt", render: true},
				{input: "dir2/file3.txt", output: "dir2/file3.txt", render: false},
				{input: "file1.txt.tmpl", output: "file1.txt", render: true},
			},
		},
		{
			Name: "the .tmpl extension enables rendering recursively",
			Files: []string{
				"dir1.tmpl/file1.txt",
				"dir1.tmpl/file2.txt",
				"dir1.tmpl/file3.txt.tmpl",
				"dir1.tmpl/dir/file1.txt",
				"dir1.tmpl/dir/file2.txt.tmpl",
				"dir2/file1.txt",
				"dir2/file2.txt.tmpl",
				"dir2/dir/file1.txt",
				"dir2/dir/file2.txt.tmpl",
			},
			Expected: []entry{
				{input: "dir1.tmpl/file1.txt", output: "dir1/file1.txt", render: true},
				{input: "dir1.tmpl/file2.txt", output: "dir1/file2.txt", render: true},
				{input: "dir1.tmpl/file3.txt.tmpl", output: "dir1/file3.txt", render: true},
				{input: "dir1.tmpl/dir/file1.txt", output: "dir1/dir/file1.txt", render: true},
				{input: "dir1.tmpl/dir/file2.txt.tmpl", output: "dir1/dir/file2.txt", render: true},
				{input: "dir2/file1.txt", output: "dir2/file1.txt", render: false},
				{input: "dir2/file2.txt.tmpl", output: "dir2/file2.txt", render: true},
				{input: "dir2/dir/file1.txt", output: "dir2/dir/file1.txt", render: false},
				{input: "dir2/dir/file2.txt.tmpl", output: "dir2/dir/file2.txt", render: true},
			},
		},
		{
			Name: "the .notmpl extension disables rendering recursively",
			Files: []string{
				"dir.tmpl/file.txt",
				"dir.tmpl/dir1/file.txt",
				"dir.tmpl/dir2.notmpl/file1.txt",
				"dir.tmpl/dir2.notmpl/file2.txt.tmpl",
				"dir.tmpl/dir2.notmpl/dir/file1.txt",
				"dir.tmpl/dir2.notmpl/dir/file2.txt.tmpl",
			},
			Expected: []entry{
				{input: "dir.tmpl/file.txt", output: "dir/file.txt", render: true},
				{input: "dir.tmpl/dir1/file.txt", output: "dir/dir1/file.txt", render: true},
				{input: "dir.tmpl/dir2.notmpl/file1.txt", output: "dir/dir2/file1.txt", render: false},
				{input: "dir.tmpl/dir2.notmpl/file2.txt.tmpl", output: "dir/dir2/file2.txt", render: true},
				{input: "dir.tmpl/dir2.notmpl/dir/file1.txt", output: "dir/dir2/dir/file1.txt", render: false},
				{input: "dir.tmpl/dir2.notmpl/dir/file2.txt.tmpl", output: "dir/dir2/dir/file2.txt", render: true},
			},
		},
		{
			Name: "conditional files",
			Files: []string{
				"dir1/file1[[.TRUE_VAR]].txt.tmpl",
				"dir1/file2[[.UNDEFINED_VAR]].txt.tmpl",
			},
			Expected: []entry{
				{input: "dir1/file1[[.TRUE_VAR]].txt.tmpl", output: "dir1/file1.txt", render: true},
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
				"dir1{{.VAR1}}/file1{{.VAR2}}.txt.tmpl",
			},
			Expected: []entry{
				{input: "dir1{{.VAR1}}/file1{{.VAR2}}.txt.tmpl", output: "dir1value1/file1value2.txt", render: true},
			},
		},
		{
			Name: "mixed variables and conditionals",
			Files: []string{
				"dir1{{.VAR1}}[[.TRUE_VAR]]/file1{{.VAR2}}[[.TRUE_VAR]].txt.tmpl",
			},
			Expected: []entry{
				{input: "dir1{{.VAR1}}[[.TRUE_VAR]]/file1{{.VAR2}}[[.TRUE_VAR]].txt.tmpl", output: "dir1value1/file1value2.txt", render: true},
			},
		},
		{
			Name: "invalid double-brace expression",
			Files: []string{
				"file1{{..}}.txt.tmpl",
			},
			Error: `failed to evaluate double-brace expression in name "file1{{..}}.txt.tmpl": parse template "file1{{..}}.txt.tmpl": template: base:1: unexpected <.> in operand`,
		},
		{
			Name: "replacements",
			Files: []string{
				"ABC_PROJEKT_DEF.txt",
				"abcprojektdef.txt",
			},
			Expected: []entry{
				{input: "ABC_PROJEKT_DEF.txt", output: "ABC_MYPROJECT_DEF.txt"},
				{input: "abcprojektdef.txt", output: "abcmyprojectdef.txt"},
			},
		},
		{
			Name: "empty folder names are collapsed in path",
			Files: []string{
				"dir1/[[.TRUE_VAR]]/dir2/file1.txt",
				"dir3/[[.UNDEFINED_VAR]]/dir4/file2.txt",
			},
			Expected: []entry{
				{input: "dir1/[[.TRUE_VAR]]/dir2/file1.txt", output: "dir1/dir2/file1.txt"},
			},
		},
	}

	getExpected := func(entries []entry, inputDir string) []entry {
		var results []entry
		for _, ent := range entries {
			results = append(results, entry{
				input:  path.Join(inputDir, ent.input),
				output: path.Join("/output", ent.output),
				render: ent.render,
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

			actual, err := getEntries(context, inputDir, outputDir, false)
			expected := getExpected(f.Expected, inputDir)

			sort.SliceStable(actual, func(i, j int) bool {
				return actual[i].input < actual[j].input
			})
			sort.SliceStable(expected, func(i, j int) bool {
				return expected[i].input < expected[j].input
			})

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
