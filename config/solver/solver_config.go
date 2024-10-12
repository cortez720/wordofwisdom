package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// Config.
type SolverConfig struct {
	ServerAddr     string `envconfig:"SERVER_ADDR"`
	ChallengeRoute string `envconfig:"CHALLENGE_ROUTE"`
	ValidateRoute  string `envconfig:"VALIDATE_ROUTE"`
}

var (
	config SolverConfig
	once   sync.Once
)

// Get reads config from environment. Once.
func GetSolverConfig() *SolverConfig {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}

		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Configuration:", string(configBytes))
	})

	return &config
}
