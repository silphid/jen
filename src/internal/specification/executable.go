package specification

type Executable interface {
	Execute(context Context) error
}

type Executables []Executable

func (executables Executables) Execute(context Context) error {
	for _, e := range executables {
		if err := e.Execute(context); err != nil {
			return err
		}
	}
	return nil
}
