package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrSessionNotFound = errors.New("session not found")
	ErrTaskNotFound = errors.New("task not found")
	ErrIternalError = errors.New("iternal Error")
	ErrWrongPassword = errors.New("wrong password")
	ErrMissingToken = errors.New("missing token")
	ErrUnauthorized = errors.New("unauthorized")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrSessionAlreadyExists = errors.New("session already exists")
	ErrTaskAlreadyExists = errors.New("task already exists")
)