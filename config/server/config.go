package config

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// Config
type ServerConfig struct {
	HTTPAddr string `envconfig:"SERVER_HTTP_HOST_ADDR"`
}

var (
	config ServerConfig
	once   sync.Once
)

// Get reads config from environment. Once.
func GetServerConfig() *ServerConfig {
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
