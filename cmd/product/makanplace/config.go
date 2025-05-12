package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var CLIENT_ID_PREFIX = ""
var CLIENT_SECRET = ""
var LISTENING_PORT = ""
var AUTH_CODE_SUCCESS_CALLBACK_PATH = ""
var ENABLE_LOG_REQUEST = false

func InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env file\n")
	}

	// auth
	CLIENT_ID_PREFIX = os.Getenv("CLIENT_ID_PREFIX")
	CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
	LISTENING_PORT = os.Getenv("LISTENING_PORT")
	AUTH_CODE_SUCCESS_CALLBACK_PATH = os.Getenv("AUTH_CODE_SUCCESS_CALLBACK_PATH")
	ENABLE_LOG_REQUEST = os.Getenv("ENABLE_LOG_REQUEST") == "true"

	return nil
}
