package main

import (
	"context"
	"github.com/segmentio/kafka-go"
	"log"

	"AdressService/internal/app"
	"AdressService/internal/config"
	"AdressService/internal/geocoder"
	"AdressService/internal/kafka"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()

	geoClient := &geocoder.Client{BaseURL: cfg.GeoURL}
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: cfg.KafkaBrokers,
		Topic:   cfg.OutputTopic,
	})
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.KafkaBrokers,
		Topic:   cfg.InputTopic,
		GroupID: "geo-service",
	})

	prod := &kafkaserv.Producer{Writer: writer}
	svc := &app.Service{Geo: geoClient, Producer: prod}

	cons := &kafkaserv.Consumer{
		Reader:  reader,
		Handler: svc.ProcessMessage,
	}

	log.Println("Service started")
	cons.Start(ctx)
}
