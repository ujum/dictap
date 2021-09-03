package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user doesn't exists")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserIncorrectPass = errors.New("incorrect password")
	ErrNotFound          = errors.New("not found")
	ErrAlreadyExists     = errors.New("already exists")
)
