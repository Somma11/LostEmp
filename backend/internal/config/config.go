package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
	JWTTTL    time.Duration
}

func Load() *Config {
	return &Config{
		Port:      envOr("PORT", "3000"),
		DBPath:    envOr("DB_PATH", "./data/lostemp.db"),
		JWTSecret: envOr("JWT_SECRET", "dev-secret-change-me"),
		JWTTTL:    time.Duration(envInt("JWT_TTL_HOURS", 12)) * time.Hour,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func envInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		n := 0
		_, err := fmt.Sscanf(v, "%d", &n)
		if err == nil {
			return n
		}
	}
	return fallback
}
