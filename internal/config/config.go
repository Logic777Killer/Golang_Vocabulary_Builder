package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                string
	DatabaseURL         string
	LogLevel            string
	ReviewIntervalHours int
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	intervalStr := os.Getenv("REVIEW_INTERVAL_HOURS")
	interval := 24
	if val, err := strconv.Atoi(intervalStr); err == nil {
		interval = val
	}

	return &Config{
		Port:                port,
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		LogLevel:            logLevel,
		ReviewIntervalHours: interval,
	}
}
