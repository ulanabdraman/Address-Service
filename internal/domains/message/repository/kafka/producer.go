package kafka

import (
	"AddressService/internal/domains/message/model"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
)

type KafkaProducer interface {
	Produce(ctx context.Context, msg *model.Message) error
}

type kafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaProducer(writer *kafka.Writer, topic string) KafkaProducer {
	return &kafkaProducer{
		writer: writer,
		topic:  topic,
	}
}

func (p *kafkaProducer) Produce(ctx context.Context, msg *model.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	kMsg := kafka.Message{
		Value: data,
	}

	return p.writer.WriteMessages(ctx, kMsg)
}
