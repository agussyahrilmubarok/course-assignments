package v4

import (
	"context"
	"encoding/json"
	"time"

	"example.com/coupon-service/internal/config"
	"example.com/coupon-service/internal/coupon"
	"example.com/coupon-service/internal/instrument/logging"
	"example.com/coupon-service/internal/instrument/tracing"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const TopicCouponIssue = "coupon-issue-requests"

type KafkaProducer struct {
	writer *kafka.Writer
}

func NewKafkaProducer(
	brokers []string,
) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Balancer:     &kafka.Hash{},
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
			RequiredAcks: kafka.RequireOne,
			MaxAttempts:  5,
			Async:        false,
		},
	}
}

func (p *KafkaProducer) SendIssueCoupon(ctx context.Context, message coupon.IssueCouponMessage) error {
	ctx, span := tracing.StartSpan(ctx, "V4.KafkaProducer.SendIssueCoupon")
	defer span.End()

	log := logging.GetLoggerFromContext(ctx)

	jsonValue, err := json.Marshal(message)
	if err != nil {
		span.RecordError(err)
		log.Error("failed to marshal issue coupon message", zap.String("policy_code", message.PolicyCode), zap.Error(err))
		return err
	}

	sendCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := p.writer.WriteMessages(sendCtx, kafka.Message{
		Topic: TopicCouponIssue,
		Key:   []byte(message.PolicyID),
		Value: jsonValue,
		Time:  time.Now(),
	}); err != nil {
		span.RecordError(err)
		log.Error("failed to send issue coupon to kafka", zap.String("policy_code", message.PolicyCode), zap.Error(err))
		return err
	}

	log.Info("success send issue coupon to kafka",
		zap.String("policy_id", message.PolicyID),
		zap.String("policy_code", message.PolicyCode),
		zap.String("coupon_id", message.CouponID),
		zap.String("user_id", message.UserID),
	)

	return nil
}

type KafkaConsumer struct {
	reader  *kafka.Reader
	service IService
}

func NewKafkaConsumer(
	cfg *config.Config,
	pg *config.Postgres,
	rdb *config.Redis,
) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  cfg.Kafka.Brokers,
		Topic:    TopicCouponIssue,
		GroupID:  cfg.Kafka.GroupID,
		MaxWait:  500 * time.Millisecond,
		MinBytes: 1,
		MaxBytes: 10e6,
	})

	service := NewService(NewRepository(pg, rdb), NewKafkaProducer(cfg.Kafka.Brokers))

	return &KafkaConsumer{
		reader:  reader,
		service: service,
	}
}

func (c *KafkaConsumer) Start(ctx context.Context) error {
	log := logging.GetLoggerFromContext(ctx)
	log.Info("kafka consumer started...", zap.String("topic", TopicCouponIssue))

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			if err == context.Canceled {
				log.Warn("kafka consumer stopping due cancellation")
				return nil
			}
			log.Error("failed read message from kafka", zap.Error(err))
			continue
		}

		go c.handleMessage(context.Background(), m)
	}
}

func (c *KafkaConsumer) handleMessage(ctx context.Context, msg kafka.Message) {
	ctx, span := tracing.StartSpan(ctx, "V4.KafkaConsumer.handleMessage")
	defer span.End()
	log := logging.GetLoggerFromContext(ctx)

	var data coupon.IssueCouponMessage
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		span.RecordError(err)
		log.Error("failed unmarshal coupon message", zap.Error(err))
		return
	}

	log.Info("received message from kafka",
		zap.String("policy_id", data.PolicyID),
		zap.String("policy_code", data.PolicyCode),
		zap.String("coupon_id", data.CouponID),
		zap.String("user_id", data.UserID),
	)

	// TODO: IDEMPOTENCY CHECK
	// example pseudo
	// key := "coupon:processed:" + data.RequestID
	// exist := redis.Exists(key)
	// if exist { return }

	if err := c.service.ProcessIssueCoupon(ctx, data); err != nil {
		span.RecordError(err)
		log.Error("failed to process issue coupon", zap.Error(err))
		return
	}

	log.Info("success processed coupon issuance",
		zap.String("policy_id", data.PolicyID),
		zap.String("policy_code", data.PolicyCode),
		zap.String("coupon_id", data.CouponID),
		zap.String("user_id", data.UserID),
	)
}

func (c *KafkaConsumer) Close() {
	_ = c.reader.Close()
}
