// Enchanted-Garden/internal/config/config.go
package config

import "os"

type Config struct {
	HTTP struct {
		Port string
	}
	DB struct {
		DSN string
	}
}

func Load() *Config {
	cfg := &Config{}
	cfg.HTTP.Port = os.Getenv("SERVER_PORT")
	if cfg.HTTP.Port == "" {
		cfg.HTTP.Port = ":8080"
	}
	cfg.DB.DSN = os.Getenv("DATABASE_URL")
	return cfg
}
