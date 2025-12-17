package config

import (
	"os"

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

	Server struct {
		Port   string
		Domain string
	}
}

func MustLoadConfig(path string) *Config {

	cfg := &Config{}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("no config file on path " + path)
	}

	_ = godotenv.Load(path)

	cfg.Database.User = os.Getenv("POSTGRES_USER")
	cfg.Database.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Database.Name = os.Getenv("POSTGRES_DB")
	cfg.Database.Host = os.Getenv("POSTGRES_HOST")
	cfg.Database.Port = os.Getenv("POSTGRES_PORT")

	cfg.Server.Port = os.Getenv("PORT")
	cfg.Server.Domain = os.Getenv("DOMAIN")

	return cfg

}
