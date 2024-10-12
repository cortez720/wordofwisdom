//go:generate mockgen -destination=./mocks/challenge_handler_mock.go -source=./challenge_handler.go -package=handler

package handler

import (
	"context"
	"fmt"
	"net/http"
)

const (
	challengeArg = "challenge"
	solutionArg  = "solution"
)

type quoteService interface {
	GetWordOfWisdom(ctx context.Context) (string, error)
}

type pow interface {
	Challenge() []byte
	Verify(challenge, solution []byte) error
}

type PowHandler struct {
	pow pow
	svc quoteService
}

func NewPow(pow pow, svc quoteService) *PowHandler {
	return &PowHandler{pow: pow, svc: svc}
}

func (hndl *PowHandler) Challenge(w http.ResponseWriter, _ *http.Request) {
	challenge := hndl.pow.Challenge()
	w.Write(challenge) //nolint:errcheck
}

func (hndl *PowHandler) Validate(w http.ResponseWriter, r *http.Request) {
	defaultCtx := context.Background()

	challenge := []byte(r.FormValue(challengeArg))
	solution := []byte(r.FormValue(solutionArg))

	if err := hndl.pow.Verify(challenge, solution); err != nil {
		http.Error(w, "Invalid PoW solution", http.StatusUnauthorized)
		return
	}

	res, err := hndl.svc.GetWordOfWisdom(defaultCtx)
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(fmt.Sprintf("Quote of the day: %s", res))) //nolint:errcheck
}
