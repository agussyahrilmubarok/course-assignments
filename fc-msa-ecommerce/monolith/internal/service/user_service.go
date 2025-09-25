package service

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/model"
	"ecommerce/internal/repository"
	"errors"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=IUserService
type IUserService interface {
	GetAll(ctx context.Context) ([]domain.User, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	Update(ctx context.Context, id uint, request model.UpdateUserRequest) (*domain.User, error)
	Delete(ctx context.Context, id uint) error
}

type userService struct {
	UserRepository repository.IUserRepository
	Logger         zerolog.Logger
}

func NewUserService(userRepo repository.IUserRepository, logger zerolog.Logger) IUserService {
	return &userService{
		UserRepository: userRepo,
		Logger:         logger,
	}
}

func (s *userService) GetAll(ctx context.Context) ([]domain.User, error) {
	users, err := s.UserRepository.FindAll(ctx)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to get all users")
		return nil, err
	}
	return users, nil
}

func (s *userService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	user, err := s.UserRepository.FindByID(ctx, id)
	if err != nil {
		s.Logger.Error().Err(err).Uint("id", id).Msg("failed to get user by ID")
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) Update(ctx context.Context, id uint, request model.UpdateUserRequest) (*domain.User, error) {
	user, err := s.UserRepository.FindByID(ctx, id)
	if err != nil {
		s.Logger.Error().Err(err).Uint("id", id).Msg("failed to find user")
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.Email != request.Email {
		exists := s.UserRepository.ExistsByEmail(ctx, request.Email)
		if exists {
			s.Logger.Warn().Str("email", request.Email).Msg("duplicate user email")
			return nil, errors.New("email already registered")
		}
	}

	user.Name = request.Name
	user.Email = request.Email

	updatedUser, err := s.UserRepository.Save(ctx, user)
	if err != nil {
		s.Logger.Error().Err(err).Uint("id", id).Msg("failed to update user")
		return nil, err
	}

	return updatedUser, nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	user, err := s.UserRepository.FindByID(ctx, id)
	if err != nil {
		s.Logger.Error().Err(err).Uint("id", id).Msg("failed to find user for delete")
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	err = s.UserRepository.DeleteByID(ctx, id)
	if err != nil {
		s.Logger.Error().Err(err).Uint("id", id).Msg("failed to delete user")
		return err
	}

	return nil
}
