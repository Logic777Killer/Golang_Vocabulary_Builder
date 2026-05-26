package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Port                string
	LogLevel            string
	ReviewIntervalHours int
	DatabaseURL         string
}

func Load() *Config {
	v := viper.New()

	v.SetDefault("PORT", "8080")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("REVIEW_INTERVAL_HOURS", 24)
	v.SetDefault("DATABASE_URL", "")

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	_ = v.ReadInConfig()

	v.AutomaticEnv()

	return &Config{
		Port:                v.GetString("PORT"),
		LogLevel:            v.GetString("LOG_LEVEL"),
		ReviewIntervalHours: v.GetInt("REVIEW_INTERVAL_HOURS"),
		DatabaseURL:         v.GetString("DATABASE_URL"),
	}
}
