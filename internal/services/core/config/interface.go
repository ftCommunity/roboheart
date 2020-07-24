package config

import (
	"github.com/ftCommunity/roboheart/internal/service"
)

type Config interface {
	GetServiceConfig(s service.Service) *ServiceConfig
}
