package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Kafka  KafkaConfig  `mapstructure:"kafka"`
}

type ServerConfig struct {
	Port int `mapstructure:"port"`
}

type KafkaConfig struct {
	Brokers       []string `mapstructure:"brokers"`
	RawTopic      string   `mapstructure:"raw_topic"`
	EnrichedTopic string   `mapstructure:"enriched_topic"`
	GroupID       string   `mapstructure:"group_id"`
}

func LoadConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("kafka.brokers", []string{"localhost:9092"})
	v.SetDefault("kafka.raw_topic", "raw-topic")
	v.SetDefault("kafka.enriched_topic", "enriched-topic")
	v.SetDefault("kafka.group_id", "address-service-group")

	if err := v.ReadInConfig(); err != nil {
		fmt.Println("Config file not found, using defaults and env")
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
