package configs

import (
	"os"
)

type EmailConfig struct {
	EmailHost   string
	EmailPort   string
	FromAddress string
	FromName    string
	UserName    string
	Password    string
}

func (env *EnvConfig) LoadEmailConfig() *EmailConfig {
	return &EmailConfig{
		EmailHost:   os.Getenv("EMAIL_HOST"),
		EmailPort:   os.Getenv("EMAIL_PORT"),
		FromAddress: os.Getenv("EMAIL_FROM_ADDRESS"),
		FromName:    os.Getenv("EMAIL_FROM_NAME"),
		UserName:    os.Getenv("EMAIL_USERNAME"),
		Password:    os.Getenv("EMAIL_PASSWORD"),
	}
}
