package consumer

import (
	"context"
	"encoding/json"

	"example.com/paymentfc/internal/service"
	"example.com/pkg/model"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaInitPaymentConsumer struct {
	reader         *kafka.Reader
	paymentService service.IPaymentService
	log            *zap.Logger
}

func NewKafkaInitPaymentConsumer(
	brokers []string,
	paymentService service.IPaymentService,
	log *zap.Logger,
) *KafkaInitPaymentConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   model.ORDER_CREATED_EVENT,
		GroupID: model.PAYMENT_GROUP,
	})

	return &KafkaInitPaymentConsumer{
		reader:         reader,
		paymentService: paymentService,
		log:            log,
	}
}

func (c *KafkaInitPaymentConsumer) Start(ctx context.Context) {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.log.Error("failed to consume message order created event", zap.Error(err))
			continue
		}

		var orderModel model.OrderModel
		err = json.Unmarshal(message.Value, &orderModel)
		if err != nil {
			c.log.Error("failed to consume message order created event", zap.String("order_id", string(message.Key)), zap.Error(err))
			continue
		}

		// Toggle for auto payment without xendit
		togglePaymentDirecly := false
		if togglePaymentDirecly {
			if err := c.paymentService.ProcessPaymentDirectly(ctx, &orderModel); err != nil {
				c.log.Error("failed to initilize payment directly when order created event", zap.Error(err))
				continue
			}
		} else {
			if err := c.paymentService.CreateInvoiceOrder(ctx, &orderModel); err != nil {
				c.log.Error("failed to initilize payment when order created event", zap.Error(err))
				continue
			}
		}
		c.log.Info("successfully initialize", zap.String("user_id", orderModel.UserID), zap.String("order_id", orderModel.ID))
	}
}

func (c *KafkaInitPaymentConsumer) Stop(ctx context.Context) {
	if err := c.reader.Close(); err != nil {
		c.log.Error("failed to close kafka reader", zap.Error(err))
	} else {
		c.log.Info("kafka reader closed successfully")
	}
}
