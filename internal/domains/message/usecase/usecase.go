package usecase

import (
	"AddressService/internal/domains/message/model"
	"AddressService/internal/domains/message/repository/geocoder"
	"AddressService/internal/domains/message/repository/kafka"
	"AddressService/internal/domains/message/trigger"
	"context"
	"sync"
	"time"
)

type MessageUseCase interface {
	ProcessMessage(ctx context.Context, msg *model.Message) error
	ProcessMessages(ctx context.Context, msgs []*model.Message) ([]*model.Message, error)
	Close()
}

type messageUseCase struct {
	trigger      *trigger.AddressTrigger
	producer     kafka.KafkaProducer
	geocoder     *geocoder.Geocoder // 👈 передаётся извне
	geoQueue     chan *model.Message
	produceQueue chan *model.Message
	wg           sync.WaitGroup
	stopCh       chan struct{}

	batchSize   int
	batchWait   time.Duration
	geoParallel int
}

// 👇 теперь принимаем готовый geocoder
func NewMessageUseCase(trigger *trigger.AddressTrigger, producer kafka.KafkaProducer, geo *geocoder.Geocoder) MessageUseCase {
	u := &messageUseCase{
		trigger:      trigger,
		producer:     producer,
		geocoder:     geo, // 👈 сохраняем сюда
		geoQueue:     make(chan *model.Message, 10_000),
		produceQueue: make(chan *model.Message, 10_000),
		stopCh:       make(chan struct{}),

		batchSize:   100,
		batchWait:   100 * time.Millisecond,
		geoParallel: 10,
	}

	// геокодер pool
	for i := 0; i < u.geoParallel; i++ {
		u.wg.Add(1)
		go u.geoWorkerBatch()
	}

	// продюсер pool
	u.wg.Add(1)
	go u.produceWorker()

	return u
}

func (u *messageUseCase) Close() {
	close(u.stopCh)
	close(u.geoQueue)
	u.wg.Wait()
	close(u.produceQueue)
	_ = u.producer.Close()
}

// ----------- GEOCODER WORKER (BATCH) -----------

func (u *messageUseCase) geoWorkerBatch() {
	defer u.wg.Done()

	ticker := time.NewTicker(u.batchWait)
	defer ticker.Stop()

	batch := make([]*model.Message, 0, u.batchSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}

		positions := make([]model.Pos, len(batch))
		for i, m := range batch {
			positions[i] = m.Pos
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		addrs, err := u.geocoder.GetAddresses(ctx, positions) // 👈 теперь через экземпляр
		cancel()

		if err != nil {
			println("❌ geoWorkerBatch: geocoder error:", err.Error())
			batch = batch[:0]
			return
		}

		for i, m := range batch {
			addr := ""
			if i < len(addrs) {
				addr = addrs[i]
			}
			m.Address = addr
			u.trigger.UpdateAddress(m.ID, m.Pos, addr)

			select {
			case u.produceQueue <- m:
			case <-u.stopCh:
				return
			}
		}

		batch = batch[:0]
	}

	for {
		select {
		case <-u.stopCh:
			flush()
			return

		case msg, ok := <-u.geoQueue:
			if !ok {
				flush()
				return
			}
			batch = append(batch, msg)
			if len(batch) >= u.batchSize {
				flush()
			}

		case <-ticker.C:
			flush()
		}
	}
}

// ----------- PRODUCER WORKER -----------

func (u *messageUseCase) produceWorker() {
	defer u.wg.Done()

	const batchSize = 500
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	batch := make([]*model.Message, 0, batchSize)

	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := u.producer.ProduceBatch(context.Background(), batch); err != nil {
			println("⚠️ produceWorker: batch produce error:", err.Error())
		}
		batch = batch[:0]
	}

	for {
		select {
		case <-u.stopCh:
			for m := range u.produceQueue {
				batch = append(batch, m)
				if len(batch) >= batchSize {
					flush()
				}
			}
			flush()
			return

		case m, ok := <-u.produceQueue:
			if !ok {
				flush()
				return
			}
			batch = append(batch, m)
			if len(batch) >= batchSize {
				flush()
			}

		case <-ticker.C:
			flush()
		}
	}
}

// ----------- ENTRY POINTS -----------

func (u *messageUseCase) ProcessMessage(ctx context.Context, msg *model.Message) error {
	shouldGeocode, cached := u.trigger.ShouldUpdateAddress(msg.ID, msg.Pos)

	if shouldGeocode {
		local := *msg
		select {
		case u.geoQueue <- &local:
			return nil
		case <-u.stopCh:
			return context.Canceled
		default:
			local.Address = cached
			return u.producer.Produce(ctx, &local)
		}
	}

	local := *msg
	local.Address = cached
	return u.producer.Produce(ctx, &local)
}

func (u *messageUseCase) ProcessMessages(ctx context.Context, msgs []*model.Message) ([]*model.Message, error) {
	results := make([]*model.Message, 0, len(msgs))

	toGeocode := make([]*model.Message, 0, len(msgs))
	positions := make([]model.Pos, 0, len(msgs))

	for _, msg := range msgs {
		local := *msg
		shouldGeocode, cached := u.trigger.ShouldUpdateAddress(local.ID, local.Pos)
		if shouldGeocode {
			toGeocode = append(toGeocode, &local)
			positions = append(positions, local.Pos)
		} else {
			local.Address = cached
			results = append(results, &local)
		}
	}

	if len(toGeocode) == 0 {
		return results, nil
	}

	addrs, err := u.geocoder.GetAddresses(ctx, positions) // 👈 тоже через u.geocoder
	if err != nil {
		return results, err
	}

	for i, msg := range toGeocode {
		addr := ""
		if i < len(addrs) {
			addr = addrs[i]
		}
		msg.Address = addr
		u.trigger.UpdateAddress(msg.ID, msg.Pos, addr)
		results = append(results, msg)
	}

	return results, nil
}
