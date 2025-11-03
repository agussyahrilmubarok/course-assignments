package v4

import (
	"context"
	"encoding/json"

	"example.com/coupon/internal/coupon"
	"example.com/coupon/pkg/instrument"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	TopicCouponIssueRequest = "coupon-issue-request"
)

type KafkaProducer struct {
	brokers []string
	log     zerolog.Logger
	tracer  trace.Tracer
}

func NewKafkaProducer(brokers []string, log zerolog.Logger, tracer trace.Tracer) *KafkaProducer {
	return &KafkaProducer{
		brokers: brokers,
		log:     log,
		tracer:  tracer,
	}
}

func (kp *KafkaProducer) SendCouponIssueRequest(ctx context.Context, message coupon.IssueCouponMessage) error {
	ctx, span := kp.tracer.Start(ctx, "feature.KafkaProducer.SendCouponIssueRequest",
		trace.WithAttributes(
			attribute.String("coupon.policy_code", message.CouponPolicyCode),
			attribute.String("user.id", message.UserID),
		),
	)
	defer span.End()

	log := instrument.GetLogger(ctx, kp.log)

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  kp.brokers,
		Topic:    TopicCouponIssueRequest,
		Balancer: &kafka.Hash{},
		Async:    true, // non-blocking send
	})
	defer func() {
		if err := writer.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close Kafka writer")
		}
	}()

	value, err := json.Marshal(message)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to marshal message")
		log.Error().
			Err(err).
			Str("topic", TopicCouponIssueRequest).
			Msg("Failed to marshal coupon issue message")
		return err
	}

	err = writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(message.CouponPolicyCode),
		Value: value,
	})
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "Failed to send Kafka message")
		log.Error().
			Err(err).
			Str("topic", TopicCouponIssueRequest).
			Str("coupon_policy_code", message.CouponPolicyCode).
			Str("user_id", message.UserID).
			Msg("Failed to send coupon issue message to Kafka")
		return err
	}

	// Success log and trace
	span.SetStatus(codes.Ok, "Message sent successfully")
	log.Info().
		Str("topic", TopicCouponIssueRequest).
		Str("coupon_policy_code", message.CouponPolicyCode).
		Str("user_id", message.UserID).
		Msg("Coupon issue message sent successfully to Kafka")
	return nil
}
