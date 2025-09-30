package usecase

import (
	"AddressService/internal/domains/message/model"
	"AddressService/internal/domains/message/repository/geocoder"
	"AddressService/internal/domains/message/repository/kafka"
	"AddressService/internal/domains/message/trigger"
	"context"
	"fmt"
)

type MessageUseCase interface {
	ProcessMessage(ctx context.Context, msg *model.Message) error
}

type messageUseCase struct {
	trigger  *trigger.AddressTrigger
	producer kafka.KafkaProducer
}

func NewMessageUseCase(
	trigger *trigger.AddressTrigger,
	producer kafka.KafkaProducer,
) MessageUseCase {
	return &messageUseCase{
		trigger:  trigger,
		producer: producer,
	}
}

func (u *messageUseCase) ProcessMessage(ctx context.Context, msg *model.Message) error {
	shouldGeocode, cachedAddress := u.trigger.ShouldUpdateAddress(msg.ID, msg.Pos)

	if shouldGeocode {
		address, err := geocoder.GetAddress(ctx, msg.Pos)
		if err != nil {
			return err
		}
		msg.Address = fmt.Sprintf("%s, %s, %s, %s, %s, %s", address.Country, address.CityDistrict, address.City, address.Road, address.Neighbourhood, address.Postcode)
		u.trigger.UpdateAddress(msg.ID, msg.Pos, address)
	} else {
		msg.Address = fmt.Sprintf("%s, %s, %s, %s, %s, %s", cachedAddress.Country, cachedAddress.CityDistrict, cachedAddress.City, cachedAddress.Road, cachedAddress.Neighbourhood, cachedAddress.Postcode)
	}

	return u.producer.Produce(ctx, msg)
}
