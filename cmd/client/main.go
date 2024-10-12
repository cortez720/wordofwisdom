package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"wordOfWisdom/internal/handler"

	clientConfig "wordOfWisdom/config/client"
	powConfig "wordOfWisdom/config/pow"
	solverConfig "wordOfWisdom/config/solver"

	hashbasedpow "wordOfWisdom/internal/pkg/hash_based_pow"
	solverSvc "wordOfWisdom/internal/service/solver"
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

	log.Printf("Running HTTP server on %s\n", clientConfig.HTTPAddr)

	go func() { log.Fatal(http.ListenAndServe(clientConfig.HTTPAddr, nil)) }()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("closing")
}
