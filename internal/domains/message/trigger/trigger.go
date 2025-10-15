package trigger

import (
	"AddressService/internal/domains/message/model"
	"sync"
)

//
// ===== Реалтайм триггер (AddressTrigger) =====
//

type cachedData struct {
	Pos     model.Pos
	Address string
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

// Реалтайм логика: город → 300м, трасса → 2000м
func (t *AddressTrigger) ShouldUpdateAddress(id int64, newPos model.Pos) (bool, string) {

	t.mu.RLock()
	last, ok := t.lastGeoMap[id]
	t.mu.RUnlock()

	if !ok {
		return true, ""
	}

	dist := DistanceMeters(last.Pos.Y, last.Pos.X, newPos.Y, newPos.X)

	var threshold float64
	threshold = 300 // город
	if newPos.S > 80 {
		threshold = 2000 // трасса
	}

	if dist >= threshold {
		return true, ""
	}

	return false, last.Address
}

func (t *AddressTrigger) UpdateAddress(id int64, pos model.Pos, address string) {
	t.mu.Lock()
	t.lastGeoMap[id] = cachedData{
		Pos:     pos,
		Address: address,
	}
	t.mu.Unlock()
}

//
// ===== Отчётный триггер (ReportAddressTrigger) =====
//

type cachedReportData struct {
	Pos     model.Pos
	Address string
}

type ReportAddressTrigger struct {
	mu         sync.RWMutex
	lastGeoMap map[int64]cachedReportData
}

func NewReportAddressTrigger() *ReportAddressTrigger {
	return &ReportAddressTrigger{
		lastGeoMap: make(map[int64]cachedReportData),
	}
}

// Отчётная логика: всегда 20 метров
func (t *ReportAddressTrigger) ShouldUpdateAddress(id int64, newPos model.Pos) (bool, string) {
	t.mu.RLock()
	last, ok := t.lastGeoMap[id]
	t.mu.RUnlock()

	if !ok {
		return true, ""
	}

	dist := DistanceMeters(last.Pos.Y, last.Pos.X, newPos.Y, newPos.X)

	if dist >= 20 {
		return true, ""
	}

	return false, last.Address
}

func (t *ReportAddressTrigger) UpdateAddress(id int64, pos model.Pos, address string) {
	t.mu.Lock()
	t.lastGeoMap[id] = cachedReportData{
		Pos:     pos,
		Address: address,
	}
	t.mu.Unlock()
}
