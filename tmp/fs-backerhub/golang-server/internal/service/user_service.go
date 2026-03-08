package service

import (
	"context"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/model"
	"example.com.backend/internal/repos"
	"example.com.backend/pkg/logger"
	"go.uber.org/zap"
)

type IUserService interface {
	FindAll(ctx context.Context) ([]model.UserDTO, error)
	FindByID(ctx context.Context, id string) (*model.UserDTO, error)
	FindByEmail(ctx context.Context, email string) (*model.UserDTO, error)
	Create(ctx context.Context, userDto model.UserDTO) error
	Update(ctx context.Context, userDto model.UserDTO) error
	DeleteByID(ctx context.Context, id string) error
	ExistsByEmailIgnoreCase(ctx context.Context, email string) (bool, error)
}

type userService struct {
	userRepo repos.IUserRepository
}

func NewUserService(userRepo repos.IUserRepository) IUserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) FindAll(ctx context.Context) ([]model.UserDTO, error) {
	log := logger.GetLoggerFromContext(ctx)

	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		log.Error("failed fetching users", zap.Error(err))
		return nil, err
	}

	var userDtos []model.UserDTO
	for _, user := range users {
		var dto model.UserDTO
		dto.FromUser(&user)
		userDtos = append(userDtos, dto)
	}

	log.Info("successfully fetched users", zap.Int("count", len(users)))
	return userDtos, nil
}

func (s *userService) FindByID(ctx context.Context, id string) (*model.UserDTO, error) {
	log := logger.GetLoggerFromContext(ctx)

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		log.Error("failed fetching user by id", zap.String("user_id", id), zap.Error(err))
		return nil, err
	}

	if user == nil {
		log.Warn("user not found by id", zap.String("user_id", id))
		return nil, nil
	}

	var dto model.UserDTO
	dto.FromUser(user)

	log.Info("successfully fetched user by id", zap.String("user_id", id))
	return &dto, nil
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*model.UserDTO, error) {
	log := logger.GetLoggerFromContext(ctx)

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Error("failed fetching user by email", zap.String("user_email", email), zap.Error(err))
		return nil, err
	}

	if user == nil {
		log.Warn("user not found by email", zap.String("user_email", email))
		return nil, nil
	}

	var dto model.UserDTO
	dto.FromUser(user)

	log.Info("successfully fetched user by email", zap.String("user_email", email))
	return &dto, nil
}

func (s *userService) Create(ctx context.Context, userDto model.UserDTO) error {
	log := logger.GetLoggerFromContext(ctx)

	userDto.HashPassword(userDto.Password)
	defaultImage := "default.png"

	user := &domain.User{
		Name:       userDto.Name,
		Email:      userDto.Email,
		Role:       string(domain.RoleUser),
		Occupation: &userDto.Occupation,
		ImageName:  &defaultImage,
	}

	_, err := s.userRepo.Create(ctx, user)
	if err != nil {
		log.Error("failed creating user", zap.String("user_email", userDto.Email), zap.Error(err))
		return err
	}

	log.Info("successfully created user",
		zap.String("user_id", user.ID),
		zap.String("user_email", user.Email),
	)
	return nil
}

func (s *userService) Update(ctx context.Context, userDto model.UserDTO) error {
	log := logger.GetLoggerFromContext(ctx)

	userDto.HashPassword(userDto.Password)

	user, err := s.userRepo.FindByID(ctx, userDto.ID)
	if err != nil {
		log.Error("failed fetching user for update", zap.String("user_id", userDto.ID), zap.Error(err))
		return err
	}
	if user == nil {
		log.Warn("user not found for update", zap.String("user_id", userDto.ID))
		return nil
	}

	user.Name = userDto.Name
	user.Email = userDto.Email
	user.Occupation = &userDto.Occupation
	user.ImageName = &userDto.ImageName

	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		log.Error("failed updating user", zap.String("user_id", userDto.ID), zap.Error(err))
		return err
	}

	log.Info("successfully updated user", zap.String("user_id", userDto.ID))
	return nil
}

func (s *userService) DeleteByID(ctx context.Context, id string) error {
	log := logger.GetLoggerFromContext(ctx)

	err := s.userRepo.DeleteByID(ctx, id)
	if err != nil {
		log.Error("failed deleting user", zap.String("user_id", id), zap.Error(err))
		return err
	}

	log.Info("successfully deleted user", zap.String("user_id", id))
	return nil
}

func (s *userService) ExistsByEmailIgnoreCase(ctx context.Context, email string) (bool, error) {
	log := logger.GetLoggerFromContext(ctx)

	exists, err := s.userRepo.ExistsByEmailIgnoreCase(ctx, email)
	if err != nil {
		log.Error("failed checking existing email", zap.String("user_email", email), zap.Error(err))
		return false, err
	}

	log.Info("checked existing email", zap.String("user_email", email), zap.Bool("exists", exists))
	return exists, nil
}
