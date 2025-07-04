package configs

import (
	"log"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
}

func LoadEnv() *EnvConfig {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("unable to load env file: %e", err)
	}

	return &EnvConfig{}
}
