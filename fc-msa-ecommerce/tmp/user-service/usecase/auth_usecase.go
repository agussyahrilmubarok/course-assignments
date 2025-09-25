package usecase

import (
	"context"
	"ecommerce/user-service/model"
	"ecommerce/user-service/service"
	"ecommerce/user-service/store"
	"errors"
	"strings"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyExists = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

//go:generate mockery --name=IAuthUseCase
type IAuthUseCase interface {
	SignUp(ctx context.Context, param model.SignUpRequest) (*model.SignUpResponse, error)
	SignIn(ctx context.Context, param model.SignInRequest) (*model.SignInResponse, error)
}

type authUseCase struct {
	userMongoStore store.IUserMongoStore
	jwtService     service.IJWTService
	log            zerolog.Logger
}

func NewAuthUseCase(
	userMongoStore store.IUserMongoStore,
	jwtService service.IJWTService,
	log zerolog.Logger,
) IAuthUseCase {
	return &authUseCase{
		userMongoStore: userMongoStore,
		jwtService:     jwtService,
		log:            log,
	}
}

func (u *authUseCase) SignUp(ctx context.Context, param model.SignUpRequest) (*model.SignUpResponse, error) {
	name := strings.TrimSpace(param.Name)
	email := strings.ToLower(strings.TrimSpace(param.Email))

	_, err := u.userMongoStore.FindByEmail(ctx, email)
	if err == nil {
		u.log.Warn().Str("email", email).Msg("email already registered")
		return nil, ErrEmailAlreadyExists
	}

	passwordHashed, err := bcrypt.GenerateFromPassword([]byte(param.Password), bcrypt.DefaultCost)
	if err != nil {
		u.log.Error().Err(err).Msg("failed to hash password")
		return nil, err
	}

	userStore, err := u.userMongoStore.Create(ctx, &store.User{
		Name:     name,
		Email:    email,
		Password: string(passwordHashed),
	})
	if err != nil {
		u.log.Error().Err(err).Msg("failed to create user")
		return nil, err
	}

	return &model.SignUpResponse{
		ID:    userStore.ID.Hex(),
		Name:  userStore.Name,
		Email: userStore.Email,
	}, nil
}

func (u *authUseCase) SignIn(ctx context.Context, param model.SignInRequest) (*model.SignInResponse, error) {
	email := strings.ToLower(strings.TrimSpace(param.Email))

	userStore, err := u.userMongoStore.FindByEmail(ctx, email)
	if err != nil {
		u.log.Warn().Str("email", email).Msg("invalid email or password")
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userStore.Password), []byte(param.Password)); err != nil {
		u.log.Warn().Str("email", email).Msg("invalid email or password")
		return nil, ErrInvalidCredentials
	}

	token, err := u.jwtService.GenerateToken(userStore.ID.Hex())
	if err != nil {
		u.log.Error().Err(err).Msg("failed to generate token")
		return nil, err
	}

	return &model.SignInResponse{
		Token: token,
	}, nil
}
