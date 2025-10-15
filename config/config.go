package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Kafka    KafkaConfig    `mapstructure:"kafka"`
	Geocoder GeocoderConfig `mapstructure:"geocoder"`
}

type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type KafkaConfig struct {
	Brokers       []string `mapstructure:"brokers"`
	RawTopic      string   `mapstructure:"raw_topic"`
	EnrichedTopic string   `mapstructure:"enriched_topic"`
	GroupID       string   `mapstructure:"group_id"`
}

type GeocoderConfig struct {
	BaseURL   string `mapstructure:"base_url"`
	TimeoutMs int    `mapstructure:"timeout_ms"`
	Workers   int    `mapstructure:"workers"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.host", "localhost")

	v.SetDefault("kafka.brokers", []string{"localhost:9092"})
	v.SetDefault("kafka.raw_topic", "raw-topic")
	v.SetDefault("kafka.enriched_topic", "enriched-topic")
	v.SetDefault("kafka.group_id", "address-service-group")

	v.SetDefault("geocoder.base_url", "http://localhost:8012")
	v.SetDefault("geocoder.timeout_ms", 800)
	v.SetDefault("geocoder.workers", 100)

	if err := v.ReadInConfig(); err != nil {
		fmt.Println("⚠️  Config file not found, using defaults and env")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
