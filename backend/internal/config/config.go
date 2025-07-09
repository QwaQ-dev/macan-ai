package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string `yaml:"env" env-default:"dev" env-requried:"true"`
	JWTSecretKey string `yaml:"jwtsecretkey"`
	Server       `yaml:"server"`
	Database     `yaml:"database"`
	Services     `yaml:"services""`
}

type Server struct {
	Port string `yaml:"port" env-default:":8080"`
}

type Services struct {
	ResumeParsingGRPCAddr string `yaml:"resume_parsing_grpc_addr" env:"RESUME_PARSING_GRPC_ADDR"`
}

type Database struct {
	Port       string `yaml:"port"`
	DBhost     string `yaml:"host"`
	DBname     string `yaml:"db_name"`
	DBpassword string `yaml:"db_password"`
	SSLMode    string `yaml:"sslmode"`
	DBusername string `yaml:"db_username"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG")

	if configPath == "" {
		log.Fatalf("No env file")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file is not exists: %s", err.Error())
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config: %s", err.Error())
	}

	return &cfg
}
