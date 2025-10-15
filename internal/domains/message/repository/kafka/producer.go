package kafka

import (
	"AddressService/internal/domains/message/model"
	"context"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
)

var json = jsoniter.ConfigFastest

type KafkaProducer interface {
	Produce(ctx context.Context, msg *model.Message) error
	ProduceBatch(ctx context.Context, msgs []*model.Message) error
	Close() error
}

type kafkaProducer struct {
	writer *kafka.Writer
	topic  string
}

func NewKafkaProducer(brokers []string, topic string) KafkaProducer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.CRC32Balancer{},
		Async:                  true, // пусть Writer сам ассинчит
		BatchSize:              512,
		BatchTimeout:           3 * time.Millisecond,
		RequiredAcks:           kafka.RequireOne,
		Compression:            kafka.Snappy,
		AllowAutoTopicCreation: true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				println("❌ Kafka write error:", err.Error())
			}
		},
	}

	return &kafkaProducer{
		writer: writer,
		topic:  topic,
	}
}

// отправка одного сообщения
func (p *kafkaProducer) Produce(_ context.Context, msg *model.Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		println("❌ KafkaProducer: Marshal error:", err.Error())
		return err
	}

	return p.writer.WriteMessages(context.Background(), kafka.Message{
		Value: data,
	})
}

// батч-отправка (вызов горутиной — отлично)
func (p *kafkaProducer) ProduceBatch(_ context.Context, msgs []*model.Message) error {
	if len(msgs) == 0 {
		return nil
	}

	kmsgs := make([]kafka.Message, 0, len(msgs))
	for _, m := range msgs {
		data, err := json.Marshal(m)
		if err != nil {
			println("⚠️ KafkaProducer: skip bad message:", err.Error())
			continue
		}
		kmsgs = append(kmsgs, kafka.Message{Value: data})
	}

	// Writer сам разобьёт на внутренние пакеты по BatchSize/BatchTimeout
	return p.writer.WriteMessages(context.Background(), kmsgs...)
}

func (p *kafkaProducer) Close() error {
	return p.writer.Close()
}
