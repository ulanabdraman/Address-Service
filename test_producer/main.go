package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type Pos struct {
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
	Z  int     `json:"z"`
	A  int     `json:"a"`
	S  int     `json:"s"`
	St int     `json:"st"`
}

type Message struct {
	T       int64                  `json:"t"`
	ST      int                    `json:"st"`
	Pos     Pos                    `json:"pos"`
	Params  map[string]interface{} `json:"p"`
	Address string                 `json:"address,omitempty"`
}

func main() {
	// Настройки
	broker := getEnv("KAFKA_BROKER", "localhost:9092")
	topic := getEnv("KAFKA_TOPIC", "raw-location")

	// Создаем Kafka writer
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{broker},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	// Пример сообщения
	msg := Message{
		T:  time.Now().Unix(),
		ST: 1,
		Pos: Pos{
			X:  76.909230,
			Y:  43.238949,
			Z:  0,
			A:  0,
			S:  0,
			St: 1,
		},
		Params: map[string]interface{}{
			"speed": 42,
			"bat":   88.5,
		},
	}

	// Сериализация
	value, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("marshal error: %v", err)
	}

	// Отправка
	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("test-key"),
		Value: value,
	})
	if err != nil {
		log.Fatalf("send error: %v", err)
	}

	log.Println("✅ Тестовое сообщение отправлено в Kafka")
}

// getEnv возвращает переменную окружения или значение по умолчанию
func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
