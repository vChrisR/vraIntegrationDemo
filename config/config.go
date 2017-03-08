package config

import (
	"os"
)

//GetPort : Get TCP PORT from os env
func GetPort() string {
	configuredPort := os.Getenv("PORT")
	if configuredPort == "" {
		return "3000"
	}

	return configuredPort
}

//GetAPICreds : Get api creds from OS env
func GetAPICreds() (string, string) {
	apiuser := os.Getenv("APIUSER")
	if apiuser == "" {
		apiuser = "api"
	}

	apipass := os.Getenv("APIPASS")
	if apipass == "" {
		apipass = "api"
	}

	return apiuser, apipass
}
