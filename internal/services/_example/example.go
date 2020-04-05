package example

import "github.com/ftCommunity/roboheart/internal/service"

type example struct{}

type Example interface {
	MethodOne()
	MethodTwo(_ string, _ int) bool
}

func (e *example) Init(map[string]service.Service, func(string), func(error)) error { return nil }
func (e *example) Stop() error                                                      { return nil }
func (e *example) Name() string                                                     { return "example" }
func (e *example) Dependencies() ([]string, []string)                               { return []string{}, []string{} }
func (e *example) SetAdditionalDependencies(map[string]service.Service) error       { return nil }
func (e *example) UnsetAdditionalDependencies(s chan interface{})                   { s <- struct{}{} }

func (e *example) MethodOne()                     {}
func (e *example) MethodTwo(_ string, _ int) bool { return true }

var Service = new(example)
