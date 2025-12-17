package config

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database struct {
		User     string
		Password string
		Host     string
		Port     string
		Name     string
	}

	Logger struct {
		Level string
	}

	JWT struct {
		Secret string
	}

	Email struct {
		From     string
		Password string
		Host     string
		Port     int

		NotifyEmail string
	}

	Port string `env:"PORT" envDefault:"8886"`
}

func LoadConfig(path string) (*Config, error) {

	config := &Config{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf(".env file not found at path: %s", path)
	}

	slog.Info("file exists")
	_ = godotenv.Load(path)

	config.Database.User = os.Getenv("MAIN_DB_USER")
	config.Database.Password = os.Getenv("MAIN_DB_PASSWORD")
	config.Database.Host = os.Getenv("MAIN_DB_HOST")
	config.Database.Port = os.Getenv("MAIN_DB_PORT")
	config.Database.Name = os.Getenv("MAIN_DB_NAME")

	config.Logger.Level = os.Getenv("LOG_LEVEL")

	config.Port = ":" + os.Getenv("MAIN_PORT")

	config.JWT.Secret = os.Getenv("MAIN_JWT_SECRET")

	config.Email.From = os.Getenv("MAIN_EMAIL_FROM")
	config.Email.Host = os.Getenv("SMTP_HOST")
	config.Email.Password = os.Getenv("MAIN_EMAIL_PASSWORD")
	config.Email.NotifyEmail = os.Getenv("NOTIFY_EMAIL")
	emailPortStr := os.Getenv("SMTP_PORT")

	emailPort, err := strconv.Atoi(emailPortStr)
	if err != nil {
		log.Fatalf("invalid EMAIL_SMTP_PORT: %v", err)
	}

	config.Email.Port = emailPort

	return config, nil
}
