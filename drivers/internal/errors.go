package internal

import "fmt"

func NewHTTPError(code int) error {
	return &HTTPError{
		code: code,
	}
}

type HTTPError struct {
	code int
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("http: response code %v", e.code)
}

func (e *HTTPError) Code() int {
	return e.code
}
