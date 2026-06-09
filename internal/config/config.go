package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"log"
	"time"
)

type Config struct {
	Port          string        `env:"PORT" envDefault:":8080"`
	DBHost        string        `env:"DB_HOST" envDefault:"localhost"`
	DBPort        string        `env:"DB_PORT" envDefault:"5432"`
	DBUser        string        `env:"DB_USER" envDefault:"postgres"`
	DBPassword    string        `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName        string        `env:"DB_NAME" envDefault:"expenses"`
	DBSSLMode     string        `env:"DB_SSLMODE" envDefault:"disable"`
	JWTSecret     string        `env:"JWT_SECRET" envDefault:"super-secret-key-change-in-prod"`
	JWTExpiration time.Duration `env:"JWT_EXPIRATION" envDefault:"24h"`
}

func Load() (*Config, error) {
	// Load .env file if it exists (safe to ignore in prod/CI)
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ No .env file found, using system env vars")
	}

	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
