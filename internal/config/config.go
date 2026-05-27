package config

import "github.com/spf13/viper"

type Config struct {
	Port                string
	LogLevel            string
	ReviewIntervalHours int
	DBHost              string
	DBPort              int
	DBName              string
	DBUser              string
	DBPass              string
	RedisAddr           string
}

func Load() *Config {
	v := viper.New()

	v.SetDefault("PORT", "8080")
	v.SetDefault("LOG_LEVEL", "info")
	v.SetDefault("REVIEW_INTERVAL_HOURS", 24)
	v.SetDefault("DB_HOST", "localhost")
	v.SetDefault("DB_PORT", 5432)
	v.SetDefault("DB_NAME", "vocab")
	v.SetDefault("DB_USER", "postgres")
	v.SetDefault("DB_PASS", "secret")
	v.SetDefault("REDIS_ADDR", "localhost:6379")

	v.AutomaticEnv()
	v.BindEnv("PORT", "LOG_LEVEL", "REVIEW_INTERVAL_HOURS", "DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASS", "REDIS_ADDR")

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	_ = v.ReadInConfig()

	return &Config{
		Port:                v.GetString("PORT"),
		LogLevel:            v.GetString("LOG_LEVEL"),
		ReviewIntervalHours: v.GetInt("REVIEW_INTERVAL_HOURS"),
		DBHost:              v.GetString("DB_HOST"),
		DBPort:              v.GetInt("DB_PORT"),
		DBName:              v.GetString("DB_NAME"),
		DBUser:              v.GetString("DB_USER"),
		DBPass:              v.GetString("DB_PASS"),
		RedisAddr:           v.GetString("REDIS_ADDR"),
	}
}
