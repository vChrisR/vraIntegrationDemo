package config

import (
	"os"
	"testing"
)

//TestGetPort tests the config.GetPort function
func TestGetPortDefault(t *testing.T) {
	os.Setenv("PORT", "")
	port := GetPort()
	if port != "3000" {
		t.Error("Expected 3000 but got ", port)
	}
}

func TestGetPortEnv(t *testing.T) {
	os.Setenv("PORT", "80")
	port := GetPort()
	if port != "80" {
		t.Error("Expected 80 but got ", port)
	}
}

func TestGetAPICredsDefault(t *testing.T) {
	os.Setenv("APIUSER", "")
	os.Setenv("APIPASS", "")
	apiUser, apiPass := GetAPICreds()
	if apiUser != "api" {
		t.Error("Expected api but got ", apiUser)
	}
	if apiPass != "api" {
		t.Error("Expected api but got ", apiPass)
	}
}

func TestGetAPICredsEnv(t *testing.T) {
	os.Setenv("APIUSER", "testuser")
	os.Setenv("APIPASS", "testpassword")
	apiUser, apiPass := GetAPICreds()
	if apiUser != "testuser" {
		t.Error("Expected api but got ", apiUser)
	}
	if apiPass != "testpassword" {
		t.Error("Expected api but got ", apiPass)
	}
}
