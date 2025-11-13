package consumer

import (
	"context"
	"encoding/json"

	"example.com/pkg/model"
	"example.com/productfc/internal/service"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaStockRollbackConsumer struct {
	reader         *kafka.Reader
	productService service.IProductService
	log            *zap.Logger
}

func NewKafkaStockRollbackConsumer(
	brokers []string,
	productService service.IProductService,
	log *zap.Logger,
) *KafkaStockRollbackConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   model.ORDER_CANCELLED_EVENT,
		GroupID: model.PRODUCT_GROUP,
	})

	return &KafkaStockRollbackConsumer{
		reader:         reader,
		productService: productService,
		log:            log,
	}
}

func (c *KafkaStockRollbackConsumer) Start(ctx context.Context) {
	for {
		message, err := c.reader.ReadMessage(ctx)
		if err != nil {
			c.log.Error("failed to consume message order cancelled event", zap.Error(err))
			continue
		}

		var orderModel model.OrderModel
		err = json.Unmarshal(message.Value, &orderModel)
		if err != nil {
			c.log.Error("failed to consume message order cancelled event", zap.String("order_id", string(message.Key)), zap.Error(err))
			continue
		}

		for _, item := range orderModel.Items {
			if err := c.productService.RollbackStockByID(ctx, item.ID, item.Quantity); err != nil {
				c.log.Error("failed to rollback stock when order cancelled event", zap.Error(err))
				continue
			}
			c.log.Info("successfully rollback stock", zap.String("product_id", item.ID), zap.Int("quantity", item.Quantity))
		}
	}
}

func (c *KafkaStockRollbackConsumer) Stop(ctx context.Context) {
	if err := c.reader.Close(); err != nil {
		c.log.Error("failed to close kafka reader", zap.Error(err))
	} else {
		c.log.Info("kafka reader closed successfully")
	}
}
