package consumer

import (
	"context"
	"encoding/json"

	"example.com/orderfc/internal/producer"
	"example.com/orderfc/internal/store"
	"example.com/pkg/model"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaPaymentFailedConsumer struct {
	reader        *kafka.Reader
	orderStore    store.IOrderStore
	kafkaProducer *producer.KafkaProducer
	log           *zap.Logger
}

func NewKafkaPaymentFailedConsumer(
	brokers []string,
	orderStore store.IOrderStore,
	kafkaProducer *producer.KafkaProducer,
	log *zap.Logger,
) *KafkaPaymentFailedConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   model.PAYMENT_FAILED_EVENT,
		GroupID: model.PRODUCT_GROUP,
	})

	return &KafkaPaymentFailedConsumer{
		reader:        reader,
		orderStore:    orderStore,
		kafkaProducer: kafkaProducer,
		log:           log,
	}
}

func (c *KafkaPaymentFailedConsumer) Start(ctx context.Context) {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.log.Error("failed to consume message payment failed event", zap.Error(err))
			continue
		}

		var paymentModel model.PaymentModel
		err = json.Unmarshal(message.Value, &paymentModel)
		if err != nil {
			c.log.Error("failed to consume message payment failed event", zap.String("order_id", paymentModel.OrderID), zap.Error(err))
			continue
		}

		order, err := c.orderStore.FindByID(ctx, paymentModel.OrderID)
		if err != nil || order == nil {
			c.log.Error("failed to find order when set to cancelled", zap.String("order_id", paymentModel.OrderID), zap.Error(err))
			continue
		}

		_, err = c.orderStore.UpdateStatus(ctx, paymentModel.OrderID, store.OrderCancelled)
		if err != nil {
			c.log.Error("failed to set order cancelled", zap.String("order_id", paymentModel.OrderID), zap.Error(err))
			continue
		}

		var itemResponses []model.OrderItemModel
		for _, item := range order.Items {
			itemResponses = append(itemResponses, model.OrderItemModel{
				ID:       item.ProductID,
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			})
		}

		orderModel := &model.OrderModel{
			ID:            order.ID,
			UserID:        order.UserID,
			Items:         itemResponses,
			TotalAmount:   order.TotalAmount,
			TotalQuantity: order.TotalQuantity,
			Status:        string(order.Status),
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
		}

		if err = c.kafkaProducer.PublishOrderCancelled(ctx, orderModel); err != nil {
			c.log.Error("failed to publish order cancelled when payment failed", zap.String("order_id", order.ID), zap.Error(err))
			continue
		}

		c.log.Info("successfully set order cancelled", zap.String("order_id", paymentModel.OrderID))
	}
}

func (c *KafkaPaymentFailedConsumer) Stop(ctx context.Context) {
	if err := c.reader.Close(); err != nil {
		c.log.Error("failed to close kafka reader", zap.Error(err))
	} else {
		c.log.Info("kafka reader closed successfully")
	}
}
