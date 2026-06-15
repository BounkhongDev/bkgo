package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      App
	Postgres Postgres
	Redis    Redis
	MinIO    MinIO
	JWT      JWT
}

type App struct {
	Name string
	Port string
	Env  string
}

type Postgres struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	DSN      string
}

type Redis struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type MinIO struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

type JWT struct {
	Secret string
}

// Load reads environment variables (from .env file if present) into a Config
// and validates required fields. Accepts optional file paths; defaults to ".env".
func Load(files ...string) (*Config, error) {
	if len(files) > 0 {
		if err := godotenv.Load(files...); err != nil {
			return nil, fmt.Errorf("config: load env file: %w", err)
		}
	} else {
		_ = godotenv.Load()
	}

	redisDB, _ := strconv.Atoi(getEnv("REDIS_DB", "0"))
	useSSL, _ := strconv.ParseBool(getEnv("MINIO_USE_SSL", "false"))

	cfg := &Config{
		App: App{
			Name: getEnv("APP_NAME", "app"),
			Port: getEnv("APP_PORT", "8080"),
			Env:  getEnv("APP_ENV", "development"),
		},
		Postgres: Postgres{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "app_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		Redis: Redis{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       redisDB,
		},
		MinIO: MinIO{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadmin"),
			Bucket:    getEnv("MINIO_BUCKET", "app-files"),
			UseSSL:    useSSL,
		},
		JWT: JWT{
			Secret: getEnv("JWT_SECRET", ""),
		},
	}

	cfg.Postgres.DSN = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Postgres.Host, cfg.Postgres.Port,
		cfg.Postgres.User, cfg.Postgres.Password,
		cfg.Postgres.DBName, cfg.Postgres.SSLMode,
	)

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.JWT.Secret == "" {
		return errors.New("config: JWT_SECRET must not be empty")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
