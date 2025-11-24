package scheduler

import (
	"context"
	"time"

	"example.com/paymentfc/internal/store"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type PaymentScheduler struct {
	cron         *cron.Cron
	paymentStore store.IPaymentStore
	log          *zap.Logger
}

func NewPaymentScheduler(
	paymentStore store.IPaymentStore,
	log *zap.Logger,
) *PaymentScheduler {
	cron := cron.New()
	return &PaymentScheduler{
		cron:         cron,
		paymentStore: paymentStore,
		log:          log,
	}
}

func (s *PaymentScheduler) Start() {
	// Schedule the payment pending sync every 10 minutes
	_, err := s.cron.AddFunc("@every 10m", func() {
		s.log.Info("Starting payment pending sync task")

		ctx := context.Background()
		pendingPayments, err := s.paymentStore.FindAllByStatus(ctx, store.StatusPending)
		if err != nil {
			s.log.Warn("failed to fetch pending payments", zap.Error(err))
			return
		}

		now := time.Now()
		for _, payment := range pendingPayments {
			if payment.ExpiredAt.Before(now) {
				s.log.Info("marking payment as failed due to expiration", zap.String("payment_id", payment.ID))
				_, err := s.paymentStore.UpdateStatus(ctx, payment.ID, store.StatusFailed)
				if err != nil {
					s.log.Error("Failed to update payment status", zap.String("payment_id", payment.ID), zap.Error(err))
					continue
				}

				// TODO: publish rollback stock event or cancel order
			}
		}

	})

	if err != nil {
		s.log.Error("failed to schedule payment sync task", zap.Error(err))
		return
	}

	s.cron.Start()
	s.log.Info("payment scheduler started")
}

func (s *PaymentScheduler) Stop() {
	s.cron.Stop()
	s.log.Info("Payment scheduler stopped")
}
