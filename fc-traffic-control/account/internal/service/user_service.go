package service

import (
	"context"
	"errors"
	"time"

	"traffic-control/account/internal/domain"
	"traffic-control/account/internal/repository"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

// IUserService defines user-related business logic (use-case layer)
type IUserService interface {
	FindAll(ctx context.Context) ([]*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Create(ctx context.Context, user *domain.User) error
	Update(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id string) error
}

// userService is the concrete implementation
type userService struct {
	userRepo repository.IUserRepository
	log      zerolog.Logger
}

// NewUserService creates a new instance of userService
func NewUserService(userRepo repository.IUserRepository, log zerolog.Logger) IUserService {
	return &userService{
		userRepo: userRepo,
		log:      log.With().Str("component", "service.user").Logger(),
	}
}

// FindAll retrieves all users
func (s *userService) FindAll(ctx context.Context) ([]*domain.User, error) {
	users, err := s.userRepo.FindAll(ctx)
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to fetch users")
		return nil, err
	}
	return users, nil
}

// FindByID retrieves a user by ID
func (s *userService) FindByID(ctx context.Context, id string) (*domain.User, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("Failed to find user")
		return nil, err
	}

	if user == nil {
		s.log.Warn().Str("id", id).Msg("User not found")
		return nil, nil
	}

	return user, nil
}

// Create validates input, hashes password, and stores the user
func (s *userService) Create(ctx context.Context, user *domain.User) error {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return errors.New("name, email, and password are required")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		s.log.Error().Err(err).Str("email", user.Email).Msg("Failed to hash password")
		return err
	}
	user.Password = string(hashedPassword)

	// Set timestamps
	now := time.Now().UTC().Format(time.RFC3339)
	user.CreatedAt = now
	user.UpdatedAt = now

	// Store in DB
	if err := s.userRepo.Create(ctx, user); err != nil {
		s.log.Error().Err(err).Str("email", user.Email).Msg("Failed to create user")
		return err
	}

	s.log.Info().Str("email", user.Email).Msg("User created successfully")
	return nil
}

// Update modifies an existing user (rehash password if changed)
func (s *userService) Update(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		return errors.New("user ID is required for update")
	}

	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			s.log.Error().Err(err).Str("id", user.ID).Msg("Failed to hash password on update")
			return err
		}
		user.Password = string(hashedPassword)
	}

	user.UpdatedAt = time.Now().UTC().Format(time.RFC3339)

	if err := s.userRepo.Update(ctx, user); err != nil {
		s.log.Error().Err(err).Str("id", user.ID).Msg("Failed to update user")
		return err
	}

	s.log.Info().Str("id", user.ID).Msg("User updated successfully")
	return nil
}

// Delete removes a user by ID
func (s *userService) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("Failed to delete user")
		return err
	}

	s.log.Info().Str("id", id).Msg("User deleted successfully")
	return nil
}
