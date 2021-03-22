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
			return nil, fmt.Errorf("must specify either or both start/end regex")
		}
		sections[i] = section
	}
	return &Insert{
		sections: sections,
	}, nil
}

func (Insert) Eval(target string) (string, error) {
	return "", fmt.Errorf("not implemented")
}
