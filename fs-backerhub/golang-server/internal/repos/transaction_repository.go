package repos

import (
	"context"

	"example.com.backend/internal/domain"
	"example.com.backend/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

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
	log := logger.GetLoggerFromContext(ctx)

	var txs []domain.Transaction
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Campaign").
		Find(&txs).Error; err != nil {

		log.Error("failed fetching all transactions", zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched all transactions", zap.Int("count", len(txs)))
	return txs, nil
}

func (r *transactionRepository) FindAllByCampaignID(ctx context.Context, campaignID string) ([]domain.Transaction, error) {
	log := logger.GetLoggerFromContext(ctx)

	var txs []domain.Transaction
	if err := r.db.WithContext(ctx).
		Where("campaign_id = ?", campaignID).
		Preload("User").
		Preload("Campaign").
		Find(&txs).Error; err != nil {

		log.Error("failed fetching transactions by campaign id",
			zap.String("campaign_id", campaignID),
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("successfully fetched transactions by campaign id",
		zap.String("campaign_id", campaignID),
		zap.Int("count", len(txs)),
	)
	return txs, nil
}

func (r *transactionRepository) FindAllByUserID(ctx context.Context, userID string) ([]domain.Transaction, error) {
	log := logger.GetLoggerFromContext(ctx)

	var txs []domain.Transaction
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Preload("User").
		Preload("Campaign").
		Preload("Campaign.CampaignImages").
		Find(&txs).Error; err != nil {

		log.Error("failed fetching transactions by user id",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("successfully fetched transactions by user id",
		zap.String("user_id", userID),
		zap.Int("count", len(txs)),
	)
	return txs, nil
}

func (r *transactionRepository) FindByID(ctx context.Context, id string) (*domain.Transaction, error) {
	log := logger.GetLoggerFromContext(ctx)

	var tx domain.Transaction
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Campaign").
		First(&tx, "id = ?", id).Error; err != nil {

		log.Error("failed fetching transaction by id", zap.String("transaction_id", id), zap.Error(err))
		return nil, err
	}

	log.Info("successfully fetched transaction by id", zap.String("transaction_id", id))
	return &tx, nil
}

func (r *transactionRepository) FindByReference(ctx context.Context, reference string) (*domain.Transaction, error) {
	log := logger.GetLoggerFromContext(ctx)

	var tx domain.Transaction
	if err := r.db.WithContext(ctx).
		Preload("User").
		Preload("Campaign").
		First(&tx, "reference = ?", reference).Error; err != nil {

		log.Error("failed fetching transaction by reference",
			zap.String("transaction_reference", reference),
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("successfully fetched transaction by reference", zap.String("transaction_reference", reference))
	return &tx, nil
}

func (r *transactionRepository) Create(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Create(tx).Error; err != nil {
		log.Error("failed creating transaction",
			zap.String("user_id", tx.UserID),
			zap.String("campaign_id", tx.CampaignID),
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("successfully created transaction",
		zap.String("id", tx.ID),
		zap.String("user_id", tx.UserID),
		zap.String("campaign_id", tx.CampaignID),
	)
	return tx, nil
}

func (r *transactionRepository) Update(ctx context.Context, tx *domain.Transaction) (*domain.Transaction, error) {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Save(tx).Error; err != nil {
		log.Error("failed updating transaction",
			zap.String("transaction_id", tx.ID),
			zap.Error(err),
		)
		return nil, err
	}

	log.Info("successfully updated transaction", zap.String("transaction_id", tx.ID))
	return tx, nil
}

func (r *transactionRepository) DeleteByID(ctx context.Context, id string) error {
	log := logger.GetLoggerFromContext(ctx)

	if err := r.db.WithContext(ctx).Delete(&domain.Transaction{}, "id = ?", id).Error; err != nil {
		log.Error("failed deleting transaction", zap.String("transaction_id", id), zap.Error(err))
		return err
	}

	log.Info("successfully deleted transaction", zap.String("transaction_id", id))
	return nil
}
