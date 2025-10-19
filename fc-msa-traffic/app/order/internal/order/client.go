package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=IClient
type IClient interface {
	GetStockOnCatalogService(ctx context.Context, productID string) (int, error)
	GetPricingOnPricingService(ctx context.Context, productID string) (*PricingResponse, error)
	ReverseStockOnCatalogService(ctx context.Context, productID string, quantity int) error
	ReleaseStockOnCatalogService(ctx context.Context, productID string, quantity int) error
}

type client struct {
	cfg *Config
	log zerolog.Logger
}

func NewClient(
	cfg *Config,
	log zerolog.Logger,
) IClient {
	return &client{
		cfg: cfg,
		log: log,
	}
}

func (c *client) GetStockOnCatalogService(ctx context.Context, productID string) (int, error) {
	url := fmt.Sprintf("%s/api/v1/catalogs/products/stocks/%s", "http://catalog-service:8082", productID)
	if c.cfg.App.Env != "production" {
		url = fmt.Sprintf("%s/api/v1/catalogs/products/stocks/%s", "http://localhost:8082", productID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.log.Error().Err(err).Msg("Failed to create request get stock product")
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.log.Error().Err(err).Msg("Failed to get stock product")
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.log.Error().Err(err).Msg("Failed to get stock - non 200 response")
		return 0, fmt.Errorf("non-200 status code: %d", res.StatusCode)
	}

	var stockData map[string]int
	if err := json.NewDecoder(res.Body).Decode(&stockData); err != nil {
		c.log.Error().Err(err).Msg("Failed to decode stock product")
		return 0, err
	}

	c.log.Info().
		Str("product_id", productID).
		Int("product_stock", stockData["stock"]).
		Msg("Get pricing product successfully")

	availableStock := stockData["stock"]
	return availableStock, nil
}

func (c *client) GetPricingOnPricingService(ctx context.Context, productID string) (*PricingResponse, error) {
	url := fmt.Sprintf("%s/api/v1/pricings/%s", "http://pricing-service:8084", productID)
	if c.cfg.App.Env != "production" {
		url = fmt.Sprintf("%s/api/v1/pricings/%s", "http://localhost:8084", productID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.log.Error().Err(err).Msg("Failed to create request get pricing product")
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.log.Error().Err(err).Msg("Failed to get pricing product")
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		c.log.Error().Err(err).Msg("Failed to get pricing - non 200 response")
		return nil, fmt.Errorf("non-200 status code: %d", res.StatusCode)
	}

	var pricing PricingResponse
	if err := json.NewDecoder(res.Body).Decode(&pricing); err != nil {
		c.log.Error().Err(err).Msg("Failed to decode stock product")
		return nil, err
	}

	c.log.Info().
		Str("product_id", productID).
		Float64("product_pricing", pricing.FinalPrice).
		Msg("Get pricing product successfully")

	return &pricing, nil
}

func (c *client) ReverseStockOnCatalogService(ctx context.Context, productID string, quantity int) error {
	url := fmt.Sprintf("%s/api/v1/catalogs/products/reverse", "http://catalog-service:8082")
	if c.cfg.App.Env != "production" {
		url = fmt.Sprintf("%s/api/v1/catalogs/products/reverse", "http://localhost:8082")
	}

	body := map[string]interface{}{
		"product_id": productID,
		"quantity":   quantity,
	}

	bodyBytes, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create reverse stock request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call reverse stock API: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("reverse stock failed: %d", res.StatusCode)
	}

	c.log.Info().
		Str("product_id", productID).
		Int("quantity", quantity).
		Msg("Reversed product stock successfully")

	return nil
}

func (c *client) ReleaseStockOnCatalogService(ctx context.Context, productID string, quantity int) error {
	url := fmt.Sprintf("%s/api/v1/catalogs/products/release", "http://catalog-service:8082")
	if c.cfg.App.Env != "production" {
		url = fmt.Sprintf("%s/api/v1/catalogs/products/release", "http://localhost:8082")
	}

	body := map[string]interface{}{
		"product_id": productID,
		"quantity":   quantity,
	}

	bodyBytes, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create release stock request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call release stock API: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("release stock failed: %d", res.StatusCode)
	}

	c.log.Info().
		Str("product_id", productID).
		Int("quantity", quantity).
		Msg("Released product stock successfully")

	return nil
}
