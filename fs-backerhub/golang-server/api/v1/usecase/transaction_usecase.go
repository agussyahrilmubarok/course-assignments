package usecaseV1

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	payloadV1 "example.com/backend/api/v1/payload"
	"example.com/backend/internal/domain"
	"example.com/backend/internal/exception"
	"example.com/backend/internal/model"
	"example.com/backend/internal/repository"
	"example.com/backend/internal/service"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=ITransactionUseCaseV1
type ITransactionUseCaseV1 interface {
	FindAllByUser(ctx context.Context, userID string) ([]payloadV1.TransactionResponse, error)
	FindAllByCampaign(ctx context.Context, campaignID string, userID string) ([]payloadV1.TransactionResponse, error)
	FindByID(ctx context.Context, txID string, userID string) (*payloadV1.TransactionResponse, error)
	Create(ctx context.Context, amount float64, campaignID string, userID string) (string, string, error)
	ProcessMidtransPayment(ctx context.Context, callback *model.MidtransCallback, status domain.TransactionStatus) error
}

type transactionUseCaseV1 struct {
	transactionRepo repository.ITransactionRepository
	midtransService service.IMidtransService
	userRepo        repository.IUserRepository
	campaignRepo    repository.ICampaignRepository
	log             zerolog.Logger
}

func NewTransactionUseCaseV1(
	transactionRepo repository.ITransactionRepository,
	midtransService service.IMidtransService,
	userRepo repository.IUserRepository,
	campaignRepo repository.ICampaignRepository,
	log zerolog.Logger,
) ITransactionUseCaseV1 {
	return &transactionUseCaseV1{
		transactionRepo: transactionRepo,
		midtransService: midtransService,
		userRepo:        userRepo,
		campaignRepo:    campaignRepo,
		log:             log,
	}
}

func (uc *transactionUseCaseV1) FindAllByUser(ctx context.Context, userID string) ([]payloadV1.TransactionResponse, error) {
	txs, err := uc.transactionRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		uc.log.Warn().Msg("transactions not found")
		return nil, exception.NewNotFound("Transactions not found", err)
	}

	var resps []payloadV1.TransactionResponse
	for _, tx := range txs {
		var resp payloadV1.TransactionResponse
		resp.FromTransaction(&tx)
		resps = append(resps, resp)
	}

	return resps, nil
}

func (uc *transactionUseCaseV1) FindAllByCampaign(ctx context.Context, campaignID string, userID string) ([]payloadV1.TransactionResponse, error) {
	campaign, err := uc.campaignRepo.FindByID(ctx, campaignID)
	if err != nil || campaign == nil {
		uc.log.Warn().Msg("campaign not found")
		return nil, exception.NewBadRequest("Campaign not found", err)
	}

	// check valid user
	if campaign.UserID != userID {
		uc.log.Warn().Msg("do not have permissions")
		return nil, exception.NewUnauthorized("Do not have permission", nil)
	}

	txs, err := uc.transactionRepo.FindAllByCampaignID(ctx, campaignID)
	if err != nil {
		uc.log.Warn().Msg("transactions not found")
		return nil, exception.NewNotFound("Transactions not found", err)
	}

	var resps []payloadV1.TransactionResponse
	for _, tx := range txs {
		var resp payloadV1.TransactionResponse
		resp.FromTransaction(&tx)
		resps = append(resps, resp)
	}

	return resps, nil
}

func (uc *transactionUseCaseV1) FindByID(ctx context.Context, txID string, userID string) (*payloadV1.TransactionResponse, error) {
	tx, err := uc.transactionRepo.FindByID(ctx, txID)
	if err != nil || tx == nil {
		uc.log.Warn().Msg("transactions not found")
		return nil, exception.NewNotFound("Transactions not found", err)
	}

	campaign, err := uc.campaignRepo.FindByID(ctx, tx.CampaignID)
	if err != nil || campaign == nil {
		uc.log.Warn().Msg("campaign not found")
		return nil, exception.NewBadRequest("Campaign not found", err)
	}

	// check valid user
	if campaign.UserID != userID {
		uc.log.Warn().Msg("do not have permissions")
		return nil, exception.NewUnauthorized("Do not have permission", nil)
	}

	var resp payloadV1.TransactionResponse
	resp.FromTransaction(tx)

	return &resp, nil
}

func (uc *transactionUseCaseV1) Create(ctx context.Context, amount float64, campaignID string, userID string) (string, string, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		uc.log.Warn().Msg("user not found")
		return "", "", exception.NewBadRequest("User not found", err)
	}

	campaign, err := uc.campaignRepo.FindByID(ctx, campaignID)
	if err != nil || campaign == nil {
		uc.log.Warn().Msg("campaign not found")
		return "", "", exception.NewBadRequest("Campaign not found", err)
	}

	ref := fmt.Sprintf("%s-%s-%s-%v", uuid.New().String(), campaignID, userID, time.Now().Unix())
	tx := &domain.Transaction{
		Amount:     amount,
		Status:     string(domain.StatusPending),
		Note:       "tx-status-getting-payment-url",
		Reference:  ref,
		CampaignID: campaign.ID,
		UserID:     user.ID,
	}

	tx, err = uc.transactionRepo.Create(ctx, tx)
	if err != nil {
		uc.log.Error().Err(err).Msg("failed when create transaction")
		return "", "", exception.NewInternal("Transaction campaign fail", err)
	}

	var midtransReq model.MidtransRequest
	midtransReq.TransactionDetails.OrderID = tx.ID
	midtransReq.TransactionDetails.GrossAmt = int64(amount)
	midtransReq.CustomerDetail.FName = user.Name
	midtransReq.CustomerDetail.Email = user.Email
	paymentUrl, err := uc.midtransService.CreateTransaction(midtransReq)
	if err != nil {
		tx.Status = string(domain.StatusFailed)
		uc.transactionRepo.Update(ctx, tx)
		uc.log.Error().Err(err).Msg("failed when create transaction via midtrans")
		return "", "", exception.NewInternal("Transaction campaign fail", err)
	}

	return tx.ID, paymentUrl, nil
}

func (uc *transactionUseCaseV1) ProcessMidtransPayment(ctx context.Context, callback *model.MidtransCallback, status domain.TransactionStatus) error {
	// Retrieve internal transaction using order_id from Midtrans
	tx, err := uc.transactionRepo.FindByID(ctx, callback.OrderID)
	if err != nil || tx == nil {
		uc.log.Error().Err(err).Str("order_id", callback.OrderID).Msg("transation not found")
		return errors.New("failed to process transaction payment")
	}

	// Retrieve associated campaign
	campaign, err := uc.campaignRepo.FindByID(ctx, tx.CampaignID)
	if err != nil || campaign == nil {
		uc.log.Error().Err(err).Msg("campaign not found")
		return errors.New("failed to process transaction payment")
	}

	switch status {
	case domain.StatusPaid:
		uc.log.Error().Msgf("processing paid transaction: %s", tx.ID)

		// Parse transaction_time
		transactionTime, err := time.Parse("2006-01-02 15:04:05", callback.TransactionTime)
		if err != nil {
			uc.log.Error().Err(err).Msgf("failed to parse transaction_time '%s': %v", callback.TransactionTime, err)
			return errors.New("failed to process transaction payment")
		}

		// Parse gross_amount
		grossAmount, err := strconv.ParseFloat(callback.GrossAmount, 64)
		if err != nil {
			uc.log.Error().Err(err).Msgf("failed to parse gross_amount '%s': %v", callback.GrossAmount, err)
			return errors.New("failed to process transaction payment")
		}

		// Update transaction
		tx.Status = string(status)
		tx.Note = fmt.Sprintf("tx,paidAt-%v,method-%v,id-%v", transactionTime, callback.PaymentType, callback.TransactionID)
		tx, err = uc.transactionRepo.Update(ctx, tx)
		if err != nil {
			uc.log.Error().Err(err).Msgf("failed to update transaction %s: %v", tx.ID, err)
			return errors.New("failed to process transaction payment")
		}

		// Update campaign
		campaign.BackerCount += 1
		campaign.CurrentAmount += grossAmount

		if _, err := uc.campaignRepo.Update(ctx, campaign); err != nil {
			// TODO: rollback
			uc.log.Error().Err(err).Msgf("failed to update campaign %s after transaction %s: %v", campaign.ID, tx.ID, err)
			return errors.New("failed to process transaction payment")
		}

	case domain.StatusPending:
		uc.log.Error().Err(err).Msgf("processing pending transaction: %s", tx.ID)

	case domain.StatusFailed:
		uc.log.Error().Err(err).Msgf("processing failed transaction: %s", tx.ID)

		// Parse transaction_time
		transactionTime, err := time.Parse("2006-01-02 15:04:05", callback.TransactionTime)
		if err != nil {
			uc.log.Error().Err(err).Msgf("failed to parse transaction_time '%s': %v", callback.TransactionTime, err)
			return errors.New("failed to process transaction payment")
		}

		// Update transaction
		tx.Status = string(status)
		tx.Note = fmt.Sprintf("tx,failedAt-%v,method-%v,id-%v", transactionTime, callback.PaymentType, callback.TransactionID)
		tx, err = uc.transactionRepo.Update(ctx, tx)
		if err != nil {
			uc.log.Error().Err(err).Msgf("failed to update transaction %s: %v", tx.ID, err)
			return errors.New("failed to process transaction payment")
		}

	case domain.StatusCanceled:
		uc.log.Error().Err(err).Msgf("processing canceled transaction: %s", tx.ID)

		// Parse transaction_time
		transactionTime, err := time.Parse("2006-01-02 15:04:05", callback.TransactionTime)
		if err != nil {
			uc.log.Error().Err(err).Msgf("failed to parse transaction_time '%s': %v", callback.TransactionTime, err)
			return errors.New("failed to process transaction payment")
		}

		// Parse gross_amount
		grossAmount, err := strconv.ParseFloat(callback.GrossAmount, 64)
		if err != nil {
			uc.log.Error().Err(err).Msgf("failed to parse gross_amount '%s': %v", callback.GrossAmount, err)
			return errors.New("failed to process transaction payment")
		}

		// Update transaction
		tx.Status = string(status)
		tx.Note = fmt.Sprintf("tx,cancelledAt-%v,method-%v,id-%v", transactionTime, callback.PaymentType, callback.TransactionID)
		tx, err = uc.transactionRepo.Update(ctx, tx)
		if err != nil {
			uc.log.Error().Err(err).Msgf("failed to update transaction %s: %v", tx.ID, err)
			return errors.New("failed to process transaction payment")
		}

		// Update campaign
		campaign.BackerCount -= 1
		campaign.CurrentAmount -= grossAmount

		if _, err := uc.campaignRepo.Update(ctx, campaign); err != nil {
			// TODO: rollback
			uc.log.Error().Err(err).Msgf("failed to update campaign %s after transaction %s: %v", campaign.ID, tx.ID, err)
			return errors.New("failed to process transaction payment")
		}

	case domain.StatusExpired:
		uc.log.Error().Err(err).Msgf("processing paid transaction: %s", tx.ID)

		// Parse transaction_time
		transactionTime, err := time.Parse("2006-01-02 15:04:05", callback.TransactionTime)
		if err != nil {
			uc.log.Error().Err(err).Msgf("failed to parse transaction_time '%s': %v", callback.TransactionTime, err)
			return errors.New("failed to process transaction payment")
		}

		// Update transaction
		tx.Status = string(status)
		tx.Note = fmt.Sprintf("tx,expiredAt-%v,method-%v,id-%v", transactionTime, callback.PaymentType, callback.TransactionID)
		tx, err = uc.transactionRepo.Update(ctx, tx)
		if err != nil {
			uc.log.Error().Err(err).Msgf("failed to update transaction %s: %v", tx.ID, err)
			return errors.New("failed to process transaction payment")
		}

	default:
		uc.log.Error().Err(err).Msgf("processing unknown transaction: %s", tx.ID)
	}

	return nil
}
