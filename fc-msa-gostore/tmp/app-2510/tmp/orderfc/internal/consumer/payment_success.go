package consumer

import (
	"context"
	"encoding/json"

	"example.com/orderfc/internal/store"
	"example.com/pkg/model"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaPaymentSuccessConsumer struct {
	reader     *kafka.Reader
	orderStore store.IOrderStore
	log        *zap.Logger
}

func NewKafkaPaymentSuccessConsumer(
	brokers []string,
	orderStore store.IOrderStore,
	log *zap.Logger,
) *KafkaPaymentSuccessConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   model.PAYMENT_SUCCESS_EVENT,
		GroupID: model.PRODUCT_GROUP,
	})

	return &KafkaPaymentSuccessConsumer{
		reader:     reader,
		orderStore: orderStore,
		log:        log,
	}
}

func (c *KafkaPaymentSuccessConsumer) Start(ctx context.Context) {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.log.Error("failed to consume message payment success event", zap.Error(err))
			continue
		}

		var paymentModel model.PaymentModel
		err = json.Unmarshal(message.Value, &paymentModel)
		if err != nil {
			c.log.Error("failed to consume message payment success event", zap.String("order_id", paymentModel.OrderID), zap.Error(err))
			continue
		}

		_, err = c.orderStore.UpdateStatus(ctx, paymentModel.OrderID, store.OrderPaid)
		if err != nil {
			c.log.Error("failed to set order paid", zap.String("order_id", paymentModel.OrderID), zap.Error(err))
			continue
		}

		c.log.Info("successfully set order paid", zap.String("order_id", paymentModel.OrderID))
	}
}

func (c *KafkaPaymentSuccessConsumer) Stop(ctx context.Context) {
	if err := c.reader.Close(); err != nil {
		c.log.Error("failed to close kafka reader", zap.Error(err))
	} else {
		c.log.Info("kafka reader closed successfully")
	}
}
