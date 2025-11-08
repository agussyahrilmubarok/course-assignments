package user

import (
	"context"

	"example.com/user/pkg/exception"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockery --name=IService
type IService interface {
	SignUp(ctx context.Context, param SignUpRequest) (*UserResponse, error)
	SignIn(ctx context.Context, param SignInRequest) (*UserWithTokenResponse, error)
	FindByID(ctx context.Context, userID string) (*UserResponse, error)
}

type service struct {
	store IStore
	log   *zerolog.Logger
}

func NewService(store IStore, log *zerolog.Logger) IService {
	return &service{
		store: store,
		log:   log,
	}
}

func (s *service) SignUp(ctx context.Context, param SignUpRequest) (*UserResponse, error) {
	exists, _ := s.store.FindByEmail(ctx, param.Email)
	if exists != nil {
		s.log.Warn().Str("email", param.Email).Msg("Email already used")
		return nil, exception.NewBadRequest("Email already used", nil)
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to hash password")
		return nil, exception.NewBadRequest("Failed to sign up user", err)
	}

	var user User
	user.Name = param.Name
	user.Email = param.Email
	user.Password = string(hashPass)
	if err := s.store.Create(ctx, &user); err != nil {
		s.log.Error().Err(err).Msg("Failed to save user")
		return nil, exception.NewInternal("Failed to sign up user", err)
	}

	var userResp UserResponse
	userResp.FromUser(&user)

	return &userResp, nil
}

// FindByID implements IService.
func (s *service) FindByID(ctx context.Context, userID string) (*UserResponse, error) {
	panic("unimplemented")
}

// SignIn implements IService.
func (s *service) SignIn(ctx context.Context, param SignInRequest) (*UserWithTokenResponse, error) {
	panic("unimplemented")
}
