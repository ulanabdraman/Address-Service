package kafkaserv

import (
	"AdressService/internal/domain"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type Consumer struct {
	Reader  *kafka.Reader
	Handler func(context.Context, domain.Message)
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		m, err := c.Reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Kafka read error: %v", err)
			continue
		}
		var msg domain.Message
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			log.Printf("Invalid message: %v", err)
			continue
		}
		c.Handler(ctx, msg)
	}
}
