package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	Mail     MailConfig
}

type AppConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// DSN returns the PostgreSQL data source name.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		d.Host, d.User, d.Password, d.Name, d.Port, d.SSLMode,
	)
}

type JWTConfig struct {
	Secret      string
	ExpireHours int
}

type RedisConfig struct {
	Addr     string // host:port
	Password string
	DB       int
}

type MailConfig struct {
	Host     string
	Port     int
	From     string
	Username string
	Password string
}

// Load reads configuration from environment variables (.env file is optional).
func Load() (*Config, error) {
	// Ignore error — .env is optional (may not exist in production).
	_ = godotenv.Load()

	expireHours, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "72"))
	if err != nil {
		expireHours = 72
	}

	redisDB, err := strconv.Atoi(getEnv("REDIS_DB", "0"))
	if err != nil {
		redisDB = 0
	}

	mailPort, err := strconv.Atoi(getEnv("MAIL_PORT", "1025"))
	if err != nil {
		mailPort = 1025
	}

	return &Config{
		App: AppConfig{
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "realworld"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "change-me-in-production"),
			ExpireHours: expireHours,
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		Mail: MailConfig{
			Host:     getEnv("MAIL_HOST", "localhost"),
			Port:     mailPort,
			From:     getEnv("MAIL_FROM", "noreply@realworld.dev"),
			Username: getEnv("MAIL_USERNAME", ""),
			Password: getEnv("MAIL_PASSWORD", ""),
		},
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
