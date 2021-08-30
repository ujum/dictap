package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user doesn't exists")
	ErrUserAlreadyExists = errors.New("user already exists")
	//ErrUserDelete        = errors.New("cant delete user")
	//ErrUserUpdate        = errors.New("cant update user")
)
