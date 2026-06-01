package config

import (
	"os"
	"time"
)

type Config struct {
	DatabaseURL      string
	JWTSecret        string
	ServerAddr       string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
}

func Load() Config {
	accessTTL := 15 * time.Minute
	refreshTTL := 7 * 24 * time.Hour

	if v := os.Getenv("ACCESS_TOKEN_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			accessTTL = d
		}
	}

	if v := os.Getenv("REFRESH_TOKEN_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			refreshTTL = d
		}
	}

	return Config{
		DatabaseURL:     envOrDefault("DATABASE_URL", "postgres://postgres:postgres@db:5432/account?sslmode=disable"),
		JWTSecret:       envOrDefault("JWT_SECRET", "change-me-secret"),
		ServerAddr:      envOrDefault("SERVER_ADDR", ":8080"),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
	}
}

func envOrDefault(name, fallback string) string {
	if v := os.Getenv(name); v != "" {
		return v
	}
	return fallback
}
