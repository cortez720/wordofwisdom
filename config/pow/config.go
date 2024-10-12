package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// Config
type PowConfig struct {
	Complexity uint64 `envconfig:"POW_COMPLEXITY"`
}

var (
	config PowConfig
	once   sync.Once
)

// Get reads config from environment. Once.
func GetPowConfig() *PowConfig {
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
