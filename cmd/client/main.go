package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	clientConfig "github.com/cortez720/wordofwisdom/config/client"
	powConfig "github.com/cortez720/wordofwisdom/config/pow"
	solverConfig "github.com/cortez720/wordofwisdom/config/solver"
	"github.com/cortez720/wordofwisdom/internal/handler"
	hashbasedpow "github.com/cortez720/wordofwisdom/internal/pkg/hash_based_pow"
	solverSvc "github.com/cortez720/wordofwisdom/internal/service/solver"
)

func main() {
	clientConfig := clientConfig.GetClientConfig()
	solverConfig := solverConfig.GetSolverConfig()
	powConfig := powConfig.GetPowConfig()

	pow, err := hashbasedpow.NewPOW(powConfig)
	if err != nil {
		log.Fatalf("hashbasedpow.NewPOW: %v", err.Error())
	}

	svc := solverSvc.NewService(solverConfig, pow)

	hndl := handler.NewSolver(svc)

	http.HandleFunc("/solve", hndl.Solve)

	log.Printf("Running Client HTTP server on %s\n", clientConfig.HTTPAddr)

	go func() { log.Fatal(http.ListenAndServe(clientConfig.HTTPAddr, nil)) }()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("closing")
}
