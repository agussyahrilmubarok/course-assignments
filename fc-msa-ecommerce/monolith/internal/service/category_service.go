package service

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/model"
	"ecommerce/internal/repository"
	"errors"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=ICategoryService
type ICategoryService interface {
	GetAll(ctx context.Context) ([]domain.Category, error)
	GetByID(ctx context.Context, id uint) (*domain.Category, error)
	Create(ctx context.Context, request model.CreateCategoryRequest) (*domain.Category, error)
	Update(ctx context.Context, id uint, request model.UpdateCategoryRequest) (*domain.Category, error)
	Delete(ctx context.Context, id uint) error
}

type categoryService struct {
	CategoryRepository repository.ICategoryRepository
	Logger             zerolog.Logger
}

func NewCategoryService(
	categoryRepository repository.ICategoryRepository,
	logger zerolog.Logger,
) ICategoryService {
	return &categoryService{
		CategoryRepository: categoryRepository,
		Logger:             logger,
	}
}

func (s *categoryService) GetAll(ctx context.Context) ([]domain.Category, error) {
	return s.CategoryRepository.FindAll(ctx)
}

func (s *categoryService) GetByID(ctx context.Context, id uint) (*domain.Category, error) {
	return s.CategoryRepository.FindByID(ctx, id)
}

func (s *categoryService) Create(ctx context.Context, request model.CreateCategoryRequest) (*domain.Category, error) {
	exists := s.CategoryRepository.ExistsByName(ctx, request.Name)
	if exists {
		s.Logger.Warn().Str("name", request.Name).Msg("duplicate category name")
		return nil, errors.New("category name already exists")
	}

	category := &domain.Category{
		Name: request.Name,
	}

	return s.CategoryRepository.Save(ctx, category)
}

func (s *categoryService) Update(ctx context.Context, id uint, request model.UpdateCategoryRequest) (*domain.Category, error) {
	category, err := s.CategoryRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	// Check if new name conflicts with other categories
	if category.Name != request.Name {
		exists := s.CategoryRepository.ExistsByName(ctx, request.Name)
		if exists {
			s.Logger.Warn().Str("name", request.Name).Msg("duplicate category name")
			return nil, errors.New("category name already exists")
		}
	}

	category.Name = request.Name
	return s.CategoryRepository.Save(ctx, category)
}

func (s *categoryService) Delete(ctx context.Context, id uint) error {
	category, err := s.CategoryRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if category == nil {
		return errors.New("category not found")
	}
	return s.CategoryRepository.DeleteByID(ctx, id)
}
