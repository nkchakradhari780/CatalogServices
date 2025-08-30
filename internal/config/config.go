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

type Config struct {
	Env        string     `yaml:"env" env:"ENV" env-require:"true" env-default:"production"`
	HTTPServer HTTPServer `yaml:"http_server" env:"HTTP_SERVER" env-require:"true"`
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