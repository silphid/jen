package prompts

type Type string

const (
	Text    Type = "text"
	Option  Type = "option"
	Options Type = "options"
	Choice  Type = "choice"
)
