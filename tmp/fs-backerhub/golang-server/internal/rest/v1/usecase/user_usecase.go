package usecaseV1

import (
	"context"

	"example.com.backend/internal/repos"
	"example.com.backend/pkg/exception"
	"example.com.backend/pkg/logger"
	"go.uber.org/zap"

	payloadV1 "example.com.backend/internal/rest/v1/payload"
)

type IUserUseCaseV1 interface {
	GetMe(ctx context.Context, userID string) (*payloadV1.UserResponse, error)
}

type userUseCaseV1 struct {
	userRepo repos.IUserRepository
}

func NewUserUseCaseV1(
	userRepo repos.IUserRepository,
) IUserUseCaseV1 {
	return &userUseCaseV1{userRepo: userRepo}
}

func (uc *userUseCaseV1) GetMe(ctx context.Context, userID string) (*payloadV1.UserResponse, error) {
	log := logger.GetLoggerFromContext(ctx)

	user, err := uc.userRepo.FindByID(ctx, userID)
	if err != nil || user == nil {
		log.Error("failed fetching user by id", zap.String("user_id", userID), zap.Error(err))
		return nil, exception.NewInternal("Failed to get user", err)
	}

	log.Info("successfully fetched user", zap.String("user_id", userID), zap.String("user_email", user.Email))
	return &payloadV1.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		ImageName: *user.ImageName,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil

}
