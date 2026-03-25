package models

import "errors"

var (
	ErrRecordNotFound     = errors.New("record not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidAuth        = errors.New("invalid authorization")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrInvalidTypeAssert  = errors.New("invalid type assert")
	ErrEmptyAssociations  = errors.New("no associations attached")
)
