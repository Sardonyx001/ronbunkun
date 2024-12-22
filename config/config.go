package config

type Config struct {
	HTTP HTTPConfig
}

func NewConfig() *Config {
	/*
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file %v", err)
		}
	*/

	return &Config{
		HTTP: LoadHTTPConfig(),
	}
}
