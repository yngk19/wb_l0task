package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Env        string `yaml:"env" env-default:"prod"`
	HTTPServer `yaml:"http_server"`
}

type HTTPServer struct {
	Address      string `yaml:"address" env-default:"localhost:8080"`
	Timeout      string `yaml:"timeout" env-default:"4s"`
	IddleTimeout string `yaml:"iddle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Cannot load the .env file!: %s", err)
	}
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("Config file is not set!")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file doesn't exist!: %s", err)
	}
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read the config!: %s", err)
	}
	return &cfg
}
