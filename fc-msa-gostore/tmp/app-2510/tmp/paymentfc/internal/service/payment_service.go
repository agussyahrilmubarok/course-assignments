package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"example.com/paymentfc/internal/client"
	"example.com/paymentfc/internal/producer"
	"example.com/paymentfc/internal/store"
	"example.com/paymentfc/pkg/utils"
	"example.com/pkg/model"
	"github.com/google/uuid"
	"github.com/xendit/xendit-go/v7"
	"github.com/xendit/xendit-go/v7/invoice"
	"go.uber.org/zap"
)

//go:generate mockery --name=IPaymentService
type IPaymentService interface {
	CreateInvoiceOrder(ctx context.Context, req *model.OrderModel) error
	ProcessXenditWebhook(ctx context.Context, externalID string, status string, amount float64) error
	ProcessPaymentDirectly(ctx context.Context, req *model.OrderModel) error
	DownloadInvoiceInPDF(ctx context.Context, orderID string) (string, error)
}

type paymentService struct {
	xenditClient    *xendit.APIClient
	paymentStore    store.IPaymentStore
	invoiceStore    store.IInvoiceStore
	userGrpcClient  *client.UserGrpcClient
	paymentProducer producer.KafkaProducer
	log             *zap.Logger
}

func NewPaymentService(
	xenditApiClient string,
	paymentStore store.IPaymentStore,
	invoiceStore store.IInvoiceStore,
	userGrpcClient *client.UserGrpcClient,
	paymentProducer producer.KafkaProducer,
	log *zap.Logger,
) IPaymentService {
	xenditClient := xendit.NewClient(xenditApiClient)
	return &paymentService{
		xenditClient:    xenditClient,
		paymentStore:    paymentStore,
		invoiceStore:    invoiceStore,
		userGrpcClient:  userGrpcClient,
		paymentProducer: paymentProducer,
		log:             log,
	}
}

func (s *paymentService) CreateInvoiceOrder(ctx context.Context, req *model.OrderModel) error {
	externalID := fmt.Sprintf("order-%s", req.ID)
	amount := req.TotalAmount

	// Step 1: Retrieve user email via gRPC to attach to the invoice
	userEmail, err := s.userGrpcClient.GetUserEmail(ctx, req.UserID)
	if err != nil {
		s.log.Error("failed to get user email", zap.String("user_id", req.UserID), zap.Error(err))
		return fmt.Errorf("get user email failed: %w", err)
	}

	// Step 2: Prepare invoice creation request for Xendit
	createInvoiceRequest := *invoice.NewCreateInvoiceRequest(externalID, amount)
	createInvoiceRequest.PayerEmail = &userEmail

	// Step 3: Execute invoice creation using Xendit API
	resp, r, errSDK := s.xenditClient.InvoiceApi.CreateInvoice(ctx).
		CreateInvoiceRequest(createInvoiceRequest).
		Execute()

	if errSDK != nil {
		s.log.Error("error when calling InvoiceApi.CreateInvoice",
			zap.Error(errSDK),
			zap.ByteString("full_error_struct", func() []byte {
				b, _ := json.Marshal(errSDK.FullError())
				return b
			}()),
			zap.Any("http_response", r),
		)
	}

	// Step 4: Persist the invoice record to local invoice store
	if _, err := s.invoiceStore.Create(ctx, &store.Invoice{
		ID:         uuid.New().String(),
		InvoiceURL: resp.GetInvoiceUrl(),
		Invoice:    resp,
	}); err != nil {
		s.log.Error("failed to save xendit invoice in local store", zap.String("order_id", req.ID), zap.Error(err))
		return fmt.Errorf("save invoice to db failed: %w", err)
	}

	// Step 5: Create a payment record based on the invoice
	payment := &store.Payment{
		ID:         uuid.NewString(),
		RefCode:    externalID,
		UserID:     req.UserID,
		OrderID:    req.ID,
		Amount:     resp.GetAmount(),
		Status:     store.StatusPending,
		InvoiceURL: resp.GetInvoiceUrl(),
		ExpiredAt:  resp.GetExpiryDate(),
	}

	createdPayment, err := s.paymentStore.Create(ctx, payment)
	if err != nil {
		s.log.Error("failed to save payment record", zap.String("order_id", req.ID), zap.Error(err))
		return fmt.Errorf("save payment failed: %w", err)
	}

	if createdPayment == nil {
		s.log.Error("payment store returned nil", zap.String("order_id", req.ID))
		return fmt.Errorf("payment store returned nil")
	}

	// Final success log
	s.log.Info("successfully created invoice order",
		zap.String("order_id", req.ID),
		zap.Any("invoice_id", resp.Id),
	)

	return nil
}

func (s *paymentService) ProcessXenditWebhook(ctx context.Context, externalID string, status string, amount float64) error {
	orderID := s.orderIDFromExternalID(externalID)
	payment, err := s.paymentStore.FindByOrderID(ctx, orderID)
	if err != nil || payment == nil {
		s.log.Error("failed to find payment by order id", zap.String("order_id", orderID), zap.Error(err))
		return err
	}

	if amount < payment.Amount || amount <= 0 {
		s.log.Error("amount paid is less than expected",
			zap.Float64("expected_amount", payment.Amount),
			zap.Float64("received_amount", amount),
			zap.String("order_id", orderID),
		)
		return errors.New("insufficient amount received")
	}

	switch status {
	case "PAID":
		payment.PaidAt = time.Now()
		payment.Status = store.StatusCompleted
		_, err := s.paymentStore.UpdateByID(ctx, payment.ID, payment)
		if err != nil {
			s.log.Error("failed to update payment status to completed",
				zap.String("payment_id", payment.ID),
				zap.Error(err),
			)
			return err
		}

		paymentModel := &model.PaymentModel{
			ID:         payment.ID,
			RefCode:    payment.RefCode,
			UserID:     payment.UserID,
			OrderID:    payment.OrderID,
			Amount:     payment.Amount,
			Status:     string(payment.Status),
			InvoiceUrl: payment.InvoiceURL,
			PaidAt:     payment.PaidAt,
		}
		if err := s.paymentProducer.PublicPaymentSuccess(ctx, paymentModel); err != nil {
			s.log.Error("failed to publish payment success", zap.String("order_id", paymentModel.OrderID), zap.Error(err))
			return err
		}

		s.log.Info("payment completed successfully", zap.String("order_id", orderID))

	case "FAILED":
		payment, err := s.paymentStore.UpdateStatus(ctx, payment.ID, store.StatusFailed)
		if err != nil {
			s.log.Error("failed to update payment status to failed",
				zap.String("payment_id", payment.ID),
				zap.Error(err),
			)
			return err
		}

		paymentModel := &model.PaymentModel{
			ID:         payment.ID,
			RefCode:    payment.RefCode,
			UserID:     payment.UserID,
			OrderID:    payment.OrderID,
			Amount:     payment.Amount,
			Status:     string(payment.Status),
			InvoiceUrl: payment.InvoiceURL,
		}
		if err := s.paymentProducer.PublicPaymentFailed(ctx, paymentModel); err != nil {
			s.log.Error("failed to publish payment failed", zap.String("order_id", paymentModel.OrderID), zap.Error(err))
			return err
		}
		s.log.Warn("payment failed", zap.String("order_id", orderID))

	case "PENDING":
		s.log.Info("payment is still pending", zap.String("order_id", orderID))
	default:
		s.log.Warn("received unknown payment status from Xendit",
			zap.String("status", status),
			zap.String("order_id", orderID),
		)
	}

	return nil
}

func (s *paymentService) ProcessPaymentDirectly(ctx context.Context, req *model.OrderModel) error {
	orderID := s.orderIDFromExternalID(req.ID)
	payment, err := s.paymentStore.FindByOrderID(ctx, orderID)
	if err != nil || payment == nil {
		s.log.Error("failed to find payment by order id", zap.String("order_id", orderID), zap.Error(err))
		return err
	}

	if req.TotalAmount < payment.Amount || req.TotalAmount <= 0 {
		s.log.Error("amount paid is less than expected",
			zap.Float64("expected_amount", payment.Amount),
			zap.Float64("received_amount", req.TotalAmount),
			zap.String("order_id", orderID),
		)
		return errors.New("insufficient amount received")
	}

	payment.PaidAt = time.Now()
	payment.Status = store.StatusCompleted
	_, err = s.paymentStore.UpdateByID(ctx, payment.ID, payment)
	if err != nil {
		s.log.Error("failed to update payment status to completed",
			zap.String("payment_id", payment.ID),
			zap.Error(err),
		)
		return err
	}

	paymentModel := &model.PaymentModel{
		ID:         payment.ID,
		RefCode:    payment.RefCode,
		UserID:     payment.UserID,
		OrderID:    payment.OrderID,
		Amount:     payment.Amount,
		Status:     string(payment.Status),
		InvoiceUrl: payment.InvoiceURL,
		PaidAt:     payment.PaidAt,
	}
	if err := s.paymentProducer.PublicPaymentSuccess(ctx, paymentModel); err != nil {
		s.log.Error("failed to publish payment success", zap.String("order_id", paymentModel.OrderID), zap.Error(err))
		return err
	}

	s.log.Info("payment completed successfully", zap.String("order_id", orderID))
	return nil
}

func (s *paymentService) orderIDFromExternalID(externalID string) string {
	return strings.TrimPrefix(externalID, "	order-")
}

func (s *paymentService) DownloadInvoiceInPDF(ctx context.Context, orderID string) (string, error) {
	payment, err := s.paymentStore.FindByOrderID(ctx, orderID)
	if err != nil || payment == nil {
		s.log.Error("failed to find payment by order id", zap.String("order_id", orderID), zap.Error(err))
		return "", err
	}

	filePath := fmt.Sprintf("/fcproject/invoice_%v", orderID)
	paymentDetail := &utils.InvoicePDF{
		ID:         payment.ID,
		OrderID:    payment.OrderID,
		Amount:     payment.Amount,
		Status:     string(payment.Status),
		PaymentUrl: payment.InvoiceURL,
	}

	err = utils.InvoicePDFGenerator(paymentDetail, filePath)
	if err != nil {
		s.log.Error("failed to generate invoice", zap.String("order_id", orderID), zap.Error(err))
		return "", err
	}

	return filePath, nil
}
