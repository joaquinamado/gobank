package utils

import (
	godotenv "github.com/joho/godotenv"
	"log"
	"sync"
)

type EnvVariables struct {
	JwtSecret string
}

var lock = &sync.Mutex{}

type single struct {
	EnvVariables *EnvVariables
}

var singleInstance *single

func GetEnvInstance() *single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = &single{
				EnvVariables: newEnvVariables(),
			}
		}
	}
	return singleInstance
}

func newEnvVariables() *EnvVariables {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	envFile, err := godotenv.Read()

	if err != nil {
		log.Fatal("Error reading .env file")
	}

	return &EnvVariables{
		JwtSecret: envFile["JWT_SECRET"],
	}
}
