package repository

import (
	"context"

	"example.com/backend/internal/domain"
	"gorm.io/gorm"
)

//go:generate mockery --name=ITransactionRepository
type ITransactionRepository interface {
	FindAll(ctx context.Context) ([]domain.Transaction, error)
	FindAllByCampaignID(ctx context.Context, campaignID string) ([]domain.Transaction, error)
	FindAllByUserID(ctx context.Context, userID string) ([]domain.Transaction, error)
	FindByID(ctx context.Context, id string) (*domain.Transaction, error)
	FindByReference(ctx context.Context, reference string) (*domain.Transaction, error)
	Create(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error)
	Update(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error)
	DeleteByID(ctx context.Context, id string) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) ITransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) FindAll(ctx context.Context) ([]domain.Transaction, error) {
	var txs []domain.Transaction
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Campaign").
		Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *transactionRepository) FindAllByCampaignID(ctx context.Context, campaignID string) ([]domain.Transaction, error) {
	var txs []domain.Transaction
	if err := r.db.WithContext(ctx).
		Where("campaign_id = ?", campaignID).
		Preload("User").
		Preload("Campaign").
		Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *transactionRepository) FindAllByUserID(ctx context.Context, userID string) ([]domain.Transaction, error) {
	var txs []domain.Transaction
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Campaign").
		Preload("Campaign.CampaignImages").
		Find(&txs).Error; err != nil {
		return nil, err
	}
	return txs, nil
}

func (r *transactionRepository) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	var tx domain.Transaction
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Campaign").
		First(&tx, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) FindByReference(ctx context.Context, reference string) (*domain.Transaction, error) {
	var tx domain.Transaction
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Campaign").
		First(&tx, "reference = ?", reference).Error; err != nil {
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	if err := r.db.WithContext(ctx).Create(tx).Error; err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	if err := r.db.WithContext(ctx).Save(tx).Error; err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *transactionRepository) DeleteByID(ctx context.Context, id string) error {
	if err := r.db.WithContext(ctx).Delete(&domain.Transaction{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
