package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	JWTSecret string


	GitHubClientID     string
	GitHubClientSecret string
	GitHubRedirectURL  string
}

func Load() *Config {
	_, thisFile, _, _ := runtime.Caller(0)
	_ = godotenv.Load(filepath.Join(filepath.Dir(thisFile), "..", "..", ".env"))

	return &Config{
		AppPort: getEnv("APP_PORT", "8081"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "k3s_deploy"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
		GitHubClientID: getEnv("GITHUB_CLIENT_ID", ""),
		GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		GitHubRedirectURL: getEnv("GITHUB_REDIRECT_URL", ""),

		JWTSecret: getEnv("JWT_SECRET", "dev-secret-change-me"),
	}
}




func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}
