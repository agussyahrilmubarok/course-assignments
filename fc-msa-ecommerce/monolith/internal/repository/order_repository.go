package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=IOrderRepository
type IOrderRepository interface {
	FindAllByUserID(ctx context.Context, userID uint) ([]domain.Order, error)
	FindByID(ctx context.Context, id uint) (*domain.Order, error)
	Save(ctx context.Context, order *domain.Order) (*domain.Order, error)
	UpdateStatus(ctx context.Context, orderID uint, status string) error
	DeleteByID(ctx context.Context, id uint) error
}

type orderRepository struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewOrderRepository(db *gorm.DB, logger zerolog.Logger) IOrderRepository {
	return &orderRepository{
		DB:     db,
		Logger: logger,
	}
}

func (r *orderRepository) FindAllByUserID(ctx context.Context, userID uint) ([]domain.Order, error) {
	var orders []domain.Order
	err := r.DB.WithContext(ctx).
		Preload("Items.Product").
		Preload("Address").
		Preload("Payment").
		Where("user_id = ?", userID).
		Find(&orders).Error

	if err != nil {
		r.Logger.Error().Err(err).Uint("user_id", userID).Msg("failed to fetch orders for user")
		return nil, err
	}

	return orders, nil
}

func (r *orderRepository) FindByID(ctx context.Context, id uint) (*domain.Order, error) {
	var order domain.Order
	err := r.DB.WithContext(ctx).
		Preload("Items.Product").
		Preload("Address").
		Preload("Payment").
		First(&order, id).Error

	if err != nil {
		r.Logger.Error().Err(err).Uint("order_id", id).Msg("failed to fetch order by id")
		return nil, err
	}

	return &order, nil
}

func (r *orderRepository) Save(ctx context.Context, order *domain.Order) (*domain.Order, error) {
	tx := r.DB.WithContext(ctx).Begin()

	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		r.Logger.Error().Err(err).Msg("failed to create order")
		return nil, err
	}

	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		if err := tx.Create(&order.Items[i]).Error; err != nil {
			tx.Rollback()
			r.Logger.Error().Err(err).Msg("failed to create order item")
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to commit order transaction")
		return nil, err
	}

	return order, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID uint, status string) error {
	err := r.DB.WithContext(ctx).
		Model(&domain.Order{}).
		Where("id = ?", orderID).
		Update("status", status).Error

	if err != nil {
		r.Logger.Error().Err(err).Uint("order_id", orderID).Msg("failed to update order status")
		return err
	}

	return nil
}

func (r *orderRepository) DeleteByID(ctx context.Context, id uint) error {
	tx := r.DB.WithContext(ctx).Begin()

	// Delete order items first
	if err := tx.Where("order_id = ?", id).Delete(&domain.OrderItem{}).Error; err != nil {
		tx.Rollback()
		r.Logger.Error().Err(err).Uint("order_id", id).Msg("failed to delete order items")
		return err
	}

	// Then delete the order
	if err := tx.Delete(&domain.Order{}, id).Error; err != nil {
		tx.Rollback()
		r.Logger.Error().Err(err).Uint("order_id", id).Msg("failed to delete order")
		return err
	}

	return tx.Commit().Error
}
