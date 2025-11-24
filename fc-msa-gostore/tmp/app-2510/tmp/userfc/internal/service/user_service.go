package service

import (
	"context"
	"errors"
	"time"

	"example.com/pkg/model"
	"example.com/userfc/internal/store"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockery --name=IUserService
type IUserService interface {
	SignUp(ctx context.Context, req *model.SignUpRequest) (*model.SignUpResponse, error)
	SignIn(ctx context.Context, req *model.SignInRequest) (*model.SignInResponse, error)
	FindByID(ctx context.Context, userID string) (*model.UserModel, error)
}

const (
	TOKEN_SECRET = "S3cr3tK#yS3cr3tK#yS3cr3tK#yS3cr3tK#yS3cr3tK#yS3cr3tK#yS3cr3tK#y"
	TOKEN_TTL    = 24
)

type userService struct {
	userStore store.IUserStore
	log       *zap.Logger
}

func NewUserService(
	userStore store.IUserStore,
	log *zap.Logger,
) IUserService {
	return &userService{
		userStore: userStore,
		log:       log,
	}
}

func (s *userService) SignUp(ctx context.Context, req *model.SignUpRequest) (*model.SignUpResponse, error) {
	exists, err := s.userStore.ExistsEmail(ctx, req.Email)
	if err != nil || exists {
		s.log.Error("failed to sign up user, email already exists", zap.Error(err))
		return nil, errors.New("email is already exists")
	}

	userID := uuid.New().String()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error("failed to sign up user, hashed password", zap.Error(err))
		return nil, err
	}

	user := &store.User{
		ID:       userID,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     store.RoleUser,
	}
	user, err = s.userStore.Create(ctx, user)
	if err != nil || user == nil {
		s.log.Error("failed to sign up user, save user", zap.Error(err))
		return nil, err
	}

	return &model.SignUpResponse{
		ID: user.ID,
	}, nil
}

func (s *userService) SignIn(ctx context.Context, req *model.SignInRequest) (*model.SignInResponse, error) {
	user, err := s.userStore.FindByEmail(ctx, req.Email)
	if err != nil || user == nil {
		s.log.Error("failed to find user by email", zap.String("user_email", req.Email), zap.Error(err))
		return nil, errors.New("email is not registered")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		s.log.Error("failed to compare password", zap.String("user_email", req.Email), zap.Error(err))
		return nil, errors.New("password does not match")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * time.Duration(TOKEN_TTL)).Unix(),
	})
	tokenString, err := token.SignedString([]byte(TOKEN_SECRET))
	if err != nil {
		s.log.Error("failed to generate token string", zap.String("user_email", req.Email), zap.Error(err))
		return nil, errors.New("failed to generate token")
	}

	return &model.SignInResponse{
		Token: tokenString,
		User: model.UserModel{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      string(user.Role),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *userService) FindByID(ctx context.Context, userID string) (*model.UserModel, error) {
	user, err := s.userStore.FindByID(ctx, userID)
	if err != nil || user == nil {
		s.log.Error("failed to find user by id", zap.String("user_id", userID), zap.Error(err))
		return nil, err
	}

	return &model.UserModel{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
