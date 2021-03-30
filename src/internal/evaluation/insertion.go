package evaluation

import (
	"fmt"
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

var regex = regexp.MustCompile(`(?m)^<<< *(.*)\n((?:.*\n)*?(?:.*))\n>>> *(.*)$\n?`)

func NewInsert(text string) (*Insert, error) {
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

func (i Insert) Eval(context Context, text string) (string, error) {
	for _, section := range i.sections {
		// Determine where to insert section body into target string
		insertionIndex, err := findInsertionIndex(text, section.start, section.end)
		if err != nil {
			return "", err
		}

		// Evaluate section body as template
		body, err := EvalTemplate(context, string(section.body))
		if err != nil {
			return "", fmt.Errorf("failed to render insertion template body %q: %w", section.body, err)
		}

		// Insert section body at given insertion index
		text = text[:insertionIndex] + body + "\n" + text[insertionIndex:]
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
