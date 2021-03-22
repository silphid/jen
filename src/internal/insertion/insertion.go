package insertion

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

type Section struct {
	start string
	end   string
	body  string
}

type Insert struct {
	sections []Section
}

func Load(file string) (*Insert, error) {
	text, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return parse(string(text))
}

var regex = regexp.MustCompile(`(?m)^<<< *(.*)\n((?:.*\n)*?(?:.*))\n>>> *(.*)$\n?`)

func parse(text string) (*Insert, error) {
	matches := regex.FindAllStringSubmatch(text, -1)
	sections := make([]Section, len(matches))
	for i, match := range matches {
		section := Section{
			start: match[1],
			body:  match[2],
			end:   match[3],
		}
		if section.start == "" && section.end == "" {
			return nil, fmt.Errorf("cannot omit both start and end regex")
		}
		sections[i] = section
	}
	return &Insert{
		sections: sections,
	}, nil
}

func (i Insert) Eval(text string) (string, error) {
	for _, section := range i.sections {
		insertionIndex, err := findInsertionIndex(text, section.start, section.end)
		if err != nil {
			return "", err
		}
		text = text[:insertionIndex] + section.body + "\n" + text[insertionIndex:]
	}
	return text, nil
}

func findInsertionIndex(text, start, end string) (int, error) {
	insertionIndex := 0

	// Find start
	if start != "" {
		startRegex, err := regexp.Compile("(?m)" + start + `.*\n`)
		if err != nil {
			return -1, fmt.Errorf("invalid start regex %q: %w", start, err)
		}
		index := startRegex.FindStringIndex(text)
		if index == nil {
			return -1, fmt.Errorf("could not locate insertion start %q", start)
		}
		insertionIndex = index[1]
	}

	// Find end
	if end != "" {
		endRegex, err := regexp.Compile("(?m)" + end)
		if err != nil {
			return -1, fmt.Errorf("invalid end regex %q: %w", end, err)
		}
		index := endRegex.FindStringIndex(text[insertionIndex:])
		if index == nil {
			if start != "" {
				return -1, fmt.Errorf("could not locate insertion end %q after start %q", end, start)
			} else {
				return -1, fmt.Errorf("could not locate insertion end %q", end)
			}
		}
		insertionIndex += index[0]
	}

	return insertionIndex, nil
}
