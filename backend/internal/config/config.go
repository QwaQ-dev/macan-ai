package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string `yaml:"env" env-required:"true"`
	Http_server `yaml:"http_server"`
	Database    `yaml:"database"`
}

type Http_server struct {
	Address      string        `yaml:"address"`
	Timeout      time.Duration `yaml:"timeout"`
	Idle_timeout time.Duration `yaml:"idle_timeout"`
}

type Database struct {
	Host       string `yaml:"host" env-required:"true"`
	Port       string `yaml:"port"`
	DBname     string `yaml:"db_name"`
	DBpassword string `yaml:"db_password"`
	DBusername string `yaml:"db_username"`
	SSLmode    string `yaml:"sslmode" env-default:"disable"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG")
	if configPath == "" {
		log.Fatalf("There is no path to config file: %s", configPath)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("There is no config file: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error with reading config: %s", err)
	}

	return &cfg
}
