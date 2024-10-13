package hashbasedpow

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"

	config "github.com/cortez720/wordofwisdom/config/pow"
)

const (
	defaultTokenSize = 16
	defaultNonceSize = 8
	maxTargetBits    = 24
)

var (
	ErrInvalidChallenge  = errors.New("invalid challenge")
	ErrInvalidSolution   = errors.New("invalid solution")
	ErrUnverified        = errors.New("unverified challenge")
	ErrInvalidComplexity = errors.New("invalid complexity")
)

// POW represents a proof of work algorithm implementation based on hashcash.
type POW struct {
	complexity uint64
}

// NewPOW creates a new POW.
func NewPOW(cfg *config.PowConfig) (*POW, error) {
	if cfg == nil || cfg.Complexity < 1 || cfg.Complexity > maxTargetBits {
		return nil, ErrInvalidComplexity
	}

	return &POW{complexity: cfg.Complexity}, nil
}

// Challenge returns a challenge.
func (p *POW) Challenge() []byte {
	return newToken(p.complexity)
}

// Verify verifies a challenge and solution.
func (p *POW) Verify(challenge, solution []byte) error {
	if len(challenge) != defaultTokenSize {
		return ErrInvalidChallenge
	}

	if len(solution) != defaultNonceSize {
		return ErrInvalidSolution
	}

	if !verify(challenge, solution) {
		return ErrUnverified
	}

	return nil
}

// Solve solves a challenge.
func (p *POW) Solve(challenge []byte) []byte {
	if len(challenge) != defaultTokenSize {
		return nil
	}

	return solve(challenge)
}

func newToken(targetBits uint64) []byte {
	buf := make([]byte, defaultTokenSize)
	target := uint64(1) << (64 - targetBits)

	binary.BigEndian.PutUint64(buf[:8], target)
	_, _ = rand.Read(buf[8:])

	return buf
}

func hash(data, nonce []byte) []byte {
	h := sha256.New()
	h.Write(data)
	h.Write(nonce)

	return h.Sum(nil)
}

func verify(token, nonce []byte) bool {
	h := hash(token, nonce)
	return bytes.Compare(h, token) < 0
}

func solve(token []byte) []byte {
	nonce := make([]byte, defaultNonceSize)

	for i := uint64(0); i < math.MaxUint64; i++ {
		binary.BigEndian.PutUint64(nonce, i)

		if verify(token, nonce) {
			return nonce
		}
	}

	return nil
}
