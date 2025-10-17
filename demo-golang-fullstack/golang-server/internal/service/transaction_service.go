package service

import (
	"context"

	"example.com/backend/internal/domain"
	"example.com/backend/internal/model"
	"example.com/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=ITransactionService
type ITransactionService interface {
	FindAll(ctx context.Context) ([]model.TransactionDTO, error)
	FindByID(ctx context.Context, id string) (*model.TransactionDTO, error)
	Create(ctx context.Context, transactionDto model.TransactionDTO) error
	Update(ctx context.Context, transactionDto model.TransactionDTO) error
	DeleteByID(ctx context.Context, id string) error
}

type transactionService struct {
	transactionRepo repository.ITransactionRepository
	log             zerolog.Logger
}

func NewTransactionService(
	transactionRepo repository.ITransactionRepository,
	log zerolog.Logger,
) ITransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		log:             log,
	}
}

func (s *transactionService) FindAll(ctx context.Context) ([]model.TransactionDTO, error) {
	transactions, err := s.transactionRepo.FindAll(ctx)
	if err != nil {
		s.log.Error().Err(err).Msgf("failed to retrieve transactions")
		return nil, err
	}

	var transationDtos []model.TransactionDTO
	for _, transaction := range transactions {
		var transactionDto model.TransactionDTO
		transactionDto.FromTransaction(&transaction)
		transationDtos = append(transationDtos, transactionDto)
	}

	return transationDtos, nil
}

func (s *transactionService) FindByID(ctx context.Context, id string) (*model.TransactionDTO, error) {
	transaction, err := s.transactionRepo.FindByID(ctx, id)
	if err != nil || transaction == nil {
		s.log.Error().Err(err).Msgf("failed to find transaction id %s", id)
		return nil, err
	}

	var transactionDto model.TransactionDTO
	transactionDto.FromTransaction(transaction)

	return &transactionDto, nil
}

func (s *transactionService) Create(ctx context.Context, transactionDto model.TransactionDTO) error {
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
		s.log.Error().Err(err).Msgf("failed to create transaction")
		return err
	}

	return nil
}

func (s *transactionService) Update(ctx context.Context, transactionDto model.TransactionDTO) error {
	transaction, err := s.transactionRepo.FindByID(ctx, transactionDto.ID)
	if err != nil || transaction == nil {
		s.log.Error().Err(err).Msgf("failed to find transaction id %s", transactionDto.ID)
		return err
	}

	transaction.Amount = transactionDto.Amount
	transaction.Status = string(transactionDto.Status)

	_, err = s.transactionRepo.Update(ctx, transaction)
	if err != nil {
		s.log.Error().Err(err).Msgf("failed to update transaction")
		return err
	}

	return nil
}

func (s *transactionService) DeleteByID(ctx context.Context, id string) error {
	if err := s.transactionRepo.DeleteByID(ctx, id); err != nil {
		s.log.Error().Err(err).Msgf("failed to delete transaction id %s", id)
		return err
	}

	return nil
}
