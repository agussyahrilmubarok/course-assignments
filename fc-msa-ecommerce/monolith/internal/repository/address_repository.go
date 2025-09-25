package repository

import (
	"context"
	"ecommerce/internal/domain"

	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

//go:generate mockery --name=IAddressRepository
type IAddressRepository interface {
	FindAllByUserID(ctx context.Context, userID uint) ([]domain.Address, error)
	FindByID(ctx context.Context, id uint) (*domain.Address, error)
	Save(ctx context.Context, address *domain.Address) (*domain.Address, error)
	Update(ctx context.Context, address *domain.Address) (*domain.Address, error)
	DeleteByID(ctx context.Context, id uint) error
	SetDefaultAddress(ctx context.Context, userID uint, addressID uint) error
}

type addressRepository struct {
	DB     *gorm.DB
	Logger zerolog.Logger
}

func NewAddressRepository(db *gorm.DB, logger zerolog.Logger) IAddressRepository {
	return &addressRepository{
		DB:     db,
		Logger: logger,
	}
}

func (r *addressRepository) FindAllByUserID(ctx context.Context, userID uint) ([]domain.Address, error) {
	var addresses []domain.Address
	if err := r.DB.WithContext(ctx).Where("user_id = ?", userID).Find(&addresses).Error; err != nil {
		r.Logger.Error().Err(err).Uint("user_id", userID).Msg("failed to fetch addresses")
		return nil, err
	}
	return addresses, nil
}

func (r *addressRepository) FindByID(ctx context.Context, id uint) (*domain.Address, error) {
	var address domain.Address
	if err := r.DB.WithContext(ctx).First(&address, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("address_id", id).Msg("address not found")
		return nil, err
	}
	return &address, nil
}

func (r *addressRepository) Save(ctx context.Context, address *domain.Address) (*domain.Address, error) {
	if err := r.DB.WithContext(ctx).Create(address).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to save address")
		return nil, err
	}
	return address, nil
}

func (r *addressRepository) Update(ctx context.Context, address *domain.Address) (*domain.Address, error) {
	if err := r.DB.WithContext(ctx).Save(address).Error; err != nil {
		r.Logger.Error().Err(err).Msg("failed to update address")
		return nil, err
	}
	return address, nil
}

func (r *addressRepository) DeleteByID(ctx context.Context, id uint) error {
	if err := r.DB.WithContext(ctx).Delete(&domain.Address{}, id).Error; err != nil {
		r.Logger.Error().Err(err).Uint("address_id", id).Msg("failed to delete address")
		return err
	}
	return nil
}

// SetDefaultAddress sets one address as default for a user, unsets others
func (r *addressRepository) SetDefaultAddress(ctx context.Context, userID uint, addressID uint) error {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Unset all other default addresses for the user
	if err := tx.Model(&domain.Address{}).
		Where("user_id = ?", userID).
		Update("is_default", false).Error; err != nil {
		tx.Rollback()
		r.Logger.Error().Err(err).Uint("user_id", userID).Msg("failed to unset default addresses")
		return err
	}

	// Set the specified address as default
	if err := tx.Model(&domain.Address{}).
		Where("id = ? AND user_id = ?", addressID, userID).
		Update("is_default", true).Error; err != nil {
		tx.Rollback()
		r.Logger.Error().Err(err).Uint("address_id", addressID).Msg("failed to set default address")
		return err
	}

	return tx.Commit().Error
}
