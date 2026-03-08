package pricing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=IClient
type IClient interface {
	GetPriceProduct(ctx context.Context, productID string) (float64, error)
	GetStockProduct(ctx context.Context, productID string) (int, error)
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

func (c *client) GetPriceProduct(ctx context.Context, productID string) (float64, error) {
	url := fmt.Sprintf("%s/api/v1/catalogs/products/%s", "http://catalog-service:8082", productID)
	if c.cfg.App.Env != "production" {
		url = fmt.Sprintf("%s/api/v1/catalogs/products/%s", "http://localhost:8082", productID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		c.log.Error().Err(err).Msg("Failed to create request get price product")
		return 0, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.log.Error().Err(err).Msg("Failed to get price product")
		return 0, err
	}
	defer res.Body.Close()

	product := &struct {
		Price float64 `json:"price"`
	}{}

	if err := json.NewDecoder(res.Body).Decode(product); err != nil {
		c.log.Error().Err(err).Msg("Failed to decode get price product")
		return 0, err
	}

	return product.Price, nil
}

func (c *client) GetStockProduct(ctx context.Context, productID string) (int, error) {
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

	availableStock := stockData["stock"]
	return availableStock, nil
}
