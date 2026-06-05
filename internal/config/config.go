// Enchanted-Garden/internal/config/config.go
package config

import "os"

type HTTPConfig struct {
	Port string
}

type DBConfig struct {
	DSN string
}

type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
}

func Load() *Config {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = ":8080"
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/garden_db?sslmode=disable"
	}

	return &Config{
		HTTP: HTTPConfig{Port: port},
		DB:   DBConfig{DSN: dsn},
	}
}
