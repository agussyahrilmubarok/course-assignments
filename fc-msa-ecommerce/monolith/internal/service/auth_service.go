package service

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/internal/model"
	"ecommerce/internal/repository"
	"ecommerce/pkg/helper"
	"errors"

	"github.com/rs/zerolog"
)

//go:generate mockery --name=IAuthService
type IAuthService interface {
	Register(ctx context.Context, request model.RegisterRequest) error
	Login(ctx context.Context, request model.LoginRequest) (string, error)
	Logout(ctx context.Context) error
}

type authService struct {
	UserRepository repository.IUserRepository
	Logger         zerolog.Logger
}

func NewAuthService(
	userRepository repository.IUserRepository,
	logger zerolog.Logger,
) IAuthService {
	return &authService{
		UserRepository: userRepository,
		Logger:         logger,
	}
}

func (s *authService) Register(ctx context.Context, request model.RegisterRequest) error {
	exists := s.UserRepository.ExistsByEmail(ctx, request.Email)
	if exists {
		s.Logger.Warn().Str("email", request.Email).Msg("duplicate user email")
		return errors.New("email already registered")
	}

	hashedPassword, err := helper.PasswordHash(request.Password)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to hash password")
		return err
	}

	user := &domain.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	_, err = s.UserRepository.Save(ctx, user)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to create user")
		return err
	}

	return nil
}

func (s *authService) Login(ctx context.Context, request model.LoginRequest) (string, error) {
	user, err := s.UserRepository.FindByEmail(ctx, request.Email)
	if err != nil || user == nil {
		s.Logger.Warn().Err(err).Msg("failed to find user by email")
		return "", errors.New("invalid email")
	}

	if err := helper.PasswordCompare(user.PasswordHash, request.Password); err != nil {
		s.Logger.Warn().Err(err).Msg("invalid password")
		return "", errors.New("invalid password")
	}

	token, err := helper.JWTGenerate(user.ID, user.Email, user.Role)
	if err != nil {
		s.Logger.Error().Err(err).Msg("failed to generate jwt")
		return "", errors.New("failed to generate jwt")
	}

	return token, nil
}

func (s *authService) Logout(ctx context.Context) error {
	// Stateless JWT: Logout is handled on the client side
	// Token blacklist can be implemented here if needed
	return nil
}
