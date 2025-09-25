package repository

import (
	"context"
	"ecommerce/internal/domain"
	"time"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=IPaymentRepository
type IPaymentRepository interface {
	FindByID(ctx context.Context, id uint) (*domain.Payment, error)
	FindByOrderID(ctx context.Context, orderID uint) (*domain.Payment, error)
	Save(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	Update(ctx context.Context, payment *domain.Payment) (*domain.Payment, error)
	DeleteByID(ctx context.Context, id uint) error
}

type paymentRepository struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewPaymentRepository(db *gorm.DB, logger zerolog.Logger) IPaymentRepository {
	return &paymentRepository{
		DB:     db,
		Logger: logger,
	}
}

func (r *paymentRepository) FindByID(ctx context.Context, id uint) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.DB.WithContext(ctx).First(&payment, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("payment_id", id).Msg("payment not found")
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) FindByOrderID(ctx context.Context, orderID uint) (*domain.Payment, error) {
	var payment domain.Payment
	if err := r.DB.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error; err != nil {
		r.Logger.Error().Err(err).Uint("order_id", orderID).Msg("payment for order not found")
		return nil, err
	}
	return &payment, nil
}

func (r *paymentRepository) Save(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	if err := r.DB.WithContext(ctx).Create(payment).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to save payment")
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) Update(ctx context.Context, payment *domain.Payment) (*domain.Payment, error) {
	payment.UpdatedAt = time.Now()
	if err := r.DB.WithContext(ctx).Save(payment).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to update payment")
		return nil, err
	}
	return payment, nil
}

func (r *paymentRepository) DeleteByID(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&domain.Payment{}, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("payment_id", id).Msg("failed to delete payment")
		return err
	}
	return nil
}
