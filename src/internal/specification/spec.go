package specification

type Spec struct {
	Name        string
	Description string
	Version     string
	Actions     map[string]Action
}
