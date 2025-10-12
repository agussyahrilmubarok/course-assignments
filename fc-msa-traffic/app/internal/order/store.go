package order

import (
	"context"

	"gorm.io/gorm"
)

//go:generate mockery --name=IOrderStore
type IOrderStore interface {
	FindByUserID(ctx context.Context, userID string) ([]Order, error)
	FindByStatus(ctx context.Context, status Status) ([]Order, error)
	FindByID(ctx context.Context, orderID string) (*Order, error)
	Save(ctx context.Context, order *Order) error
	UpdateStatus(ctx context.Context, orderID string, status Status) error
	DeleteByID(ctx context.Context, orderID string) error
}

type orderStore struct {
	db *gorm.DB
}

func NewOrderStore(db *gorm.DB) IOrderStore {
	return &orderStore{db: db}
}

func (s *orderStore) FindByUserID(ctx context.Context, userID string) ([]Order, error) {
	var orders []Order
	if err := s.db.WithContext(ctx).
		Preload("Items").
		Where("user_id = ?", userID).
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *orderStore) FindByStatus(ctx context.Context, status Status) ([]Order, error) {
	var orders []Order
	if err := s.db.WithContext(ctx).
		Preload("Items").
		Where("status = ?", status).
		Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *orderStore) FindByID(ctx context.Context, orderID string) (*Order, error) {
	var o Order
	if err := s.db.WithContext(ctx).
		Preload("Items").
		First(&o, "id = ?", orderID).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (s *orderStore) Save(ctx context.Context, order *Order) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(order).Error; err != nil {
			return err
		}

		for _, item := range order.Items {
			item.OrderID = order.ID
			if err := tx.Save(&item).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *orderStore) UpdateStatus(ctx context.Context, orderID string, status Status) error {
	return s.db.WithContext(ctx).
		Model(&Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error
}

func (s *orderStore) DeleteByID(ctx context.Context, orderID string) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&OrderItem{}, "order_id = ?", orderID).Error; err != nil {
			return err
		}
		return tx.Delete(&Order{}, "id = ?", orderID).Error
	})
}
