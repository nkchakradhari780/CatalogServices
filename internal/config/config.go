package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServer struct {
	Addr string `yaml:"address" env:"HTTP_SERVER_ADDRESS" env-require:"true"`
}

type Database struct {
    Host     string `yaml:"host" env:"DATABASE_HOST" env-required:"true"`
    Port     int    `yaml:"port" env:"DATABASE_PORT" env-required:"true"`
    Name     string `yaml:"name" env:"DATABASE_NAME" env-required:"true"`
    Username string `yaml:"username" env:"DATABASE_USERNAME" env-required:"true"`
    Password string `yaml:"password" env:"DATABASE_PASSWORD" env-required:"true"`
    SSLMode  string `yaml:"sslmode" env:"DATABASE_SSLMODE" env-default:"disable"`
}


type Config struct {
    Env        string     `yaml:"env" env:"ENV" env-required:"true" env-default:"production"`
    HTTPServer HTTPServer `yaml:"http_server" env-required:"true"`
    Database   Database   `yaml:"database" env-required:"true"`
}


func MustLoad() *Config {
	var configPath string
	configPath = os.Getenv("CONFIG_PATH")

	if configPath == "" {
		flags := flag.String("config", "D:/Infosoft_solutions/CatalogServices/config/local.yaml", "path to the configration file")
		flag.Parse()
  
		configPath = *flags

		if configPath == "" {
			log.Fatal("Config path is not set")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", configPath)
	}

	var cfg Config 

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("Can not read config file %s: %s", configPath, err.Error())
	}

	return &cfg
}