package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	config "wordOfWisdom/config/client"
	powConfig "wordOfWisdom/config/pow"
	"wordOfWisdom/internal/handler"
	hashbasedpow "wordOfWisdom/internal/pkg/hash_based_pow"
)

func main() {
	clientConfig := config.GetClientConfig()

	//pow config
	pow, err := hashbasedpow.NewPOW(powConfig.GetPowConfig())
	if err != nil {
		log.Fatalf("hashbasedpow.NewPOW: %v", err.Error())
	}

	hndl := handler.NewSolver(pow, clientConfig)

	// routes
	http.HandleFunc("/solve", hndl.Solve)

	log.Printf("Running HTTP server on %s\n", clientConfig.HTTPAddr)

	go func() { log.Fatal(http.ListenAndServe(clientConfig.HTTPAddr, nil)) }()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("closing")

}
