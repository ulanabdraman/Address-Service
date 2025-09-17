package kafkaserv

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"AdressService/internal/domain"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func (p *Producer) Send(ctx context.Context, msg domain.Message) error {
	val, err := json.Marshal(msg)
	log.Println(msg)
	if err != nil {
		return err
	}
	return p.Writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(msg.ID, 10)),
		Value: val,
	})
}
