//go:generate mockgen -destination=./mocks.go -source=./services.go -package=handlers

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

// QuoteService ...
type QuoteService struct {
}

// NewService ...
func NewService(ctx context.Context) *QuoteService {
	return &QuoteService{}
}

// GetWordOfWisdom ...
func (_ *QuoteService) GetWordOfWisdom(ctx context.Context) (string, error) {
	return quotes[rand.Intn(len(quotes))], nil
}
