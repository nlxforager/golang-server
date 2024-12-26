package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const CONFIG_LOGGER_TYPE = "LOGGER_TYPE"

const CONFIG_NATS_EMBEDDED = "NATS_EMBEDDED"
const CONFIG_NATS_SERVER_URL = "NATS_SERVER_URL"

func Init() error {
	return godotenv.Load()
}

func NatsConfig() ( /* _embedded */ bool, string, error) {
	embedded := os.Getenv(CONFIG_NATS_EMBEDDED) == "TRUE"
	if embedded && os.Getenv(CONFIG_NATS_SERVER_URL) != "" {
		return false, "", fmt.Errorf("either embed or specify serverurl")
	}
	return embedded, os.Getenv(CONFIG_NATS_SERVER_URL), nil
}
