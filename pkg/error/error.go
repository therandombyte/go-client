package error

type Error struct {
	Operation string
	User string
	Type string
	Err error
}

func (e *Error) Error() string {
	return e.Err.Error()
}



