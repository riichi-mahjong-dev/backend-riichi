package configs

import (
	"os"
)

type JwtConfig struct {
	SecretKey  string
	RefreshKey string
}

func (env *EnvConfig) LoadJwtConfig() *JwtConfig {
	return &JwtConfig{
		SecretKey:  os.Getenv("SECRET_KEY_ACCESS"),
		RefreshKey: os.Getenv("SECRET_KEY_REFRESH"),
	}
}
