package order

import (
	"context"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=IService
type IService interface {
	CalculateAndCreateOrder(ctx context.Context, param CreateOrderRequest) (*Order, error)
}

type service struct {
	cfg *
	store IStore
	log   zerolog.Logger
}

func NewService(
	store IStore,
	log zerolog.Logger,
) IService {
	return &service{
		store: store,
		log:   log,
	}
}

func (s *service) CalculateAndCreateOrder(ctx context.Context, param CreateOrderRequest) (*Order, error) {
	panic("unimplemented")
}
