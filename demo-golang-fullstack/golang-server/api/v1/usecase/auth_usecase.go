package usecaseV1

import (
	"context"

	payloadV1 "example.com/backend/api/v1/payload"
	"example.com/backend/internal/domain"
	"example.com/backend/internal/exception"
	"example.com/backend/internal/repository"
	"example.com/backend/internal/service"
	"example.com/backend/pkg/password"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=IAuthUseCaseV1
type IAuthUseCaseV1 interface {
	SignUp(ctx context.Context, param payloadV1.SignUpRequest) (*payloadV1.SignUpResponse, error)
	SignIn(ctx context.Context, param payloadV1.SignInRequest) (*payloadV1.SignInResponse, error)
}

type authUseCaseV1 struct {
	userRepo   repository.IUserRepository
	jwtService service.IJwtService
	log        zerolog.Logger
}

func NewAuthUseCaseV1(
	userRepo repository.IUserRepository,
	jwtService service.IJwtService,
	log zerolog.Logger,
) IAuthUseCaseV1 {
	return &authUseCaseV1{
		userRepo:   userRepo,
		jwtService: jwtService,
		log:        log,
	}
}

func (uc *authUseCaseV1) SignUp(ctx context.Context, param payloadV1.SignUpRequest) (*payloadV1.SignUpResponse, error) {
	exists, err := uc.userRepo.FindByEmail(ctx, param.Email)
	if err == nil || exists != nil {
		uc.log.Warn().Msg("email already registered")
		return nil, exception.NewConflict("Email already registered", nil)
	}

	defaultImage := "default.png"
	hashPassword, err := password.GenerateHash(param.Password)
	if err != nil {
		uc.log.Error().Err(err).Msg("failed to generate hash from password")
		return nil, exception.NewInternal("Failed to sign up user", nil)
	}

	user := &domain.User{
		Name:       param.Name,
		Email:      param.Email,
		Password:   hashPassword,
		Role:       string(domain.RoleUser),
		Occupation: &param.Occupation,
		ImageName:  &defaultImage,
	}

	user, err = uc.userRepo.Create(ctx, user)
	if err != nil {
		uc.log.Error().Err(err).Msg("failed when saving user in database")
		return nil, exception.NewInternal("Failed to sign up user", nil)
	}

	return &payloadV1.SignUpResponse{ID: user.ID}, nil
}

func (uc *authUseCaseV1) SignIn(ctx context.Context, param payloadV1.SignInRequest) (*payloadV1.SignInResponse, error) {
	user, err := uc.userRepo.FindByEmail(ctx, param.Email)
	if err != nil || user == nil {
		uc.log.Warn().Err(err).Msg("email is not registered")
		return nil, exception.NewBadRequest("Email is not registered", err)
	}

	if !password.CompareHash(param.Password, user.Password) {
		uc.log.Warn().Msg("password does not match")
		return nil, exception.NewBadRequest("Password does not match", nil)
	}

	tokenString, err := uc.jwtService.Generate(user.ID)
	if err != nil {
		uc.log.Error().Err(err).Msg("failed to generate token string")
		return nil, exception.NewInternal("Failed to sign in user", err)
	}

	return &payloadV1.SignInResponse{Token: tokenString}, nil
}
