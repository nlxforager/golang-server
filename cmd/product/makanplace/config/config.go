package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	ConnString string
}

type Config struct {
	DatabaseConfig
	GoogleAuthConfig
	ServerConfig
}

func InitConfig() (c Config, e error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, fmt.Errorf("error loading .env file\n")
	}
	c.ServerConfig.Port = os.Getenv("LISTENING_PORT")
	// database

	dbUrl := os.Getenv("DATABASE_URL")
	c.DatabaseConfig.ConnString = dbUrl

	// auth
	c.GoogleAuthConfig, err = gauth(c.ServerConfig.Port)
	if err != nil {
		return Config{}, err
	}

	// server

	return c, nil
}
