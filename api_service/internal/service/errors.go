package service

import "errors"

var (
	ErrParseUUID     = errors.New("Error parsing uuid value")
	ErrNotFound      = errors.New("Not found")
	ErrAlreadyExists = errors.New("Already exists")
)
