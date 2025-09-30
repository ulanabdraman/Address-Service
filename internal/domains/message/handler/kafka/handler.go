package kafka

import (
	"AddressService/internal/domains/message/usecase"
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
)

type MessageConsumer struct {
	usecase usecase.MessageUseCase
}

func NewMessageConsumer(uc usecase.MessageUseCase) *MessageConsumer {
	return &MessageConsumer{usecase: uc}
}

func (h *MessageConsumer) Consume(ctx context.Context, r *kafka.Reader) {
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("kafka read error: %v", err)
			continue
		}

		var raw []MessageDTO
		if err := json.Unmarshal(m.Value, &raw); err != nil {
			log.Printf("invalid message array format: %v", err)
			continue
		}

		for _, dto := range raw {
			msg := dto.ToModel()
			if err := h.usecase.ProcessMessage(ctx, msg); err != nil {
				log.Printf("failed to process message: %v", err)
			}
		}
	}
}
