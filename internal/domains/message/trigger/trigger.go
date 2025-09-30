package trigger

import (
	"AddressService/internal/domains/message/model"
	"sync"
)

type cachedData struct {
	Pos     model.Pos
	Address model.Address
}

type AddressTrigger struct {
	mu         sync.RWMutex
	lastGeoMap map[int64]cachedData
}

func NewAddressTrigger() *AddressTrigger {
	return &AddressTrigger{
		lastGeoMap: make(map[int64]cachedData),
	}
}

// Возвращает: (нужно ли геокодировать, адрес если есть)
func (t *AddressTrigger) ShouldUpdateAddress(id int64, newPos model.Pos) (bool, model.Address) {
	t.mu.RLock()
	last, ok := t.lastGeoMap[id]
	t.mu.RUnlock()

	if !ok {
		return true, model.Address{}
	}

	dist := DistanceMeters(last.Pos.Y, last.Pos.X, newPos.Y, newPos.X)

	// 🧠 Логика приоритета:
	// если в прошлом адресе не было города → трасса → 2 км
	// если город есть → город → 300 м
	var threshold float64
	if last.Address.City == "" {
		threshold = 2000
	} else {
		threshold = 300
	}

	if dist >= threshold {
		return true, model.Address{}
	}

	return false, last.Address
}

func (t *AddressTrigger) UpdateAddress(id int64, pos model.Pos, address model.Address) {
	t.mu.Lock()
	t.lastGeoMap[id] = cachedData{
		Pos:     pos,
		Address: address,
	}
	t.mu.Unlock()
}
