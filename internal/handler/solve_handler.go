//go:generate mockgen -destination=./mocks/solve_handler_mock.go -source=./solve_handler.go -package=handler

package handler

import (
	"log"
	"net/http"
)

type SolverService interface {
	Solve() ([]byte, error)
}

type SolverHandler struct {
	svc SolverService
}

func NewSolver(solver SolverService) *SolverHandler {
	return &SolverHandler{svc: solver}
}

func (hndl *SolverHandler) Solve(w http.ResponseWriter, _ *http.Request) {
	quote, err := hndl.svc.Solve()
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		log.Printf("hndl.svc.Solve: %v", err.Error())

		return
	}

	w.Write([]byte(string(quote))) //nolint:errcheck
}
