package service

import (
	"context"

	"example.com/backend/internal/domain"
	"example.com/backend/internal/model"
	"example.com/backend/internal/repository"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=IUserService
type IUserService interface {
	FindAll(ctx context.Context) ([]model.UserDTO, error)
	FindByID(ctx context.Context, id string) (*model.UserDTO, error)
	FindByEmail(ctx context.Context, email string) (*model.UserDTO, error)
	Create(ctx context.Context, userDto model.UserDTO) error
	Update(ctx context.Context, userDto model.UserDTO) error
	DeleteByID(ctx context.Context, id string) error
}

type userService struct {
	userRepo repository.IUserRepository
	log      zerolog.Logger
}

func NewUserService(
	userRepo repository.IUserRepository,
	log zerolog.Logger,
) IUserService {
	return &userService{
		userRepo: userRepo,
		log:      log,
	}
}

func (s *userService) FindAll(ctx context.Context) ([]model.UserDTO, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		s.log.Error().Err(err).Msgf("failed to find users")
		return nil, err
	}

	var userDtos []model.UserDTO
	for _, user := range users {
		var userDto model.UserDTO
		userDto.FromUser(&user)
		userDtos = append(userDtos, userDto)
	}

	return userDtos, nil
}

func (s *userService) FindByID(ctx context.Context, id string) (*model.UserDTO, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil || user == nil {
		s.log.Error().Err(err).Msgf("failed to find user by id %s", id)
		return nil, err
	}

	var userDto model.UserDTO
	userDto.FromUser(user)

	return &userDto, nil
}

func (s *userService) FindByEmail(ctx context.Context, email string) (*model.UserDTO, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		s.log.Error().Err(err).Msgf("failed to find user by email %s", email)
		return nil, err
	}

	var userDto model.UserDTO
	userDto.FromUser(user)

	return &userDto, nil
}

func (s *userService) Create(ctx context.Context, userDto model.UserDTO) error {
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
		s.log.Error().Err(err).Msgf("failed to save user with email %s", userDto.Email)
		return err
	}

	return nil
}

func (s *userService) Update(ctx context.Context, userDto model.UserDTO) error {
	userDto.HashPassword(userDto.Password)

	user, err := s.userRepo.FindByID(ctx, userDto.ID)
	if err != nil || user == nil {
		s.log.Error().Err(err).Msgf("failed to find user by id %s", userDto.ID)
		return err
	}

	user.Name = userDto.Name
	user.Email = userDto.Email
	user.Occupation = &userDto.Occupation
	user.ImageName = &userDto.ImageName

	_, err = s.userRepo.Update(ctx, user)
	if err != nil {
		s.log.Error().Err(err).Msgf("failed to update user by id %s", userDto.ID)
		return err
	}

	return nil
}

func (s *userService) DeleteByID(ctx context.Context, id string) error {
	err := s.userRepo.DeleteByID(ctx, id)
	if err != nil {
		s.log.Error().Err(err).Msgf("failed to delete user by id %s", id)
		return err
	}

	return nil
}
