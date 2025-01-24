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

type NatsConfig struct {
	Embedded bool
	Url      string
}

func GetNatsConfig() (NatsConfig, error) {
	embedded, url := os.Getenv(CONFIG_NATS_EMBEDDED) == "TRUE", os.Getenv(CONFIG_NATS_SERVER_URL)
	if embedded && url != "" {
		return NatsConfig{}, fmt.Errorf("either set embedded server or specify serverurl, not both")
	}
	return NatsConfig{
		Embedded: embedded,
		Url:      url,
	}, nil
}
