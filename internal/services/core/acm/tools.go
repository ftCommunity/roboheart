package acm

import "errors"

var (
	NotPermittedError = errors.New("Operation not permitted")
)

func isUserError(err error) bool {
	return err == NotPermittedError || err == TokenNotFoundError
}
