package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=ICartRepository
type ICartRepository interface {
	GetByUserID(ctx context.Context, userID uint) (*domain.Cart, error)
	Save(ctx context.Context, cart *domain.Cart) (*domain.Cart, error)
	DeleteByID(ctx context.Context, id uint) error
	ClearCartItems(ctx context.Context, cartID uint) error
}

type cartRepository struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewCartRepository(db *gorm.DB, logger zerolog.Logger) ICartRepository {
	return &cartRepository{
		DB:     db,
		Logger: logger,
	}
}

// GetByUserID returns a cart with all items and related product info
func (r *cartRepository) GetByUserID(ctx context.Context, userID uint) (*domain.Cart, error) {
	var cart domain.Cart
	if err := r.DB.WithContext(ctx).
		Preload("Items.Product").
		Where("user_id = ?", userID).
		First(&cart).Error; err != nil {
		r.Logger.Error().Err(err).Uint("user_id", userID).Msg("failed to get cart by user ID")
		return nil, err
	}
	return &cart, nil
}

// Save inserts or updates the cart and its items
func (r *cartRepository) Save(ctx context.Context, cart *domain.Cart) (*domain.Cart, error) {
	if err := r.DB.WithContext(ctx).Save(cart).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to save cart")
		return nil, err
	}

	// Save cart items
	for _, item := range cart.Items {
		item.CartID = cart.ID
		if err := r.DB.WithContext(ctx).Save(&item).Error; err != nil {
			r.Logger.Error().Err(err).Msg("failed to save cart item")
			return nil, err
		}
	}

	return cart, nil
}

// DeleteByID deletes the cart and its items
func (r *cartRepository) DeleteByID(ctx context.Context, id uint) error {
	tx := r.DB.WithContext(ctx).Begin()

	// Delete cart items
	if err := tx.Where("cart_id = ?", id).Delete(&domain.CartItem{}).Error; err != nil {
		tx.Rollback()
		r.Logger.Error().Err(err).Uint("cart_id", id).Msg("failed to delete cart items")
		return err
	}

	// Delete the cart
	if err := tx.Delete(&domain.Cart{}, id).Error; err != nil {
		tx.Rollback()
		r.Logger.Error().Err(err).Uint("cart_id", id).Msg("failed to delete cart")
		return err
	}

	return tx.Commit().Error
}

// ClearCartItems removes all items from a given cart
func (r *cartRepository) ClearCartItems(ctx context.Context, cartID uint) error {
	if err := r.DB.WithContext(ctx).
		Where("cart_id = ?", cartID).
		Delete(&domain.CartItem{}).Error; err != nil {
		r.Logger.Error().Err(err).Uint("cart_id", cartID).Msg("failed to clear cart items")
		return err
	}
	return nil
}
