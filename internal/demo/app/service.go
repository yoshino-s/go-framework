package app

import (
	"context"

	"github.com/yoshino-s/go-framework/application"
)

type Service struct {
	*application.EmptyApplication
}

func NewService() *Service {
	return &Service{
		EmptyApplication: application.NewEmptyApplication("DemoService"),
	}
}

func (s *Service) Setup(ctx context.Context) {
	s.Logger.Info("DemoService setup completed")
}

func (s *Service) GetRandomNumber() int {
	// Simulate some processing
	s.Logger.Info("Generating a random number")
	return 42 // Placeholder for actual random number generation logic
}
