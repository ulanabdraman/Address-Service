package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type Config struct {
	KafkaBrokers []string
	InputTopic   string
	OutputTopic  string
	GeoURL       string
}

func LoadConfig() *Config {
	_ = godotenv.Load() // загружаем .env

	cfg := &Config{
		KafkaBrokers: parseList("KAFKA_BROKERS"),
		InputTopic:   mustGet("KAFKA_INPUT_TOPIC"),
		OutputTopic:  mustGet("KAFKA_OUTPUT_TOPIC"),
		GeoURL:       mustGet("GEO_URL"),
	}

	return cfg
}

func mustGet(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		log.Fatalf("⚠️ missing required env variable: %s", key)
	}
	return val
}

func parseList(key string) []string {
	val := mustGet(key)
	parts := strings.Split(val, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
