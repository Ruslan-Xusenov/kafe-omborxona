package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword       string
	DBName           string
	JWTSecret        string
	ServerPort       string
	FrontendURL      string
	TelegramBotToken string
}

func Load() (*Config, error) {
	loadEnvFile(".env")

	cfg := &Config{
		DBHost:           getEnv("DB_HOST", "localhost"),
		DBPort:           getEnv("DB_PORT", "5432"),
		DBUser:           getEnv("DB_USER", "postgres"),
		DBPassword:       getEnv("DB_PASSWORD", "postgres"),
		DBName:           getEnv("DB_NAME", "kafe_omborxona"),
		JWTSecret:        getEnv("JWT_SECRET", "default-secret"),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		FrontendURL:      getEnv("FRONTEND_URL", "http://localhost:3000"),
		TelegramBotToken: getEnv("TELEGRAM_BOT_TOKEN", ""),
	}

	if cfg.JWTSecret == "default-secret" {
		return nil, fmt.Errorf("JWT_SECRET must be set")
	}

	return cfg, nil
}

func (c *Config) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func loadEnvFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			if _, exists := os.LookupEnv(key); !exists {
				os.Setenv(key, val)
			}
		}
	}
}
