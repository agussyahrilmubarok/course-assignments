package usecaseV1

import (
	"context"

	payloadV1 "example.com/backend/api/v1/payload"
	"example.com/backend/internal/exception"
	"example.com/backend/internal/repository"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=IUserUseCaseV1
type IUserUseCaseV1 interface {
	GetMe(ctx context.Context, userID string) (*payloadV1.UserResponse, error)
}

type userUseCaseV1 struct {
	userRepo repository.IUserRepository
	log      zerolog.Logger
}

func NewUserUseCaseV1(
	userRepo repository.IUserRepository,
	log zerolog.Logger,
) IUserUseCaseV1 {
	return &userUseCaseV1{
		userRepo: userRepo,
		log:      log,
	}
}

func (uc *userUseCaseV1) GetMe(ctx context.Context, userID string) (*payloadV1.UserResponse, error) {
	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		uc.log.Error().Err(err).Msg("failed to find user by id")
		return nil, exception.NewNotFound("User is not found", err)
	}

	return &payloadV1.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		ImageName: *user.ImageName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
