package service

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/model"
	"ecommerce/internal/repository"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=IProductService
type IProductService interface {
	GetAll(ctx context.Context) ([]domain.Product, error)
	GetByID(ctx context.Context, id uint) (*domain.Product, error)
	Create(ctx context.Context, request model.CreateProductRequest) (*domain.Product, error)
	Update(ctx context.Context, request model.UpdateProductRequest) (*domain.Product, error)
	Delete(ctx context.Context, id uint) error
}

type productService struct {
	ProductRepository repository.IProductRepository
	Logger            zerolog.Logger
}

func NewProductService(
	productRepository repository.IProductRepository,
	logger zerolog.Logger,
) IProductService {
	return &productService{
		ProductRepository: productRepository,
		Logger:            logger,
	}
}
