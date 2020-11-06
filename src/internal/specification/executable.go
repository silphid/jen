package specification

type Executable interface {
	Execute(context Context) error
}
