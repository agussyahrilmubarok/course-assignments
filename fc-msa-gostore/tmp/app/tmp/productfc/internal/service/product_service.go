package service

import (
	"context"

	"example.com/pkg/model"
	"example.com/productfc/internal/store"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

//go:generate mockery --name=IProductService
type IProductService interface {
	Create(ctx context.Context, req *model.CreateProductRequest) (*model.ProductModel, error)
	UpdateByID(ctx context.Context, productID string, req *model.UpdateProductRequest) (*model.ProductModel, error)
	DeleteByID(ctx context.Context, productID string) error
	FindByID(ctx context.Context, productID string) (*model.ProductModel, error)
	DeductStockByID(ctx context.Context, productID string, quantity int) error
	RollbackStockByID(ctx context.Context, productID string, quantity int) error

	//SearchByID(ctx context.Context) ([]model.ProductModel, error)
}

type productService struct {
	productStore store.IProductStore
	log          *zap.Logger
}

func NewProductService(productStore store.IProductStore, log *zap.Logger) IProductService {
	return &productService{
		productStore: productStore,
		log:          log,
	}
}

func (s *productService) Create(ctx context.Context, req *model.CreateProductRequest) (*model.ProductModel, error) {
	product := &store.Product{
		ID:    uuid.New().String(),
		Name:  req.Name,
		Price: req.Price,
		Stock: req.Stock,
	}

	product, err := s.productStore.Create(ctx, product)
	if err != nil || product == nil {
		s.log.Error("failed to save new product", zap.Error(err))
		return nil, err
	}

	return &model.ProductModel{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil
}

func (s *productService) UpdateByID(ctx context.Context, productID string, req *model.UpdateProductRequest) (*model.ProductModel, error) {
	product, err := s.productStore.FindByID(ctx, productID)
	if err != nil || product == nil {
		s.log.Error("failed to find product by id", zap.String("product_id", productID), zap.Error(err))
		return nil, err
	}

	if req.Name != nil {
		product.Name = *req.Name
	}

	if req.Price != nil {
		product.Price = *req.Price
	}

	if req.Stock != nil {
		product.Stock = *req.Stock
	}

	product, err = s.productStore.UpdateByID(ctx, productID, product)
	if err != nil || product == nil {
		s.log.Error("failed to update product by id", zap.String("product_id", productID), zap.Error(err))
		return nil, err
	}

	return &model.ProductModel{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil
}

func (s *productService) DeleteByID(ctx context.Context, productID string) error {
	if err := s.productStore.DeleteByID(ctx, productID); err != nil {
		s.log.Error("failed to delete product by id", zap.String("product_id", productID), zap.Error(err))
		return err
	}

	return nil
}

func (s *productService) FindByID(ctx context.Context, productID string) (*model.ProductModel, error) {
	product, err := s.productStore.FindByID(ctx, productID)
	if err != nil || product == nil {
		s.log.Error("failed to find product by id", zap.String("product_id", productID), zap.Error(err))
		return nil, err
	}

	return &model.ProductModel{
		ID:        product.ID,
		Name:      product.Name,
		Price:     product.Price,
		Stock:     product.Stock,
		CreatedAt: product.CreatedAt,
		UpdatedAt: product.UpdatedAt,
	}, nil
}

func (s *productService) DeductStockByID(ctx context.Context, productID string, quantity int) error {
	product, err := s.productStore.FindByID(ctx, productID)
	if err != nil || product == nil {
		s.log.Error("failed to find product by id", zap.String("product_id", productID), zap.Error(err))
		return err
	}

	product.Stock = product.Stock - quantity

	product, err = s.productStore.UpdateByID(ctx, productID, product)
	if err != nil || product == nil {
		s.log.Error("failed to update product by id", zap.String("product_id", productID), zap.Error(err))
		return err
	}

	return nil
}

func (s *productService) RollbackStockByID(ctx context.Context, productID string, quantity int) error {
	product, err := s.productStore.FindByID(ctx, productID)
	if err != nil || product == nil {
		s.log.Error("failed to find product by id", zap.String("product_id", productID), zap.Error(err))
		return err
	}

	product.Stock = product.Stock + quantity

	product, err = s.productStore.UpdateByID(ctx, productID, product)
	if err != nil || product == nil {
		s.log.Error("failed to update product by id", zap.String("product_id", productID), zap.Error(err))
		return err
	}

	return nil
}
