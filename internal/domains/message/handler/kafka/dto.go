package kafka

import "AddressService/internal/domains/message/model"

type MessageDTO struct {
	ID     int64                  `json:"id"`
	DT     int64                  `json:"dt"`
	ST     int64                  `json:"st"`
	Pos    model.Pos              `json:"pos"`
	Params map[string]interface{} `json:"p"`
}

// Конвертация DTO → Model
func (dto *MessageDTO) ToModel() *model.Message {
	return &model.Message{
		ID:     dto.ID,
		DT:     dto.DT,
		ST:     dto.ST,
		Pos:    dto.Pos,
		Params: dto.Params,
	}
}
