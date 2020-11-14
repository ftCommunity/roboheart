package rherrors

import "errors"

var (
	DataNotReadyError = errors.New("data is not ready")
)
