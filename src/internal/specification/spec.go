package specification

type Spec struct {
	Name        string
	Description string
	Version     string
	InputDir    string
	Actions     map[string]Action
}
