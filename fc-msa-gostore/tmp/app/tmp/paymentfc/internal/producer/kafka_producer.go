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

func (p *KafkaProducer) PublicPaymentSuccess(ctx context.Context, message *model.PaymentModel) error {
	value, err := json.Marshal(message)
	if err != nil {
		p.log.Error("failed to marshal payment success message", zap.Error(err))
		return err
	}

	msg := kafka.Message{
		Topic: model.PAYMENT_SUCCESS_EVENT,
		Key:   []byte(fmt.Sprintf("payment-%s", message.ID)),
		Value: value,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.log.Error("failed to publish payment success event", zap.Error(err))
		return err
	}

	return nil
}

func (p *KafkaProducer) PublicPaymentFailed(ctx context.Context, message *model.PaymentModel) error {
	value, err := json.Marshal(message)
	if err != nil {
		p.log.Error("failed to marshal payment failed message", zap.Error(err))
		return err
	}

	msg := kafka.Message{
		Topic: model.PAYMENT_FAILED_EVENT,
		Key:   []byte(fmt.Sprintf("payment-%s", message.ID)),
		Value: value,
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		p.log.Error("failed to publish payment failed event", zap.Error(err))
		return err
	}

	return nil
}
