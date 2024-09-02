package gofemarterrors

type RequestError struct {
	InternalError error
	HTTPStatus    int
}

func (e *RequestError) Error() string {
	return e.InternalError.Error()
}

func (e *RequestError) Unwrap() error {
	return e.InternalError
}
