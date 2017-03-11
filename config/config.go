package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/cloudfoundry-community/go-cfenv"
	redis "gopkg.in/redis.v5"
)

//GetPort : Get TCP PORT from CF env
func GetPort() string {
	var portNumber int

	cfEnv, err := cfenv.Current()
	if err != nil {
		portNumber = 3000
	} else {
		portNumber = cfEnv.Port
	}

	return strconv.Itoa(portNumber)
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

//GetRedisConf get Redis connection information
func GetRedisConf() redis.Options {
	//set default
	var options = redis.Options{Addr: "localhost:6379",
		Password: "",
		DB:       0,
	}

	//read cfEnv
	cfEnv, err := cfenv.Current()
	if err == nil {
		fmt.Println(err)
		redisServices, err := cfEnv.Services.WithTag("redis")
		if err != nil {
			panic(err)
		}

		redisCreds := redisServices[0].Credentials

		options = redis.Options{Addr: redisCreds["hostname"].(string) + ":" + redisCreds["port"].(string),
			Password: redisCreds["password"].(string),
			DB:       0,
		}
	}
	return options
}
