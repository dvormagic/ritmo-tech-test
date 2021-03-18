package errors

import "fmt"

type Error struct {
	Message string `json:"message"`
}

func NewError(message string) *Error {
	return &Error{
		Message: fmt.Sprintf("error: %s", message),
	}
}
