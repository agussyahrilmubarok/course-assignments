package usecaseV1

import (
	"context"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/repos"
	"example.com.backend/internal/service"
	"example.com.backend/pkg/exception"
	"example.com.backend/pkg/logger"
	"example.com.backend/pkg/password"
	"go.uber.org/zap"

	payloadV1 "example.com.backend/internal/rest/v1/payload"
)

type IAuthUseCaseV1 interface {
	SignUp(ctx context.Context, param payloadV1.SignUpRequest) (*payloadV1.SignUpResponse, error)
	SignIn(ctx context.Context, param payloadV1.SignInRequest) (*payloadV1.SignInResponse, error)
}

type authUseCaseV1 struct {
	userRepo   repos.IUserRepository
	jwtService service.IJwtService
}

func NewAuthUseCaseV1(
	userRepo repos.IUserRepository,
	jwtService service.IJwtService,
) IAuthUseCaseV1 {
	return &authUseCaseV1{
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}

func (uc *authUseCaseV1) SignUp(ctx context.Context, param payloadV1.SignUpRequest) (*payloadV1.SignUpResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	exists, _ := uc.userRepo.ExistsByEmailIgnoreCase(ctx, param.Email)
	if exists {
		log.Warn("email already registered", zap.String("user_email", param.Email))
		return nil, exception.NewConflict("Email already registered", nil)
	}

	hashPassword, err := password.GenerateHash(param.Password)
	if err != nil {
		log.Error("failed to generate hash from password", zap.Error(err))
		return nil, exception.NewInternal("Failed to sign up user", err)
	}

	defaultImage := "default.png"
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
		log.Error("failed when saving user in database", zap.String("user_email", param.Email), zap.Error(err))
		return nil, exception.NewInternal("Failed to sign up user", err)
	}

	log.Info("successfully created user", zap.String("user_id", user.ID), zap.String("user_email", user.Email))
	return &payloadV1.SignUpResponse{ID: user.ID}, nil
}

func (uc *authUseCaseV1) SignIn(ctx context.Context, param payloadV1.SignInRequest) (*payloadV1.SignInResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	user, err := uc.userRepo.FindByEmail(ctx, param.Email)
	if err != nil || user == nil {
		log.Warn("email is not registered", zap.String("user_email", param.Email))
		return nil, exception.NewBadRequest("Email is not registered", err)
	}

	if !password.CompareHash(param.Password, user.Password) {
		log.Warn("password does not match", zap.String("user_email", param.Email))
		return nil, exception.NewBadRequest("Password does not match", nil)
	}

	tokenString, err := uc.jwtService.Generate(ctx, user.ID)
	if err != nil {
		log.Error("failed to generate token string", zap.String("user_id", user.ID), zap.Error(err))
		return nil, exception.NewInternal("Failed to sign in user", err)
	}

	log.Info("successfully signed in user", zap.String("user_id", user.ID), zap.String("user_email", user.Email))
	return &payloadV1.SignInResponse{Token: tokenString}, nil
}
