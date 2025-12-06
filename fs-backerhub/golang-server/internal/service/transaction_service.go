package service

import (
	"context"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/model"
	"example.com.backend/internal/repos"
	"example.com.backend/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ITransactionService interface {
	FindAll(ctx context.Context) ([]model.TransactionDTO, error)
	FindByID(ctx context.Context, id string) (*model.TransactionDTO, error)
	Create(ctx context.Context, transactionDto model.TransactionDTO) error
	Update(ctx context.Context, transactionDto model.TransactionDTO) error
	DeleteByID(ctx context.Context, id string) error
}

type transactionService struct {
	transactionRepo repos.ITransactionRepository
}

func NewTransactionService(
	transactionRepo repos.ITransactionRepository,
) ITransactionService {
	return &transactionService{transactionRepo: transactionRepo}
}

func (s *transactionService) FindAll(ctx context.Context) ([]model.TransactionDTO, error) {
	log := logger.GetLoggerFromContext(ctx)

	transactions, err := s.transactionRepo.FindAll(ctx)
	if err != nil {
		log.Error("failed to retrieve transactions", zap.Error(err))
		return nil, err
	}

	var transactionDtos []model.TransactionDTO
	for _, transaction := range transactions {
		var dto model.TransactionDTO
		dto.FromTransaction(&transaction)
		transactionDtos = append(transactionDtos, dto)
	}

	log.Info("successfully retrieved all transactions", zap.Int("count", len(transactions)))
	return transactionDtos, nil
}

func (s *transactionService) FindByID(ctx context.Context, id string) (*model.TransactionDTO, error) {
	log := logger.GetLoggerFromContext(ctx)

	transaction, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("failed to retrieve transaction by id", zap.String("transaction_id", id), zap.Error(err))
		return nil, err
	}
	if transaction == nil {
		log.Warn("transaction not found", zap.String("transaction_id", id))
		return nil, nil
	}

	var dto model.TransactionDTO
	dto.FromTransaction(transaction)

	log.Info("successfully retrieved transaction", zap.String("transaction_id", id))
	return &dto, nil
}

func (s *transactionService) Create(ctx context.Context, transactionDto model.TransactionDTO) error {
	log := logger.GetLoggerFromContext(ctx)

	transaction := &domain.Transaction{
		Amount:     transactionDto.Amount,
		UserID:     transactionDto.UserID,
		CampaignID: transactionDto.CampaignID,
		Status:     string(domain.StatusPending),
		Note:       "process-by-admin",
		Reference:  uuid.New().String(),
	}

	_, err := s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		log.Error("failed to create transaction",
			zap.String("user_id", transactionDto.UserID),
			zap.String("campaign_id", transactionDto.CampaignID),
			zap.Error(err),
		)
		return err
	}

	log.Info("successfully created transaction",
		zap.String("transaction_id", transaction.ID),
		zap.String("transaction_reference", transaction.Reference),
	)
	return nil
}

func (s *transactionService) Update(ctx context.Context, transactionDto model.TransactionDTO) error {
	log := logger.GetLoggerFromContext(ctx)

	transaction, err := s.transactionRepo.FindByID(ctx, transactionDto.ID)
	if err != nil {
		log.Error("failed to retrieve transaction for update",
			zap.String("transaction_id", transactionDto.ID),
			zap.Error(err),
		)
		return err
	}
	if transaction == nil {
		log.Warn("transaction not found for update", zap.String("transaction_id", transactionDto.ID))
		return nil
	}

	transaction.Amount = transactionDto.Amount
	transaction.Status = string(transactionDto.Status)

	_, err = s.transactionRepo.Update(ctx, transaction)
	if err != nil {
		log.Error("failed to update transaction",
			zap.String("transaction_id", transactionDto.ID),
			zap.Error(err),
		)
		return err
	}

	log.Info("successfully updated transaction", zap.String("transaction_id", transactionDto.ID))
	return nil
}

func (s *transactionService) DeleteByID(ctx context.Context, id string) error {
	log := logger.GetLoggerFromContext(ctx)

	err := s.transactionRepo.DeleteByID(ctx, id)
	if err != nil {
		log.Error("failed to delete transaction", zap.String("transaction_id", id), zap.Error(err))
		return err
	}

	log.Info("successfully deleted transaction", zap.String("transaction_id", id))
	return nil
}
