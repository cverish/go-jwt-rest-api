package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server struct {
		Port         string
		Host         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
	}

	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
		SSLMode  string
	}

	API struct {
		BasePath string
	}

	JWT struct {
		AccessTokenKey     string
		AccessTokenSecret  string
		AccessTokenExpiry  time.Duration
		RefreshTokenKey    string
		RefreshTokenSecret string
		RefreshTokenExpiry time.Duration
		BackendRefresh     bool
		CookieDomain       string
	}

	Environment string
	IsDev       bool
}

func Load(envpath string) (*Config, error) {
	if _, err := os.Stat(envpath); err != nil {
		return nil, err
	}

	godotenv.Load(envpath)

	cfg := &Config{}

	// Server config
	cfg.Server.Port = getEnv("SERVER_PORT", "8080")
	cfg.Server.Host = getEnv("SERVER_HOST", "0.0.0.0")
	cfg.Server.ReadTimeout = time.Second * 15
	cfg.Server.WriteTimeout = time.Second * 15

	// Database config
	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "postgres")
	cfg.Database.Password = getEnv("DB_PASSWORD", "")
	cfg.Database.DBName = getEnv("DB_NAME", "postgres")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE", "disable")

	// API config
	cfg.API.BasePath = getEnv("API_BASE_PATH", "/api")

	// JWT Config
	cfg.JWT.AccessTokenKey = "access_token"
	cfg.JWT.AccessTokenSecret = getEnv("JWT_TOKEN_SECRET", "your-secret-key")
	cfg.JWT.AccessTokenExpiry = stringToDuration(
		getEnv("JWT_TOKEN_EXPIRY_MINUTES", "10"),
	) * time.Minute
	cfg.JWT.RefreshTokenKey = "refresh_token"
	cfg.JWT.RefreshTokenSecret = getEnv("JWT_REFRESH_SECRET", "your-secret-refresh-key")
	cfg.JWT.RefreshTokenExpiry = stringToDuration(
		getEnv("JWT_REFRESH_EXPIRY_DAYS", "7"),
	) * time.Hour * 24
	cfg.JWT.BackendRefresh = getEnv("JWT_BACKEND_REFRESH", "false") == "true"
	cfg.JWT.CookieDomain = getEnv("JWT_COOKIE_DOMAIN", "localhost")

	cfg.Environment = getEnv("ENV", "development")
	cfg.IsDev = cfg.Environment == "local" || cfg.Environment == "development"

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func stringToDuration(value string) time.Duration {
	n, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	return time.Duration(n)
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
