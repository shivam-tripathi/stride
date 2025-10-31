package service

import (
	"quizizz.com/internal/config"
)

// AppService defines the interface for application business logic
type AppService interface {
	GetPingMessage() string
}

// appService implements the AppService interface
type appService struct {
	config *config.Config
}

// NewAppService creates a new AppService
func NewAppService(config *config.Config) AppService {
	return &appService{
		config: config,
	}
}

// GetPingMessage returns a ping message
func (s *appService) GetPingMessage() string {
	return "pong from " + s.config.AppName
}
