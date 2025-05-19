package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	ConnString string
}

type SuperUser struct {
	Gmails []string
}

type AdminConfig struct {
	SuperUser
}

type Config struct {
	DatabaseConfig
	GoogleAuthConfig
	ServerConfig
	AdminConfig
}

func InitConfig() (c Config, e error) {

	err := godotenv.Load()
	if !(os.Getenv("OPTIONAL_LOAD_ENV_FILE") == "TRUE") && err != nil {
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

	superUserGmails := os.Getenv("SUPER_USER_GMAIL")
	c.SuperUser.Gmails = strings.Split(superUserGmails, ",")
	return c, nil
}
