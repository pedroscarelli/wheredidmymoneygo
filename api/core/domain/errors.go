package domain

import "errors"

var (
	ErrAccountNotFound    = errors.New("account not found")
	ErrDuplicateAccount   = errors.New("an account with this e-mail or CPF already exists")
	ErrInvalidCredentials = errors.New("incorrect email or password")
)
