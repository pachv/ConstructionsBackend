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

	cfg.Database.User = os.Getenv("ADMIN_DB_USER")
	cfg.Database.Password = os.Getenv("ADMIN_DB_PASSWORD")
	cfg.Database.Name = os.Getenv("ADMIN_DB_NAME")
	cfg.Database.Host = os.Getenv("ADMIN_DB_HOST")
	cfg.Database.Port = os.Getenv("ADMIN_DB_PORT")

	cfg.Server.Port = os.Getenv("ADMIN_PORT")
	cfg.Server.Domain = os.Getenv("ADMIN_DOMAIN")

	return cfg

}
