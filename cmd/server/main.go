package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	powConfig "github.com/cortez720/wordofwisdom/config/pow"
	serverConfig "github.com/cortez720/wordofwisdom/config/server"
	"github.com/cortez720/wordofwisdom/internal/handler"
	hashbasedpow "github.com/cortez720/wordofwisdom/internal/pkg/hash_based_pow"
	quoteSvc "github.com/cortez720/wordofwisdom/internal/service/quote"
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
