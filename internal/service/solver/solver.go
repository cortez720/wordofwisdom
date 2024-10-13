package handler

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	config "github.com/cortez720/wordofwisdom/config/solver"
)

const (
	challengeArg = "challenge"
	solutionArg  = "solution"
)

type solver interface {
	Solve(challenge []byte) []byte
}

type Solver struct {
	cfg    *config.SolverConfig
	solver solver
}

func NewService(cfg *config.SolverConfig, solver solver) *Solver {
	return &Solver{cfg: cfg, solver: solver}
}

func (s *Solver) Solve() ([]byte, error) {
	resp, err := http.Get(s.cfg.ServerAddr + s.cfg.ChallengeRoute)
	if err != nil {
		return nil, fmt.Errorf("http.Get: %w", err)
	}

	defer resp.Body.Close()

	challenge, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	solution := s.solver.Solve(challenge)

	form := url.Values{}
	form.Add(challengeArg, string(challenge))
	form.Add(solutionArg, string(solution))

	resp, err = http.PostForm(s.cfg.ServerAddr+s.cfg.ValidateRoute, form)
	if err != nil {
		return nil, fmt.Errorf("http.PostForm: %w", err)
	}
	defer resp.Body.Close()

	quote, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}

	return quote, nil
}
