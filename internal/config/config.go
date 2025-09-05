package config

import "os"

type Config struct {
	Port        string
	Environment string
	IsDev       bool
}

func NewConfig() *Config {
	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	env := os.Getenv("ENV")
	isDev := env != "production"

	return &Config{
		Port:        port,
		Environment: env,
		IsDev:       isDev,
	}
}