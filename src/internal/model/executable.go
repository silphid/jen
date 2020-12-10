package model

type Executable interface {
	Execute(config *Config) error
}

type Executables []Executable

func (executables Executables) Execute(config *Config) error {
	for _, e := range executables {
		if err := e.Execute(config); err != nil {
			return err
		}
	}
	return nil
}
