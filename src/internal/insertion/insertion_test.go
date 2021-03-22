package insertion

import (
	"testing"

	"github.com/go-test/deep"
	_assert "github.com/stretchr/testify/assert"
	_require "github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {

	items := []struct {
		name     string
		text     string
		expected *Insert
		error    string
	}{
		{
			name: "discard external text",
			text: `
discarded

<<< start
body
>>> end

discarded
`,
			expected: &Insert{
				sections: []Section{{
					start: "start",
					body:  "body",
					end:   "end",
				}},
			},
		},
		{
			name: "multi-line body",
			text: `
<<< start
this is
a multi-line
body
>>> end
`,
			expected: &Insert{
				sections: []Section{
					{
						start: "start",
						body: `this is
a multi-line
body`,
						end: "end",
					},
				},
			},
		},
		{
			name: "multi-section",
			text: `
<<< start1
body1
>>> end1
<<< start2
body2
>>> end2
`,
			expected: &Insert{
				sections: []Section{
					{
						start: "start1",
						body:  "body1",
						end:   "end1",
					},
					{
						start: "start2",
						body:  "body2",
						end:   "end2",
					},
				},
			},
		},
		{
			name: "no end regex",
			text: `
<<< start
body
>>>`,
			expected: &Insert{
				sections: []Section{
					{
						start: "start",
						body:  "body",
					},
				},
			},
		},
		{
			name: "no start regex",
			text: `
<<<
body
>>> end`,
			expected: &Insert{
				sections: []Section{
					{
						body: "body",
						end:  "end",
					},
				},
			},
		},
		{
			name: "neither start nor end regex is invalid",
			text: `
<<<
body
>>>
`,
			error: "must specify either or both start/end regex",
		},
		{
			name: "start flush on first line",
			text: `<<< start
body
>>> end
`,
			expected: &Insert{
				sections: []Section{
					{
						start: "start",
						body:  "body",
						end:   "end",
					},
				},
			},
		},
		{
			name: "end flush on last line with regex",
			text: `
<<< start
body
>>> end`,
			expected: &Insert{
				sections: []Section{
					{
						start: "start",
						body:  "body",
						end:   "end",
					},
				},
			},
		},
		{
			name: "end flush on last line without regex",
			text: `
<<< start
body
>>>`,
			expected: &Insert{
				sections: []Section{
					{
						start: "start",
						body:  "body",
					},
				},
			},
		},
	}

	for _, item := range items {
		t.Run(item.name, func(t *testing.T) {
			assert := _assert.New(t)
			oldCompareUnexportedFields := deep.CompareUnexportedFields
			deep.CompareUnexportedFields = true
			defer func() { deep.CompareUnexportedFields = oldCompareUnexportedFields }()

			actual, err := parse(item.text)

			if item.error != "" {
				assert.EqualError(err, item.error)
			} else {
				assert.NoError(err, "parse should complete successfully")
				assert.NotNil(actual)
				if actual != nil {
					assert.Equal(len(item.expected.sections), len(actual.sections), "number of sections")
					if diff := deep.Equal(item.expected, actual); diff != nil {
						t.Error(diff)
					}
				}
			}
		})
	}
}

func TestEval(t *testing.T) {

	items := []struct {
		name     string
		text     string
		insert   *Insert
		expected string
		error    string
	}{
		{
			name: "only start - should insert after first start match",
			text: `line 1
line 2
line 3
line 4
line 2`,
			insert: &Insert{
				sections: []Section{{
					start: "^line 2",
					body:  "body 1\nbody 2",
				}},
			},
			expected: `line 1
line 2
body 1
body 2
line 3
line 4
line 2`,
		},
		{
			name: "only end - should insert before first end match",
			text: `line 1
line 2
line 3
line 4`,
			insert: &Insert{
				sections: []Section{{
					body: "body 1\nbody 2",
					end:  "^line 4",
				}},
			},
			expected: `line 1
line 2
line 3
body 1
body 2
line 4`,
		},
		{
			name: "both start and end - should insert before first end match located after first start match",
			text: `line 1
line 4
line 2
line 3
line 4`,
			insert: &Insert{
				sections: []Section{{
					start: "^line 2",
					body:  "body 1\nbody 2",
					end:   "^line 4",
				}},
			},
			expected: `line 1
line 4
line 2
line 3
body 1
body 2
line 4`,
		},
		{
			name: "multiple sections",
			text: `line 1
line 2
line 3
line 4`,
			insert: &Insert{
				sections: []Section{
					{
						start: "^line 1",
						body:  "body 1\nbody 2",
					},
					{
						body: "body 3\nbody 4",
						end:  "^line 4",
					},
				},
			},
			expected: `line 1
body 1
body 2
line 2
line 3
body 3
body 4
line 4`,
		},
		{
			name: "no start match",
			text: `line 1
line 2
line 3
line 4`,
			insert: &Insert{
				sections: []Section{{
					start: "^line 5",
				}},
			},
			error: `could not locate insertion start "^line 5"`,
		},
		{
			name: "no end match",
			text: `line 1
line 2
line 3
line 4`,
			insert: &Insert{
				sections: []Section{{
					end: "^line 5",
				}},
			},
			error: `could not locate insertion end "^line 5"`,
		},
		{
			name: "no end match after given start",
			text: `line 1
line 2
line 3
line 4`,
			insert: &Insert{
				sections: []Section{{
					start: "^line 2",
					end:   "^line 1",
				}},
			},
			error: `could not locate insertion end "^line 1" after start "^line 2"`,
		},
	}

	for _, item := range items {
		t.Run(item.name, func(t *testing.T) {
			assert := _assert.New(t)
			require := _require.New(t)

			actual, err := item.insert.Eval(item.text)

			if item.error != "" {
				require.EqualError(err, item.error)
			} else {
				require.NoError(err, "parse should complete successfully")
				assert.Equal(item.expected, actual)
			}
		})
	}
}
