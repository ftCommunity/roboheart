package acm

import (
	"errors"
)

var (
	TokenNotFoundError = errors.New("Token not found")
	NotPermittedError  = errors.New("Operation not permitted")
)
