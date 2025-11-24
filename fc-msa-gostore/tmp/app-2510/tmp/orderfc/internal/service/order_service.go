package service

import (
	"context"
	"errors"

	"example.com/orderfc/internal/client"
	"example.com/orderfc/internal/producer"
	"example.com/orderfc/internal/store"
	"example.com/pkg/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type IOrderService interface {
	Checkout(ctx context.Context, req *model.CheckoutOrderRequest) (*model.OrderModel, error)
	Cancel(ctx context.Context, orderID string) (*model.OrderModel, error)
	GetHistory(ctx context.Context, userID string) ([]model.OrderModel, error)
}

type orderService struct {
	orderStore        store.IOrderStore
	productGrpcClient *client.ProductGrpcClient
	kafkaProducer     *producer.KafkaProducer
	log               *zap.Logger
}

func NewOrderService(
	orderStore store.IOrderStore,
	productGrpcClient *client.ProductGrpcClient,
	kafkaProducer *producer.KafkaProducer,
	log *zap.Logger,
) IOrderService {
	return &orderService{
		orderStore:        orderStore,
		productGrpcClient: productGrpcClient,
		kafkaProducer:     kafkaProducer,
		log:               log,
	}
}

func (s orderService) Checkout(ctx context.Context, req *model.CheckoutOrderRequest) (*model.OrderModel, error) {
	var (
		orderItems    []store.OrderItem
		totalAmount   float64
		totalQuantity int
	)

	for _, item := range req.Items {
		name, price, stock, err := s.productGrpcClient.GetProductInfo(ctx, item.ID)
		if err != nil {
			s.log.Error("failed to get product info", zap.String("product_id", item.ID), zap.Error(err))
			return nil, errors.New("failed to get product info")
		}

		if item.Quantity > stock {
			s.log.Error("not enough stock for product", zap.String("product_id", item.ID), zap.Error(err))
			return nil, errors.New("product out of stock")
		}

		orderItem := store.OrderItem{
			ProductID: item.ID,
			Name:      name,
			Price:     price,
			Quantity:  item.Quantity,
		}

		totalAmount += price * float64(item.Quantity)
		totalQuantity += item.Quantity
		orderItems = append(orderItems, orderItem)
	}

	order := &store.Order{
		ID:            uuid.New().String(),
		UserID:        req.UserID,
		Items:         orderItems,
		TotalAmount:   totalAmount,
		TotalQuantity: totalQuantity,
		Status:        store.OrderCreated,
	}

	order, err := s.orderStore.Create(ctx, order)
	if err != nil || order == nil {
		s.log.Error("failed to checkout create order", zap.Error(err))
		return nil, err
	}

	var itemResponses []model.OrderItemModel
	for _, item := range order.Items {
		itemResponses = append(itemResponses, model.OrderItemModel{
			ID:       item.ProductID,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}

	orderModel := &model.OrderModel{
		ID:            order.ID,
		UserID:        order.UserID,
		Items:         itemResponses,
		TotalAmount:   order.TotalAmount,
		TotalQuantity: order.TotalQuantity,
		Status:        string(order.Status),
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}

	if err = s.kafkaProducer.PublishOrderCreated(ctx, orderModel); err != nil {
		s.orderStore.DeleteByID(ctx, order.ID)
		s.log.Error("failed to publish order created event", zap.Error(err))
		return nil, err
	}

	return orderModel, nil
}

func (s orderService) GetHistory(ctx context.Context, userID string) ([]model.OrderModel, error) {
	orders, err := s.orderStore.FindAllByUserID(ctx, userID)
	if err != nil || orders == nil {
		s.log.Error("failed to get order history", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	var responses []model.OrderModel
	for _, order := range orders {
		var items []model.OrderItemModel
		for _, item := range order.Items {
			items = append(items, model.OrderItemModel{
				ID:       item.ProductID,
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			})
		}

		responses = append(responses, model.OrderModel{
			ID:            order.ID,
			UserID:        order.UserID,
			Items:         items,
			TotalAmount:   order.TotalAmount,
			TotalQuantity: order.TotalQuantity,
			Status:        string(order.Status),
			CreatedAt:     order.CreatedAt,
			UpdatedAt:     order.UpdatedAt,
		})
	}

	return responses, nil
}

func (s *orderService) Cancel(ctx context.Context, orderID string) (*model.OrderModel, error) {
	order, err := s.orderStore.FindByID(ctx, orderID)
	if err != nil || order == nil {
		s.log.Error("failed to find order by id", zap.String("order_id", orderID), zap.Error(err))
		return nil, err
	}

	switch order.Status {
	case store.OrderPaid, store.OrderProcessed:
		s.log.Error("cannot cancel order already paid or processed",
			zap.String("order_id", order.ID),
			zap.String("status", string(order.Status)),
		)
		return nil, errors.New("order cannot be cancelled")
	}

	order, err = s.orderStore.UpdateStatus(ctx, orderID, store.OrderCancelled)
	if err != nil || order == nil {
		s.log.Error("failed to update status cancelled", zap.String("order_id", orderID), zap.Error(err))
		return nil, err
	}

	var itemResponses []model.OrderItemModel
	for _, item := range order.Items {
		itemResponses = append(itemResponses, model.OrderItemModel{
			ID:       item.ProductID,
			Name:     item.Name,
			Price:    item.Price,
			Quantity: item.Quantity,
		})
	}

	orderModel := &model.OrderModel{
		ID:            order.ID,
		UserID:        order.UserID,
		Items:         itemResponses,
		TotalAmount:   order.TotalAmount,
		TotalQuantity: order.TotalQuantity,
		Status:        string(order.Status),
		CreatedAt:     order.CreatedAt,
		UpdatedAt:     order.UpdatedAt,
	}

	if err = s.kafkaProducer.PublishOrderCancelled(ctx, orderModel); err != nil {
		s.orderStore.UpdateStatus(ctx, orderID, store.OrderCreated)
		s.log.Error("failed to publish order cancelled", zap.String("order_id", order.ID), zap.Error(err))
		return nil, err
	}

	return orderModel, nil
}
