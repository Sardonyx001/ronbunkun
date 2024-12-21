package config

import "os"

type HTTPConfig struct {
	Host       string
	Port       string
	ExposePort string
}

func LoadHTTPConfig() HTTPConfig {
	return HTTPConfig{
		Host:       os.Getenv("RONBUNKUN_API_HOST"),
		Port:       os.Getenv("RONBUNKUN_API_PORT"),
		ExposePort: os.Getenv("RONBUNKUN_API_EXPOSE_PORT"),
	}
}
