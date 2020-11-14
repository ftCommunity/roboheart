package rherrors

import (
	"errors"
	"github.com/ftCommunity-roboheart/roboheart/package/instance"
)

func NewInstanceAssertionError(i instance.Instance) error {
	return errors.New("Failed to assert instance " + i.ID().String())
}
