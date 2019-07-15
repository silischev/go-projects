package main

import (
	"errors"
)

type httpError struct {
	OriginalError error
	Status        int
}

func (error httpError) Error() string {
	return error.OriginalError.Error()
}

func NewHttpError(msg string, status int) error {
	return httpError{errors.New(msg), status}
}
