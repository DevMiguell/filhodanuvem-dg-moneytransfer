package errors

type Error struct {
	Code    int
	Message string
	err     error
}

func (e Error) Error() string {
	return e.Message
}

func New(code int, message string, err error) Error {
	return Error{
		Code:    code,
		Message: message,
		err:     err,
	}
}
