package config

import "os"

var GOAUTH_CLIENT_ID_PREFIX = ""
var GOAUTH_CLIENT_SECRET = ""
var GOAUTH_LISTENING_PORT = ""
var GOAUTH_AUTH_CODE_SUCCESS_CALLBACK_PATH = ""
var GOAUTH_ENABLE_LOG_REQUEST = false

type GoogleAuthConfig struct {
	CLIENT_ID_PREFIX                string
	CLIENT_SECRET                   string
	AUTH_CODE_SUCCESS_CALLBACK_PATH string
	ENABLE_LOG_REQUEST              bool
	Port                            string
}

func gauth(port string) (GoogleAuthConfig, error) {
	c := GoogleAuthConfig{
		CLIENT_ID_PREFIX:                os.Getenv("CLIENT_ID_PREFIX"),
		CLIENT_SECRET:                   os.Getenv("CLIENT_SECRET"),
		AUTH_CODE_SUCCESS_CALLBACK_PATH: os.Getenv("AUTH_CODE_SUCCESS_CALLBACK_PATH"),
		ENABLE_LOG_REQUEST:              os.Getenv("ENABLE_LOG_REQUEST") == "true",
		Port:                            port,
	}

	return c, nil
}
