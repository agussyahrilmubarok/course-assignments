package service

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"strings"

	"example.com/backend/internal/model"
	"example.com/backend/pkg/connections"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=IMidtransService
type IMidtransService interface {
	CreateTransaction(request model.MidtransRequest) (string, error)
	HandlerNotification(callback model.MidtransCallback) error
}

type midtransService struct {
	cfg connections.Midtrans
	log zerolog.Logger
}

func NewMidtransService(cfg connections.Midtrans, log zerolog.Logger) IMidtransService {
	return &midtransService{
		cfg: cfg,
		log: log,
	}
}

// CreateTransaction
// INFO: This method initializes a Snap transaction using Midtrans API,
//
//	generating a redirect URL for the user to complete the payment.
//	The created transaction will trigger a callback (webhook) from Midtrans
//	once the payment status changes (e.g. settlement, expire, cancel).
//
// NEXT: Setup Midtrans Webhook
// 1. Go to https://dashboard.midtrans.com/
// 2. Navigate to Settings > Configuration
// 3. Set the "Payment Notification URL" to your backend endpoint, e.g., https://yourdomain.com/api/v1/donations/callback
// 4. Ensure the endpoint is publicly accessible (use ngrok or Cloudflare Tunnel during local development)
// 5. Optionally verify the `signature_key` from the webhook payload for added security
func (s *midtransService) CreateTransaction(request model.MidtransRequest) (string, error) {
	// Create Snap Client
	snapClient := connections.NewMidtransSnapClient(s.cfg)

	// Initialize Snap Request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  request.TransactionDetails.OrderID, // Change with code or id transaction
			GrossAmt: request.TransactionDetails.GrossAmt,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: request.CustomerDetail.FName,
			Email: request.CustomerDetail.Email,
		},
	}

	// Create Midtrans Transaction
	snapRes, err := snapClient.CreateTransaction(snapReq)
	if err != nil {
		s.log.Error().Err(err).Msg("create midtrans transaction fail")
		return "", errors.New("Transaction fail")
	}

	return snapRes.RedirectURL, nil
}

// HandleMidtransWebhook
// INFO: This method handles incoming HTTP POST notifications (webhooks) from Midtrans
//
//	when a transaction status changes (e.g., settlement, pending, expire, cancel, deny).
//	It should parse the JSON payload, verify its authenticity (optional),
//	and update the corresponding donation record in the database accordingly.
//
// NEXT:
// 1. Parse the JSON body sent by Midtrans into a struct (e.g., midtrans.TransactionStatusResponse)
// 2. (Optional) Verify `signature_key` using: sha512(order_id + status_code + gross_amount + server_key)
// 3. Respond with 200 OK to acknowledge receipt, otherwise Midtrans will retry the callback
func (s *midtransService) HandlerNotification(callback model.MidtransCallback) error {
	// Verify signature
	serverKey := s.cfg.ServerKey
	orderID := callback.OrderID
	statusCode := callback.StatusCode
	grossAmount := callback.GrossAmount
	signatureKey := callback.SignatureKey

	raw := orderID + statusCode + grossAmount + serverKey
	hash := sha512.Sum512([]byte(raw))
	expectedSignature := fmt.Sprintf("%x", hash[:])
	if !strings.EqualFold(signatureKey, expectedSignature) {
		s.log.Error().Msgf("Signature mismatch for order_id=%s", orderID)
		return errors.New("Invalid signature")
	}

	return nil
}
