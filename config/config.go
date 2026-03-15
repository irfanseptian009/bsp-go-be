package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   time.Duration
}

// Load reads configuration from environment variables.
func Load() *Config {
	// Load .env file if it exists (ignored in production)
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	jwtExpiry := 24 * time.Hour
	if exp := os.Getenv("JWT_EXPIRY"); exp != "" {
		if parsed, err := time.ParseDuration(exp); err == nil {
			jwtExpiry = parsed
		}
	}

	return &Config{
		Port:        port,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTExpiry:   jwtExpiry,
	}
}
