package consumer

import (
	"context"
	"encoding/json"

	"example.com/pkg/model"
	"example.com/productfc/internal/service"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaStockUpdateConsumer struct {
	reader         *kafka.Reader
	productService service.IProductService
	log            *zap.Logger
}

func NewKafkaStockUpdateConsumer(
	brokers []string,
	productService service.IProductService,
	log *zap.Logger,
) *KafkaStockUpdateConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   model.ORDER_CREATED_EVENT,
		GroupID: model.PRODUCT_GROUP,
	})

	return &KafkaStockUpdateConsumer{
		reader:         reader,
		productService: productService,
		log:            log,
	}
}

func (c *KafkaStockUpdateConsumer) Start(ctx context.Context) {
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

		for _, item := range orderModel.Items {
			if err := c.productService.DeductStockByID(ctx, item.ID, item.Quantity); err != nil {
				c.log.Error("failed to deduct stock when order created event", zap.Error(err))
				continue
			}
			c.log.Info("successfully deducted stock", zap.String("product_id", item.ID), zap.Int("quantity", item.Quantity))
		}
	}
}

func (c *KafkaStockUpdateConsumer) Stop(ctx context.Context) {
	if err := c.reader.Close(); err != nil {
		c.log.Error("failed to close kafka reader", zap.Error(err))
	} else {
		c.log.Info("kafka reader closed successfully")
	}
}
