package main

import (
	"AddressService/config"
	"AddressService/internal/domains/message/handler/http"
	HandKafka "AddressService/internal/domains/message/handler/kafka"
	"AddressService/internal/domains/message/repository/geocoder"
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

	gin.SetMode(gin.ReleaseMode)

	producer := ProdKafka.NewKafkaProducer(cfg.Kafka.Brokers, cfg.Kafka.EnrichedTopic)

	geo := geocoder.New(cfg.Geocoder.BaseURL, cfg.Geocoder.TimeoutMs, cfg.Geocoder.Workers)

	addressTrigger := trigger.NewAddressTrigger()
	messageUC := usecase.NewMessageUseCase(addressTrigger, producer, geo)

	r := gin.Default()
	httpHandler := http.NewMessageHandler(messageUC)
	r.POST("/message", httpHandler.Handle)
	r.POST("/report", httpHandler.HandleReport)

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers:        cfg.Kafka.Brokers,
		Topic:          cfg.Kafka.RawTopic,
		GroupID:        cfg.Kafka.GroupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		QueueCapacity:  10000,
		MaxWait:        20 * time.Millisecond,
		CommitInterval: 2 * time.Second,
	})

	log.Println("started")
	kafkaConsumer := HandKafka.NewMessageConsumer(messageUC, reader, 200, 50000)
	go kafkaConsumer.Consume(context.Background())

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
