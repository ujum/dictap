package error

import "errors"

var (
	ErrUserIncorrectPass = errors.New("incorrect password")
	ErrNotFound          = errors.New("not found")
	ErrAlreadyExists     = errors.New("already exists")
)
