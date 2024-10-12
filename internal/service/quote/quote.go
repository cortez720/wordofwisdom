package quote

import (
	"context"

	"golang.org/x/exp/rand"
)

var quotes = []string{
	"Do not dwell in the past, do not dream of the future, concentrate the mind on the present moment.",
	"The only way to do great work is to love what you do.",
	"Life is what happens when you're busy making other plans.",
	"Get busy living or get busy dying.",
}

// Service ...
type Service struct{}

// NewService ...
func NewService(_ context.Context) *Service {
	return &Service{}
}

// GetWordOfWisdom ...
func (svc *Service) GetWordOfWisdom(_ context.Context) (string, error) {
	return quotes[rand.Intn(len(quotes))], nil
}
