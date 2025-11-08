package user

import (
	"context"
	"fmt"
	"time"

	"example.com/user/pkg/config"
	"example.com/user/pkg/exception"
	"github.com/golang-jwt/jwt/v5"
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
	cfg   *config.Config
	log   *zerolog.Logger
}

func NewService(store IStore, cfg *config.Config, log *zerolog.Logger) IService {
	return &service{
		store: store,
		cfg:   cfg,
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

func (s *service) SignIn(ctx context.Context, param SignInRequest) (*UserWithTokenResponse, error) {
	user, err := s.store.FindByEmail(ctx, param.Email)
	if err != nil || user == nil {
		s.log.Warn().Err(err).Str("email", param.Email).Msg("Failed to find user by email")
		return nil, exception.NewBadRequest("Email not registered", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(param.Password)); err != nil {
		err := fmt.Errorf("compare hash and password error")
		s.log.Warn().Err(err).Str("email", param.Email).Msg("Failed to compare password")
		return nil, exception.NewBadRequest("Password does not match", err)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * time.Duration(s.cfg.JWT.TTL)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.Secret))
	if err != nil {
		s.log.Error().Err(err).Str("email", param.Email).Msg("Failed to generate jwt")
		return nil, exception.NewInternal("Failed to sign in user", err)
	}

	var userResp UserResponse
	userResp.FromUser(user)

	return &UserWithTokenResponse{
		Token: tokenString,
		User:  userResp,
	}, nil
}

func (s *service) FindByID(ctx context.Context, userID string) (*UserResponse, error) {
	user, err := s.store.FindByID(ctx, userID)
	if err != nil || user == nil {
		s.log.Warn().Err(err).Str("user_id", userID).Msg("Failed to find user by id")
		return nil, exception.NewBadRequest("User not found", err)
	}

	var userResp UserResponse
	userResp.FromUser(user)

	return &userResp, nil
}
