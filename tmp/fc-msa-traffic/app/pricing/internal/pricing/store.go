package pricing

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=IStore
type IStore interface {
	SetPricing(ctx context.Context, pricing *Pricing) error
	GetPricing(ctx context.Context, productID string) (*Pricing, error)
	SetPricingRule(ctx context.Context, rule *PricingRule) error
	GetPricingRule(ctx context.Context, productID string) (*PricingRule, error)
}

const (
	PRICING_KEY       = "pricing:%s"      // Redis key format for Pricing
	PRICING_TTL       = time.Hour * 24    // 1 hour TTL
	PRICING_RULES_KEY = "pricing:rule:%s" // Redis key format for PricingRule
	PRICING_RULES_TTL = time.Hour * 24    // 1 day TTL
)

type store struct {
	rdb *redis.Client
	log zerolog.Logger
}

func NewStore(rdb *redis.Client, log zerolog.Logger) IStore {
	return &store{
		rdb: rdb,
		log: log,
	}
}

func (s *store) SetPricing(ctx context.Context, pricing *Pricing) error {
	key := fmt.Sprintf(PRICING_KEY, pricing.ProductID)
	data, err := json.Marshal(pricing)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to marshal pricing data")
		return err
	}

	if err := s.rdb.Set(ctx, key, data, PRICING_TTL).Err(); err != nil {
		s.log.Error().Err(err).Str("key", key).Msg("Failed to set pricing in redis")
		return err
	}

	s.log.Info().Str("key", key).Msg("Pricing data cached successfully")
	return nil
}

func (s *store) GetPricing(ctx context.Context, productID string) (*Pricing, error) {
	key := fmt.Sprintf(PRICING_KEY, productID)
	data, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			s.log.Debug().Str("key", key).Msg("Pricing not found in redis")
			return nil, nil
		}
		s.log.Error().Err(err).Str("key", key).Msg("Failed to get pricing from redis")
		return nil, err
	}

	var pricing Pricing
	if err := json.Unmarshal([]byte(data), &pricing); err != nil {
		s.log.Error().Err(err).Str("key", key).Msg("Failed to unmarshal pricing data")
		return nil, err
	}

	s.log.Info().Str("key", key).Msg("Get Pricing data from cached successfully")
	return &pricing, nil
}

func (s *store) SetPricingRule(ctx context.Context, rule *PricingRule) error {
	key := fmt.Sprintf(PRICING_RULES_KEY, rule.ProductID)
	data, err := json.Marshal(rule)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to marshal pricing rule data")
		return err
	}

	if err := s.rdb.Set(ctx, key, data, PRICING_RULES_TTL).Err(); err != nil {
		s.log.Error().Err(err).Str("key", key).Msg("Failed to set pricing rule in redis")
		return err
	}

	s.log.Info().Str("key", key).Msg("Pricing rule cached successfully")
	return nil
}

func (s *store) GetPricingRule(ctx context.Context, productID string) (*PricingRule, error) {
	key := fmt.Sprintf(PRICING_RULES_KEY, productID)
	data, err := s.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			s.log.Debug().Str("key", key).Msg("Pricing rule not found in redis")
			return nil, nil
		}
		s.log.Error().Err(err).Str("key", key).Msg("Failed to get pricing rule from redis")
		return nil, err
	}

	var rule PricingRule
	if err := json.Unmarshal([]byte(data), &rule); err != nil {
		s.log.Error().Err(err).Str("key", key).Msg("Failed to unmarshal pricing rule data")
		return nil, err
	}

	s.log.Info().Str("key", key).Msg("Get Pricing rule from cached successfully")
	return &rule, nil
}
