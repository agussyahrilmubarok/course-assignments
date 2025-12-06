package usecaseV1

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/model"
	"example.com.backend/internal/repos"
	"example.com.backend/internal/service"
	"example.com.backend/pkg/exception"
	"example.com.backend/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"

	payloadV1 "example.com.backend/internal/rest/v1/payload"
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
	transactionRepo repos.ITransactionRepository
	midtransService service.IMidtransService
	userRepo        repos.IUserRepository
	campaignRepo    repos.ICampaignRepository
}

func NewTransactionUseCaseV1(
	transactionRepo repos.ITransactionRepository,
	midtransService service.IMidtransService,
	userRepo repos.IUserRepository,
	campaignRepo repos.ICampaignRepository,
) ITransactionUseCaseV1 {
	return &transactionUseCaseV1{
		transactionRepo: transactionRepo,
		midtransService: midtransService,
		userRepo:        userRepo,
		campaignRepo:    campaignRepo,
	}
}

func (uc *transactionUseCaseV1) FindAllByUser(ctx context.Context, userID string) ([]payloadV1.TransactionResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	txs, err := uc.transactionRepo.FindAllByUserID(ctx, userID)
	if err != nil {
		log.Error("failed fetching transactions for user", zap.String("user_id", userID), zap.Error(err))
		return nil, exception.NewNotFound("Transactions not found", err)
	}

	resps := make([]payloadV1.TransactionResponse, 0, len(txs))
	for _, tx := range txs {
		var resp payloadV1.TransactionResponse
		resp.FromTransaction(&tx)
		resps = append(resps, resp)
	}

	log.Info("successfully fetched transactions for user", zap.String("user_id", userID), zap.Int("count", len(resps)))
	return resps, nil
}

func (uc *transactionUseCaseV1) FindAllByCampaign(ctx context.Context, campaignID string, userID string) ([]payloadV1.TransactionResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	campaign, err := uc.campaignRepo.FindByID(ctx, campaignID)
	if err != nil {
		log.Error("failed fetching campaign", zap.String("campaign_id", campaignID), zap.Error(err))
		return nil, exception.NewBadRequest("Campaign not found", err)
	}
	if campaign == nil {
		log.Warn("campaign not found", zap.String("campaign_id", campaignID))
		return nil, exception.NewBadRequest("Campaign not found", nil)
	}

	if campaign.UserID != userID {
		log.Warn("unauthorized access to campaign transactions", zap.String("campaign_id", campaignID), zap.String("user_id", userID))
		return nil, exception.NewUnauthorized("Do not have permission", nil)
	}

	txs, err := uc.transactionRepo.FindAllByCampaignID(ctx, campaignID)
	if err != nil {
		log.Error("failed fetching transactions for campaign", zap.String("campaign_id", campaignID), zap.Error(err))
		return nil, exception.NewNotFound("Transactions not found", err)
	}

	resps := make([]payloadV1.TransactionResponse, 0, len(txs))
	for _, tx := range txs {
		var resp payloadV1.TransactionResponse
		resp.FromTransaction(&tx)
		resps = append(resps, resp)
	}

	log.Info("successfully fetched transactions for campaign", zap.String("campaign_id", campaignID), zap.Int("count", len(resps)))
	return resps, nil
}

func (uc *transactionUseCaseV1) FindByID(ctx context.Context, txID string, userID string) (*payloadV1.TransactionResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	tx, err := uc.transactionRepo.FindByID(ctx, txID)
	if err != nil || tx == nil {
		log.Warn("transaction not found", zap.String("transaction_id", txID), zap.Error(err))
		return nil, exception.NewNotFound("Transaction not found", err)
	}

	campaign, err := uc.campaignRepo.FindByID(ctx, tx.CampaignID)
	if err != nil || campaign == nil {
		log.Error("campaign not found for transaction", zap.String("transaction_id", txID), zap.String("campaign_id", tx.CampaignID), zap.Error(err))
		return nil, exception.NewBadRequest("Campaign not found", err)
	}

	if campaign.UserID != userID {
		log.Warn("unauthorized access to transaction", zap.String("transaction_id", txID), zap.String("user_id", userID))
		return nil, exception.NewUnauthorized("Do not have permission", nil)
	}

	var resp payloadV1.TransactionResponse
	resp.FromTransaction(tx)

	log.Info("successfully fetched transaction", zap.String("transaction_id", txID))
	return &resp, nil
}

func (uc *transactionUseCaseV1) Create(ctx context.Context, amount float64, campaignID string, userID string) (string, string, error) {
	log := logger.GetLoggerFromContext(ctx)

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		log.Warn("user not found", zap.String("user_id", userID), zap.Error(err))
		return "", "", exception.NewBadRequest("User not found", err)
	}

	campaign, err := uc.campaignRepo.FindByID(ctx, campaignID)
	if err != nil || campaign == nil {
		log.Warn("campaign not found", zap.String("campaign_id", campaignID), zap.Error(err))
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
		log.Error("failed creating transaction", zap.String("user_id", userID), zap.String("campaign_id", campaignID), zap.Error(err))
		return "", "", exception.NewInternal("Transaction creation failed", err)
	}

	log.Info("transaction created", zap.String("transaction_id", tx.ID), zap.String("status", tx.Status), zap.Float64("amount", tx.Amount))

	var midtransReq model.MidtransRequest
	midtransReq.TransactionDetails.OrderID = tx.ID
	midtransReq.TransactionDetails.GrossAmt = int64(amount)
	midtransReq.CustomerDetail.FName = user.Name
	midtransReq.CustomerDetail.Email = user.Email

	paymentUrl, err := uc.midtransService.CreateTransaction(ctx, midtransReq)
	if err != nil {
		tx.Status = string(domain.StatusFailed)
		uc.transactionRepo.Update(ctx, tx)
		log.Error("failed creating transaction via Midtrans", zap.String("transaction_id", tx.ID), zap.String("user_id", userID), zap.String("campaign_id", campaignID), zap.Error(err))
		return "", "", exception.NewInternal("Transaction payment creation failed", err)
	}

	log.Info("successfully created transaction and generated payment URL",
		zap.String("transaction_id", tx.ID),
		zap.String("user_id", userID),
		zap.String("campaign_id", campaignID),
		zap.String("payment_url", paymentUrl),
	)
	return tx.ID, paymentUrl, nil
}

func (uc *transactionUseCaseV1) ProcessMidtransPayment(ctx context.Context, callback *model.MidtransCallback, status domain.TransactionStatus) error {
	log := logger.GetLoggerFromContext(ctx)

	// Retrieve internal transaction using order_id from Midtrans
	tx, err := uc.transactionRepo.FindByID(ctx, callback.OrderID)
	if err != nil || tx == nil {
		log.Error("transaction not found for midtrans callback", zap.String("order_id", callback.OrderID), zap.Error(err))
		return errors.New("failed to process transaction payment")
	}

	// Retrieve associated campaign
	campaign, err := uc.campaignRepo.FindByID(ctx, tx.CampaignID)
	if err != nil || campaign == nil {
		log.Error("campaign not found for midtrans transaction", zap.String("transaction_id", tx.ID), zap.Error(err))
		return errors.New("failed to process transaction payment")
	}

	// Parse transaction_time
	transactionTime, err := time.Parse("2006-01-02 15:04:05", callback.TransactionTime)
	if err != nil {
		log.Error("failed parsing transaction_time", zap.String("transaction_time", callback.TransactionTime), zap.String("transaction_id", tx.ID), zap.Error(err))
		return errors.New("failed to process transaction payment")
	}

	// Parse gross_amount if needed
	var grossAmount float64
	if status == domain.StatusPaid || status == domain.StatusCanceled {
		grossAmount, err = strconv.ParseFloat(callback.GrossAmount, 64)
		if err != nil {
			log.Error("failed parsing gross_amount", zap.String("gross_amount", callback.GrossAmount), zap.String("transaction_id", tx.ID), zap.Error(err))
			return errors.New("failed to process transaction payment")
		}
	}

	// Update transaction
	tx.Status = string(status)
	switch status {
	case domain.StatusPaid:
		tx.Note = fmt.Sprintf("tx,paidAt-%v,method-%v,id-%v", transactionTime, callback.PaymentType, callback.TransactionID)
	case domain.StatusFailed:
		tx.Note = fmt.Sprintf("tx,failedAt-%v,method-%v,id-%v", transactionTime, callback.PaymentType, callback.TransactionID)
	case domain.StatusCanceled:
		tx.Note = fmt.Sprintf("tx,cancelledAt-%v,method-%v,id-%v", transactionTime, callback.PaymentType, callback.TransactionID)
	case domain.StatusExpired:
		tx.Note = fmt.Sprintf("tx,expiredAt-%v,method-%v,id-%v", transactionTime, callback.PaymentType, callback.TransactionID)
	default:
		log.Warn("processing unknown transaction status", zap.String("transaction_id", tx.ID), zap.String("status", string(status)))
		tx.Note = fmt.Sprintf("tx,unknownStatus-%v,method-%v,id-%v", transactionTime, callback.PaymentType, callback.TransactionID)
	}

	tx, err = uc.transactionRepo.Update(ctx, tx)
	if err != nil {
		log.Error("failed updating transaction", zap.String("transaction_id", tx.ID), zap.String("status", string(status)), zap.Error(err))
		return errors.New("failed to process transaction payment")
	}

	// Update campaign if status affects amounts
	switch status {
	case domain.StatusPaid:
		campaign.BackerCount += 1
		campaign.CurrentAmount += grossAmount
	case domain.StatusCanceled:
		campaign.BackerCount -= 1
		campaign.CurrentAmount -= grossAmount
	}

	if status == domain.StatusPaid || status == domain.StatusCanceled {
		_, err := uc.campaignRepo.Update(ctx, campaign)
		if err != nil {
			log.Error("failed updating campaign after transaction", zap.String("campaign_id", campaign.ID), zap.String("transaction_id", tx.ID), zap.Error(err))
			// TODO: implement rollback if needed
			return errors.New("failed to process transaction payment")
		}
	}

	log.Info("successfully processed transaction callback",
		zap.String("transaction_id", tx.ID),
		zap.String("campaign_id", campaign.ID),
		zap.String("status", string(status)),
		zap.Time("transaction_time", transactionTime),
	)
	return nil
}
