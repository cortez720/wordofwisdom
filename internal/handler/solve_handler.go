package handler

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	
	config "wordOfWisdom/config/client"
)

type solver interface {
	Solve(challenge []byte) []byte
}

type SolverHandler struct {
	cfg    *config.ClientConfig
	solver solver
}

func NewSolver(solver solver, cfg *config.ClientConfig) *SolverHandler {
	return &SolverHandler{solver: solver, cfg: cfg}
}

func (hndl *SolverHandler) Solve(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(hndl.cfg.ServerAddr + hndl.cfg.ChallengeRoute)
	if err != nil {
		log.Fatalf("http.Get: %v", err.Error())
	}

	defer resp.Body.Close()
	challenge, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll: %v", err.Error())
	}

	solution := hndl.solver.Solve(challenge)

	form := url.Values{}
	form.Add(challengeArg, string(challenge))
	form.Add(solutionArg, string(solution))

	resp, err = http.PostForm(hndl.cfg.ServerAddr+hndl.cfg.ValidateRoute, form)
	if err != nil {
		log.Fatalf("http.PostForm: %v", err.Error())
	}
	defer resp.Body.Close()

	quote, _ := ioutil.ReadAll(resp.Body)
	w.Write([]byte("Quote of the day: " + string(quote)))
}
