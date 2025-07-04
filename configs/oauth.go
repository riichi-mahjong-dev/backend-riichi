package configs

import (
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type ConfigOauth struct {
	Config *oauth2.Config
	Key    string
}

func (env *EnvConfig) LoadOAuthConfig() *ConfigOauth {
	return &ConfigOauth{
		Config: &oauth2.Config{
			ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
			ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("OAUTH_REDIRECT_URL"),
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     google.Endpoint,
		},
		Key: os.Getenv("OAUTH_KEY"),
	}
}
