package main

import (
	"AddressService/config"
	"AddressService/internal/domains/message/handler/http"
	HandKafka "AddressService/internal/domains/message/handler/kafka"
	ProdKafka "AddressService/internal/domains/message/repository/kafka"
	"AddressService/internal/domains/message/trigger"
	"AddressService/internal/domains/message/usecase"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	kafkago "github.com/segmentio/kafka-go"
	"log"
	"time"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Kafka writer
	kafkaWriter := kafkago.NewWriter(kafkago.WriterConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   cfg.Kafka.EnrichedTopic,
		Async:   true,
	})
	producer := ProdKafka.NewKafkaProducer(kafkaWriter, cfg.Kafka.EnrichedTopic)

	// Usecase
	addressTrigger := trigger.NewAddressTrigger()
	messageUC := usecase.NewMessageUseCase(addressTrigger, producer)

	// HTTP
	r := gin.Default()
	httpHandler := http.NewMessageHandler(messageUC)
	r.POST("/message", httpHandler.Handle)

	// Kafka consumer
	fmt.Printf("[CONFIG DEBUG] RawTopic: '%s'\n", cfg)

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:  cfg.Kafka.Brokers,
		GroupID:  cfg.Kafka.GroupID,
		Topic:    cfg.Kafka.RawTopic,
		MinBytes: 10e3,
		MaxBytes: 10e6,
		MaxWait:  1 * time.Second,
	})
	kafkaConsumer := HandKafka.NewMessageConsumer(messageUC)

	go kafkaConsumer.Consume(context.Background(), reader)

	// Start HTTP
	addr := ":" + fmt.Sprint(cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
	select {}
}
