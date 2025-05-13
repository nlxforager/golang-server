package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var CLIENT_ID_PREFIX = ""
var CLIENT_SECRET = ""
var LISTENING_PORT = ""
var AUTH_CODE_SUCCESS_ENDPOINT_PATH = ""
var ENABLE_LOG_REQUEST = false

func InitConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file")
	} else {
		fmt.Println("Loading .env file ok.")
	}

	CLIENT_ID_PREFIX = os.Getenv("CLIENT_ID_PREFIX")
	CLIENT_SECRET = os.Getenv("CLIENT_SECRET")
	LISTENING_PORT = os.Getenv("LISTENING_PORT")
	AUTH_CODE_SUCCESS_ENDPOINT_PATH = os.Getenv("AUTH_CODE_SUCCESS_ENDPOINT_PATH")
	ENABLE_LOG_REQUEST = os.Getenv("ENABLE_LOG_REQUEST") == "true"

	return nil
}
