package trail

import "net/http"

type PublicError struct {
	Status  int
	Message string
}

func NewPublicError(Status int, message string) PublicError {
	return PublicError{
		Status:  Status,
		Message: message,
	}
}

type Error struct {
	PublicError
	error
}

func NewError(publicError PublicError, innerError error) Error {
	return Error{
		PublicError: publicError,
		error:       innerError,
	}
}

type Success struct {
	Status  int
	Headers http.Header
	Body    any
}
