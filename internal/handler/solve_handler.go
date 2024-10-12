//go:generate mockgen -destination=./mocks/solve_handler_mock.go -source=./solve_handler.go -package=handler

package handler

import (
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

func (hndl *SolverHandler) Solve(w http.ResponseWriter, r *http.Request) {
	quote, err := hndl.svc.Solve()
	if err != nil {
		http.Error(w, "Internal error.", http.StatusInternalServerError)
		return
	}

	w.Write([]byte(string(quote)))
}
