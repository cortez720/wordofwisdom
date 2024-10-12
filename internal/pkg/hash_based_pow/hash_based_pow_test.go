package hashbasedpow

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	pow "wordOfWisdom/config/pow"
)

func TestNewPOW(t *testing.T) {
	tests := []struct {
		name     string
		cfg      *pow.PowConfig
		expected error
	}{
		{
			name:     "valid config",
			cfg:      &pow.PowConfig{Complexity: 20},
			expected: nil,
		},
		{
			name:     "nil config",
			cfg:      nil,
			expected: ErrInvalidComplexity,
		},
		{
			name:     "complexity too low",
			cfg:      &pow.PowConfig{Complexity: 0},
			expected: ErrInvalidComplexity,
		},
		{
			name:     "complexity too high",
			cfg:      &pow.PowConfig{Complexity: 25},
			expected: ErrInvalidComplexity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPOW(tt.cfg)
			require.Equal(t, tt.expected, err)
		})
	}
}

func TestChallenge(t *testing.T) {
	cfg := &pow.PowConfig{Complexity: 20}
	powInstance, _ := NewPOW(cfg)
	challenge := powInstance.Challenge()

	require.Equal(t, defaultTokenSize, len(challenge))
}

func TestVerify(t *testing.T) {
	cfg := &pow.PowConfig{Complexity: 20}
	powInstance, _ := NewPOW(cfg)

	challenge := powInstance.Challenge()
	solution := powInstance.Solve(challenge)

	err := powInstance.Verify(challenge, solution)
	require.NoError(t, err)

	err = powInstance.Verify([]byte("invalid"), solution)
	require.Equal(t, ErrInvalidChallenge, err)

	err = powInstance.Verify(challenge, []byte("invalid"))
	require.Equal(t, ErrInvalidSolution, err)

	err = powInstance.Verify(challenge, []byte(strings.Repeat(" ", defaultNonceSize)))
	require.Equal(t, ErrUnverified, err)
}

func TestSolve(t *testing.T) {
	cfg := &pow.PowConfig{Complexity: 20}
	powInstance, _ := NewPOW(cfg)

	challenge := powInstance.Challenge()
	solution := powInstance.Solve(challenge)

	require.NotNil(t, solution)

	err := powInstance.Verify(challenge, solution)
	require.NoError(t, err)
}
