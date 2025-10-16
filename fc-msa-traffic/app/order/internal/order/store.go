package order

import (
	"context"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=IStore
type IStore interface {
	SaveOrder(ctx context.Context, order *Order) error
}

type store struct {
	db  *gorm.DB
	log zerolog.Logger
}

func NewStore(db *gorm.DB, log zerolog.Logger) IStore {
	return &store{db: db, log: log}
}

func (s *store) SaveOrder(ctx context.Context, order *Order) error {
	for i := range order.OrderItems {
		order.OrderItems[i].TotalPrice = order.OrderItems[i].CalculateTotalPrice()
	}

	order.FinalAmount = order.CalculateFinalAmount()

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			s.log.Error().Err(err).Msg("Failed to create order")
			return err
		}

		for _, item := range order.OrderItems {
			item.OrderID = order.ID
			if err := tx.Create(&item).Error; err != nil {
				s.log.Error().Err(err).Msg("Failed to create order item")
				return err
			}
		}

		return nil
	})

	if err != nil {
		s.log.Error().Err(err).Msg("Transaction failed when saving order")
		return err
	}

	return nil
}
