package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"wordOfWisdom/internal/handler"

	powConfig "wordOfWisdom/config/pow"
	serverConfig "wordOfWisdom/config/server"

	hashbasedpow "wordOfWisdom/internal/pkg/hash_based_pow"
	quoteSvc "wordOfWisdom/internal/service/quote"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// default context
	ctx := context.Background()

	// server config
	srvCfg := serverConfig.GetServerConfig()

	// pow config
	pow, err := hashbasedpow.NewPOW(powConfig.GetPowConfig())
	if err != nil {
		return fmt.Errorf("hashbasedpow.NewPOW: %w", err)
	}

	svc := quoteSvc.NewService(ctx)

	hndl := handler.NewPow(pow, svc)

	// routes
	http.HandleFunc("/validate", hndl.Validate)
	http.HandleFunc("/challenge", hndl.Challenge)

	log.Printf("Running HTTP server on %s\n", srvCfg.HTTPAddr)

	go func() { log.Fatal(http.ListenAndServe(srvCfg.HTTPAddr, nil)) }()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	fmt.Println("closing")

	return nil
}
