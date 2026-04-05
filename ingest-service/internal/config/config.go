package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Host string
	Port string
}

type KafkaConfig struct {
	Brokers []string
}

type PostgresConfig struct {
	ConnectionString string
}

type Config struct {
	Server   ServerConfig
	Kafka    KafkaConfig
	Postgres PostgresConfig
}

func LoadConfig() (*Config, error) {
	godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			Host: os.Getenv("SERVER_HOST"),
			Port: os.Getenv("SERVER_PORT"),
		},
		Kafka: KafkaConfig{
			Brokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		},
		Postgres: PostgresConfig{
			ConnectionString: os.Getenv("DB_CONN"),
		},
	}

	return config, nil
}
