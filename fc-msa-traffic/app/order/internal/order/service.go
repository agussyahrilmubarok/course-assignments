package order

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=IService
type IService interface {
	CalculateAndCreateOrder(ctx context.Context, param CreateOrderRequest) (*Order, error)
	CancelOrderAndRestockProduct(ctx context.Context, param CancelOrderRequest) (*Order, error)
}

type service struct {
	cfg    *Config
	store  IStore
	client IClient
	log    zerolog.Logger
}

func NewService(
	cfg *Config,
	store IStore,
	client IClient,
	log zerolog.Logger,
) IService {
	return &service{
		cfg:    cfg,
		store:  store,
		client: client,
		log:    log,
	}
}

func (s *service) CalculateAndCreateOrder(ctx context.Context, param CreateOrderRequest) (*Order, error) {
	// Make channel for check stock and pricing product
	stockChan := make(chan stockResult, len(param.OrderItems))
	pricingChan := make(chan pricingResult, len(param.OrderItems))

	for _, val := range param.OrderItems {
		orderItem := val // capture by value for race conditions
		s.log.Info().Str("order_item_id", orderItem.ProductID).Msg("Capture order item")

		go s.checkStockAsync(ctx, orderItem, stockChan)
		go s.checkPricingAsync(ctx, orderItem, pricingChan)
	}

	// Temporary maps to match result with productID
	stockMap := make(map[string]bool)
	pricingMap := make(map[string]PricingResponse)

	// Collect results
	for i := 0; i < len(param.OrderItems); i++ {
		// Stock
		sr := <-stockChan
		if sr.Error != nil {
			s.log.Error().Err(sr.Error).Str("product_id", sr.ProductID).Msg("Check stock product fail")
			return nil, sr.Error
		}
		if !sr.Available {
			s.log.Error().Err(sr.Error).Str("product_id", sr.ProductID).Msg("Product out of stock")
			return nil, fmt.Errorf("product %s out of stock", sr.ProductID)
		}
		stockMap[sr.ProductID] = true
	}

	for i := 0; i < len(param.OrderItems); i++ {
		// Pricing
		pr := <-pricingChan
		if pr.Error != nil {
			s.log.Error().Err(pr.Error).Str("product_id", pr.ProductID).Msg("Check pricing product fail")
			return nil, pr.Error
		}
		pricingMap[pr.ProductID] = PricingResponse{
			ProductID:  pr.ProductID,
			Markup:     pr.MarkUp,
			Discount:   pr.Discount,
			FinalPrice: pr.FinalPrice,
		}
	}

	var (
		entityOrderItems    []OrderItem
		entityTotalPrice    float64
		entityTotalMarkUp   float64
		entityTotalDiscount float64
		entityTotalQuantity int
	)

	for _, orderItem := range param.OrderItems {
		price := pricingMap[orderItem.ProductID]
		qty := orderItem.Quantity

		// basePrice := price.FinalPrice / (1 + price.Markup - price.Discount)
		basePrice := price.FinalPrice
		totalPrice := basePrice * float64(qty)
		markUp := basePrice * price.Markup * float64(qty)
		discount := basePrice * price.Discount * float64(qty)

		entityOrderItems = append(entityOrderItems, OrderItem{
			ID:           uuid.NewString(),
			ProductID:    orderItem.ProductID,
			ProductPrice: price.FinalPrice,
			Quantity:     qty,
			MarkUp:       price.Markup,
			Discount:     price.Discount,
			TotalPrice:   price.FinalPrice * float64(qty),
		})

		entityTotalPrice += totalPrice
		entityTotalMarkUp += markUp
		entityTotalDiscount += discount
		entityTotalQuantity += qty
	}

	order := &Order{
		ID:         uuid.NewString(),
		UserID:     param.UserID,
		Status:     StatusProcessed,
		OrderItems: entityOrderItems,
	}

	if err := s.store.SaveOrder(ctx, order); err != nil {
		s.log.Error().Err(err).Msg("Failed create order when saving")
		return nil, err
	}

	for _, item := range order.OrderItems {
		if err := s.client.ReserveStockOnCatalogService(ctx, item.ProductID, item.Quantity); err != nil {
			s.log.Error().
				Err(err).
				Str("product_id", item.ProductID).
				Int("quantity", item.Quantity).
				Msg("Failed to reverse stock after order creation")

			for _, rollbackItem := range order.OrderItems {
				if rollbackItem.ProductID == item.ProductID {
					break
				}
				_ = s.client.ReleaseStockOnCatalogService(ctx, rollbackItem.ProductID, rollbackItem.Quantity)
			}

			_ = s.store.UpdateStatus(ctx, order, StatusFailed)
			return nil, fmt.Errorf("failed to reverse stock for product %s: %w", item.ProductID, err)
		}
	}

	s.log.Info().
		Str("order_id", order.ID).
		Int("total_items", len(order.OrderItems)).
		Msg("Successfully created order and reserved stock")

	return order, nil
}

func (s *service) CancelOrderAndRestockProduct(ctx context.Context, param CancelOrderRequest) (*Order, error) {
	order, err := s.store.FindOrderByID(ctx, param.OrderID)
	if err != nil || order == nil {
		s.log.Error().Err(err).Str("order_id", param.OrderID).Msg("Order not found")
		return nil, err
	}

	if order.UserID != param.UserID {
		s.log.Error().Str("order_id", param.OrderID).Str("user_id", param.UserID).Msg("Order is not permission to cancel")
		return nil, errors.New("do not have permissions to cancel")
	}

	if err := s.store.UpdateStatus(ctx, order, StatusCancelled); err != nil {
		s.log.Error().Str("order_id", param.OrderID).Str("user_id", param.UserID).Msg("Failed to set cancelled on order")
		return nil, err
	}

	for _, item := range order.OrderItems {
		if err := s.client.ReleaseStockOnCatalogService(ctx, item.ProductID, item.Quantity); err != nil {
			s.log.Error().
				Err(err).
				Str("product_id", item.ProductID).
				Int("quantity", item.Quantity).
				Msg("Failed to release stock when cancelling order")
			return nil, fmt.Errorf("failed to release stock for product %s: %w", item.ProductID, err)
		}
	}

	s.log.Info().
		Str("order_id", order.ID).
		Int("total_items", len(order.OrderItems)).
		Msg("Successfully cancelled order and released stock")

	return order, nil
}

type stockResult struct {
	ProductID string
	Available bool
	Error     error
}

func (s *service) checkStockAsync(ctx context.Context, orderItem OrderItemRequest, ch chan<- stockResult) {
	defer func() {
		if r := recover(); r != nil {
			s.log.Error().Any("recover", r).Msg("Recovered from panic in stock goroutine")
			ch <- stockResult{ProductID: orderItem.ProductID, Error: fmt.Errorf("error: %v", r)}
		}
	}()

	stock, err := s.client.GetStockOnCatalogService(ctx, orderItem.ProductID)
	if err != nil {
		s.log.Error().Err(err).Str("product_id", orderItem.ProductID).Msg("Failed to check stock in catalog service")
		ch <- stockResult{
			ProductID: orderItem.ProductID,
			Available: false,
			Error:     err,
		}

		return
	}

	ch <- stockResult{
		ProductID: orderItem.ProductID,
		Available: stock >= orderItem.Quantity,
		Error:     nil,
	}
}

type pricingResult struct {
	ProductID  string
	Discount   float64
	MarkUp     float64
	FinalPrice float64
	Error      error
}

func (s *service) checkPricingAsync(ctx context.Context, orderItem OrderItemRequest, ch chan<- pricingResult) {
	defer func() {
		if r := recover(); r != nil {
			s.log.Error().Any("recover", r).Msg("Recovered from panic in pricing goroutine")
			ch <- pricingResult{ProductID: orderItem.ProductID, Error: fmt.Errorf("error: %v", r)}
		}
	}()

	pricing, err := s.client.GetPricingOnPricingService(ctx, orderItem.ProductID)
	if err != nil {
		s.log.Error().Err(err).
			Str("product_id", orderItem.ProductID).
			Msg("Failed to get pricing from pricing service")

		ch <- pricingResult{
			ProductID: orderItem.ProductID,
			Error:     err,
		}
		return
	}

	ch <- pricingResult{
		ProductID:  orderItem.ProductID,
		Discount:   pricing.Discount,
		MarkUp:     pricing.Markup,
		FinalPrice: pricing.FinalPrice,
		Error:      nil,
	}
}
