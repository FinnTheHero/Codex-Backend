package common

type Error struct {
	Err    error
	Status int
}

func (e *Error) Error() string {
	return e.Err.Error()
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) StatusCode() int {
	return e.Status
}
