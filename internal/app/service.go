package app

import (
	"context"
	"log"

	"AdressService/internal/domain"
	"AdressService/internal/geocoder"
	"AdressService/internal/kafka"
)

type Service struct {
	Geo      *geocoder.Client
	Producer *kafkaserv.Producer
}

func (s *Service) ProcessMessage(ctx context.Context, msg domain.Message) {
	addr, err := s.Geo.ReverseGeocode(msg.Pos.X, msg.Pos.Y)
	if err != nil {
		log.Printf("Failed to geocode: %v", err)
		return
	}
	msg.Address = addr
	if err := s.Producer.Send(ctx, msg); err != nil {
		log.Printf("Failed to send: %v", err)
	}
}
