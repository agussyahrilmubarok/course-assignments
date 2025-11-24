package service

import (
	"context"

	pb "example.com/pkg/proto/v1"
	"example.com/userfc/internal/store"
	"go.uber.org/zap"
)

type userGrpcService struct {
	pb.UnimplementedUserServiceServer
	userStore store.IUserStore
	log       *zap.Logger
}

func NewUserGrpcService(userStore store.IUserStore, log *zap.Logger) *userGrpcService {
	return &userGrpcService{
		userStore: userStore,
		log:       log,
	}
}

func (s *userGrpcService) GetUserByID(ctx context.Context, req *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.userStore.FindByID(ctx, req.Id)
	if err != nil || user == nil {
		s.log.Error("failed to find user by id", zap.String("user_id", req.Id), zap.Error(err))
		return nil, err
	}

	return &pb.UserResponse{
		Id:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}, nil
}
