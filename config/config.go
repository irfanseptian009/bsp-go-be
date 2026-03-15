package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	Port                   string
	DatabaseURL            string
	JWTSecret              string
	JWTExpiry              time.Duration
	SupabaseURL            string
	SupabaseServiceRoleKey string
	SupabaseStorageBucket  string
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

	databaseURL := normalizeEnvValue(os.Getenv("DATABASE_URL"), "DATABASE_URL")
	jwtSecret := normalizeEnvValue(os.Getenv("JWT_SECRET"), "JWT_SECRET")
	supabaseURL := normalizeEnvValue(os.Getenv("SUPABASE_URL"), "SUPABASE_URL")
	supabaseServiceRoleKey := normalizeEnvValue(os.Getenv("SUPABASE_SERVICE_ROLE_KEY"), "SUPABASE_SERVICE_ROLE_KEY")
	supabaseStorageBucket := normalizeEnvValue(os.Getenv("SUPABASE_STORAGE_BUCKET"), "SUPABASE_STORAGE_BUCKET")
	if supabaseStorageBucket == "" {
		supabaseStorageBucket = "profile"
	}

	return &Config{
		Port:                   port,
		DatabaseURL:            databaseURL,
		JWTSecret:              jwtSecret,
		JWTExpiry:              jwtExpiry,
		SupabaseURL:            supabaseURL,
		SupabaseServiceRoleKey: supabaseServiceRoleKey,
		SupabaseStorageBucket:  supabaseStorageBucket,
	}
}

func normalizeEnvValue(value string, key string) string {
	v := strings.TrimSpace(value)
	v = strings.Trim(v, "\"'")

	prefix := key + "="
	if strings.HasPrefix(v, prefix) {
		v = strings.TrimSpace(strings.TrimPrefix(v, prefix))
		v = strings.Trim(v, "\"'")
	}

	return v
}
