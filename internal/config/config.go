package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Service
	DB
	Nats
}

type Service struct {
	Env        string `yaml:"env" env-default:"prod"`
	HTTPServer `yaml:"http_server"`
}

type Nats struct {
	Host      string
	Port      string
	User      string
	Password  string
	ClusterID string
	StoreType string
	ClientID  string
}

type DB struct {
	Host           string
	Port           string
	User           string
	Password       string
	SSLMode        string
	DBName         string
	MigrationsPath string
}

type HTTPServer struct {
	Address      string        `yaml:"address" env-default:"localhost:8080"`
	Timeout      time.Duration `yaml:"timeout" env-default:"4s"`
	IddleTimeout time.Duration `yaml:"iddle_timeout" env-default:"60s"`
	Port         int           `yaml:"port" env-default:"80"`
}

func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Cannot load the .env file!: %s", err)
	}
	configPath := os.Getenv("CONFIG_PATH")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbSSLMode := os.Getenv("SSL_MODE")
	dbMigrationsPath := os.Getenv("MIGRATIONS_PATH")
	natsHost := os.Getenv("NATS_HOST")
	natsUser := os.Getenv("NATS_USER")
	natsPassword := os.Getenv("NATS_PASSWORD")
	natsPort := os.Getenv("NATS_PORT")
	natsClusterID := os.Getenv("NATS_STREAMING_CLUSTER_ID")
	natsStoreType := os.Getenv("NATS_STREAMING_STORE_TYPE")
	natsClientID := os.Getenv("NATS_CLIENT_ID")
	var nats Nats = Nats{
		Host:      natsHost,
		User:      natsUser,
		Password:  natsPassword,
		Port:      natsPort,
		ClusterID: natsClusterID,
		StoreType: natsStoreType,
		ClientID:  natsClientID,
	}
	var db DB = DB{
		Host:           dbHost,
		Port:           dbPort,
		User:           dbUser,
		Password:       dbPassword,
		SSLMode:        dbSSLMode,
		DBName:         dbName,
		MigrationsPath: dbMigrationsPath,
	}
	if configPath == "" {
		log.Fatal("Config file is not set!")
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Config file doesn't exist!: %s", err)
	}
	var Service Service
	if err := cleanenv.ReadConfig(configPath, &Service); err != nil {
		log.Fatalf("Cannot read the config!: %s", err)
	}
	var cfg Config = Config{
		DB:      db,
		Service: Service,
		Nats:    nats,
	}
	return &cfg
}
