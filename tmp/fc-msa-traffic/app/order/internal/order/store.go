package order

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=IStore
type IStore interface {
	SaveOrder(ctx context.Context, order *Order) error
	FindOrderByID(ctx context.Context, orderID string) (*Order, error)
	UpdateStatus(ctx context.Context, order *Order, status OrderStatus) error
}

type store struct {
	shards []*gorm.DB
	router *ShardRouter
	log    zerolog.Logger
}

func NewStore(shards []*gorm.DB, router *ShardRouter, log zerolog.Logger) IStore {
	return &store{
		shards: shards,
		router: router,
		log:    log,
	}
}

func (s *store) getDBByUserID(userID string) *gorm.DB {
	index := s.router.GetShard(userID)
	return s.shards[index]
}

func (s *store) getDBByOrderID(orderID string) *gorm.DB {
	index := s.router.GetShard(orderID)
	return s.shards[index]
}

func (s *store) SaveOrder(ctx context.Context, order *Order) error {
	if order.ID == "" {
		order.ID = uuid.NewString()
	}

	for i := range order.OrderItems {
		if order.OrderItems[i].ID == "" {
			order.OrderItems[i].ID = uuid.NewString()
		}
		order.OrderItems[i].OrderID = order.ID
		order.OrderItems[i].TotalPrice = order.OrderItems[i].CalculateTotalPrice()
	}

	order.FinalAmount = order.CalculateFinalAmount()

	db := s.getDBByUserID(order.UserID)

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			s.log.Error().Err(err).
				Str("order_id", order.ID).
				Msg("Failed to create order and its items")
			return err
		}
		return nil
	})

	if err != nil {
		s.log.Error().Err(err).Str("order_id", order.ID).Msg("Transaction failed when saving order")
		return err
	}

	s.log.Info().Str("order_id", order.ID).Msg("Save order successfully")
	return nil
}

func (s *store) FindOrderByID(ctx context.Context, orderID string) (*Order, error) {
	db := s.getDBByOrderID(orderID)

	var order Order
	err := db.WithContext(ctx).
		Preload("OrderItems").
		First(&order, "id = ?", orderID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		s.log.Warn().Str("order_id", orderID).Msg("Order not found")
		return nil, nil
	}

	if err != nil {
		s.log.Error().Err(err).Str("order_id", orderID).Msg("Failed to find order by ID")
		return nil, err
	}

	s.log.Info().Str("order_id", orderID).Msg("Order found successfully")
	return &order, nil
}

func (s *store) UpdateStatus(ctx context.Context, order *Order, status OrderStatus) error {
	if order == nil || order.ID == "" {
		s.log.Error().Msg("Order object or ID is required to update status")
		return errors.New("order object or ID is required")
	}

	db := s.getDBByOrderID(order.ID)

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing Order
		if err := tx.First(&existing, "id = ?", order.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				s.log.Warn().Str("order_id", order.ID).Msg("Order not found when updating status")
				return gorm.ErrRecordNotFound
			}
			s.log.Error().Err(err).Str("order_id", order.ID).Msg("Failed to find order before updating status")
			return err
		}

		order.Status = status
		if err := tx.Model(&existing).Update("status", status).Error; err != nil {
			s.log.Error().Err(err).Str("order_id", order.ID).Msg("Failed to update order status")
			return err
		}
		return nil
	})

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		s.log.Error().Err(err).Str("order_id", order.ID).Msg("Transaction failed when updating order status")
		return err
	}

	s.log.Info().
		Str("order_id", order.ID).
		Str("new_status", string(status)).
		Msg("Order status updated successfully")

	return nil
}
