package rherrors

import (
	"errors"
	"github.com/ftCommunity-roboheart/roboheart/package/instance"
)

func NewInstanceAssertionError(i interface{ ID() instance.ID }) error {
	return errors.New("Failed to assert instance " + i.ID().String())
}
