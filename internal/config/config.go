package config

import (
	"log"
	"os"

	// "github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

const (
	configPathEnvKey 	= "CONFIG_PATH"
	postgresURLEnvKey 	= "POSTGRES_URL"
	valkeyURLEnvKey 	= "VALKEY_URL"
)

// Config represents the configuration structure
type Config struct {
	Database   		StorageConfig	
	Valkey     		RedisConfig
}

type StorageConfig struct {
	URL 			string
}

type RedisConfig struct {
	URL				string
}

// MustLoadConfig loads the configuration from the specified path
func MustLoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("no .env file found, falling back to environment only")
	}
	
	configPath := os.Getenv(configPathEnvKey)
	if configPath == "" {
		log.Fatalf("%s is not set up", configPathEnvKey)
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist: %s", configPath, err.Error())
	}

	var cfg Config
	// if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
	// 	log.Fatalf("failed to read config file: %s", err.Error())
	// }

	postgresURL := os.Getenv(postgresURLEnvKey)
	valkeyURL := os.Getenv(valkeyURLEnvKey)

	if postgresURL == "" || valkeyURL == "" {
		log.Fatalf("failed to read URLs")
	}

	cfg.Database.URL = postgresURL
	cfg.Valkey.URL = valkeyURL

	return &cfg
}
