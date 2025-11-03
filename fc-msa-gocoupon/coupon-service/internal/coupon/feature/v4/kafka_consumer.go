package v4

import (
	"context"
	"encoding/json"
	"time"

	"example.com/coupon/internal/coupon"
	"example.com/coupon/pkg/instrument"
	"github.com/rs/zerolog"
	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type KafkaConsumer struct {
	reader        *kafka.Reader
	couponFeature ICouponFeature
	log           zerolog.Logger
	tracer        trace.Tracer
}

func NewKafkaConsumer(
	brokers []string,
	groupID string,
	couponFeature ICouponFeature,
	log zerolog.Logger,
	tracer trace.Tracer,
) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		GroupID:        groupID,
		Topic:          TopicCouponIssueRequest,
		MinBytes:       10e3,        // 10KB
		MaxBytes:       10e6,        // 10MB
		CommitInterval: time.Second, // auto-commit interval
	})

	return &KafkaConsumer{
		reader:        reader,
		couponFeature: couponFeature,
		log:           log,
		tracer:        tracer,
	}
}

func (kc *KafkaConsumer) ConsumeCouponIssueRequest(ctx context.Context) error {
	ctx, span := kc.tracer.Start(ctx, "feature.KafkaConsumer.ConsumeCouponIssueRequest")
	defer span.End()

	log := instrument.GetLogger(ctx, kc.log)
	log = log.With().Str("func", "ConsumeCouponIssueRequest").Logger()

	for {
		m, err := kc.reader.ReadMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				// context canceled → graceful shutdown
				log.Info().Msg("Kafka consumer context canceled, stopping consumer loop")
				return nil
			}

			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to read Kafka message")
			log.Error().Err(err).Msg("Failed to read message from Kafka")
			continue
		}

		var msg coupon.IssueCouponMessage
		if err := json.Unmarshal(m.Value, &msg); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, "Failed to unmarshal Kafka message")
			log.Error().
				Err(err).
				Str("topic", m.Topic).
				Msg("Failed to unmarshal Kafka message value")
			continue
		}

		childCtx, childSpan := kc.tracer.Start(ctx, "feature.KafkaConsumer.ProcessIssueCoupon",
			trace.WithAttributes(
				attribute.String("coupon.policy_code", msg.CouponPolicyCode),
				attribute.String("user.id", msg.UserID),
			),
		)

		log.Info().
			Str("topic", m.Topic).
			Str("coupon_policy_code", msg.CouponPolicyCode).
			Str("user_id", msg.UserID).
			Msg("Received coupon issue message")

		if err := kc.couponFeature.ProcessIssueCoupon(childCtx, msg); err != nil {
			childSpan.RecordError(err)
			childSpan.SetStatus(codes.Error, "Failed to process coupon issue")
			log.Error().
				Err(err).
				Str("coupon_policy_code", msg.CouponPolicyCode).
				Str("user_id", msg.UserID).
				Msg("Failed to process coupon issue request")

			// TODO: optionally publish to dead-letter topic
		} else {
			childSpan.SetStatus(codes.Ok, "Coupon issue processed successfully")
			log.Info().
				Str("coupon_policy_code", msg.CouponPolicyCode).
				Str("user_id", msg.UserID).
				Msg("Coupon issue processed successfully")
		}

		childSpan.End()
	}
}

func (kc *KafkaConsumer) Stop(ctx context.Context) error {
	log := instrument.GetLogger(ctx, kc.log)

	if kc.reader == nil {
		return nil
	}

	log.Info().Msg("Shutting down Kafka consumer...")
	if err := kc.reader.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close Kafka reader")
		return err
	}

	log.Info().Msg("Kafka consumer stopped successfully")
	return nil
}
