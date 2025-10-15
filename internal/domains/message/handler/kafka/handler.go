package kafka

import (
	"AddressService/internal/domains/message/model"
	"AddressService/internal/domains/message/usecase"
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/segmentio/kafka-go"
)

var json = jsoniter.ConfigFastest

type MessageConsumer struct {
	usecase     usecase.MessageUseCase
	reader      *kafka.Reader
	workerCount int
	queue       chan model.Message
	total       atomic.Int64
	batchSize   int
	wg          sync.WaitGroup
}

func NewMessageConsumer(uc usecase.MessageUseCase, reader *kafka.Reader, workers int, queueSize int) *MessageConsumer {
	if workers <= 0 {
		workers = 200
	}
	if queueSize <= 0 {
		queueSize = 50_000
	}

	c := &MessageConsumer{
		usecase:     uc,
		reader:      reader,
		workerCount: workers,
		queue:       make(chan model.Message, queueSize),
		batchSize:   500,
	}

	for i := 0; i < workers; i++ {
		c.wg.Add(1)
		go c.worker()
	}

	return c
}

func (c *MessageConsumer) worker() {
	defer c.wg.Done()
	for msg := range c.queue {
		if err := c.usecase.ProcessMessage(context.Background(), &msg); err != nil {
			log.Printf("❌ worker: %v", err)
		}
		c.total.Add(1)
	}
}

func (c *MessageConsumer) Consume(ctx context.Context) error {
	defer close(c.queue)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			batch := make([]kafka.Message, 0, c.batchSize)
			for i := 0; i < c.batchSize; i++ {
				m, err := c.reader.FetchMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return nil
					}
					log.Printf("❌ FetchMessage: %v", err)
					break
				}
				batch = append(batch, m)
			}

			if len(batch) == 0 {
				continue
			}

			// 🧠 Decode batch без лишних аллокаций
			for _, km := range batch {
				var raw []MessageDTO
				if err := json.Unmarshal(km.Value, &raw); err != nil {
					continue
				}
				for _, dto := range raw {
					msg := dto.ToModel()
					select {
					case c.queue <- *msg:
					default:
						// бэкпрешер: блокируем на долю секунды, если очередь заполнена
						time.Sleep(10 * time.Millisecond)
						c.queue <- *msg
					}
				}
			}

			// ✅ Commit одной пачкой
			if err := c.reader.CommitMessages(ctx, batch...); err != nil {
				log.Printf("⚠️ Commit failed: %v", err)
			}
		}
	}
}

func (c *MessageConsumer) Close() {
	close(c.queue)
	c.wg.Wait()
}
