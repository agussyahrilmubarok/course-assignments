package service

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/model"
	"ecommerce/internal/repository"
	"errors"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=ITagService
type ITagService interface {
	GetAll(ctx context.Context) ([]domain.Tag, error)
	GetByID(ctx context.Context, id uint) (*domain.Tag, error)
	Create(ctx context.Context, request model.CreateTagRequest) (*domain.Tag, error)
	Update(ctx context.Context, id uint, request model.UpdateTagRequest) (*domain.Tag, error)
	Delete(ctx context.Context, id uint) error
}

type TagService struct {
	TagRepository repository.ITagRepository
	Logger        zerolog.Logger
}

func NewTagService(
	tagRepository repository.ITagRepository,
	logger zerolog.Logger,
) ITagService {
	return &TagService{
		TagRepository: tagRepository,
		Logger:        logger,
	}
}

func (s *TagService) GetAll(ctx context.Context) ([]domain.Tag, error) {
	return s.TagRepository.FindAll(ctx)
}

func (s *TagService) GetByID(ctx context.Context, id uint) (*domain.Tag, error) {
	return s.TagRepository.FindByID(ctx, id)
}

func (s *TagService) Create(ctx context.Context, request model.CreateTagRequest) (*domain.Tag, error) {
	exists := s.TagRepository.ExistsByName(ctx, request.Name)
	if exists {
		s.Logger.Warn().Str("name", request.Name).Msg("duplicate tag name")
		return nil, errors.New("tag name already exists")
	}

	Tag := &domain.Tag{
		Name: request.Name,
	}

	return s.TagRepository.Save(ctx, Tag)
}

func (s *TagService) Update(ctx context.Context, id uint, request model.UpdateTagRequest) (*domain.Tag, error) {
	Tag, err := s.TagRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if Tag == nil {
		return nil, errors.New("tag not found")
	}

	// Check if new name conflicts with other categories
	if Tag.Name != request.Name {
		exists := s.TagRepository.ExistsByName(ctx, request.Name)
		if exists {
			s.Logger.Warn().Str("name", request.Name).Msg("duplicate tag name")
			return nil, errors.New("tag name already exists")
		}
	}

	Tag.Name = request.Name
	return s.TagRepository.Save(ctx, Tag)
}

func (s *TagService) Delete(ctx context.Context, id uint) error {
	Tag, err := s.TagRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if Tag == nil {
		return errors.New("tag not found")
	}
	return s.TagRepository.DeleteByID(ctx, id)
}
