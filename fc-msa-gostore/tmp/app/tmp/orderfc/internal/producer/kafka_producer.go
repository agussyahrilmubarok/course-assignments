package producer

import (
	"context"
	"encoding/json"
	"fmt"

	"example.com/pkg/model"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type KafkaProducer struct {
	writer *kafka.Writer
	log    *zap.Logger
}

func NewKafkaProducer(brokers []string, log *zap.Logger) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Balancer: &kafka.LeastBytes{},
	}

	return &KafkaProducer{
		writer: writer,
		log:    log,
	}
}

func (p *KafkaProducer) PublishOrderCreated(ctx context.Context, message *model.OrderModel) error {
	value, err := json.Marshal(message)
	if err != nil {
		p.log.Error("failed to marshal order created message", zap.Error(err))
		return err
	}

	msg := kafka.Message{
		Topic: model.ORDER_CREATED_EVENT,
		Key:   []byte(fmt.Sprintf("order-%s", message.ID)),
		Value: value,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.log.Error("failed to publish order created event", zap.Error(err))
		return err
	}

	return nil
}

func (p *KafkaProducer) PublishOrderCancelled(ctx context.Context, message *model.OrderModel) error {
	value, err := json.Marshal(message)
	if err != nil {
		p.log.Error("failed to marshal order cancelled message", zap.Error(err))
		return err
	}

	msg := kafka.Message{
		Topic: model.ORDER_CANCELLED_EVENT,
		Key:   []byte(fmt.Sprintf("order-%s", message.ID)),
		Value: value,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.log.Error("failed to publish order cancelled event", zap.Error(err))
		return err
	}

	return nil
}
