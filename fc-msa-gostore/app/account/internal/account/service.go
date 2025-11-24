package account

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

type IService interface {
	SignUp(ctx context.Context, request SignUpRequest) (*AccountResponse, error)
	SignIn(ctx context.Context, request SignInRequest) (*AccountWithTokenResponse, error)
	FindByID(ctx context.Context, accountID string) (*AccountResponse, error)
}

type service struct {
	repository IRepository
	logger     *logrus.Logger
}

func NewService(repository IRepository, logger *logrus.Logger) IService {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s *service) SignUp(ctx context.Context, request SignUpRequest) (*AccountResponse, error) {
	account, err := s.repository.FindByEmail(ctx, request.Email)
	if account != nil || err == nil {
		s.logger.Warnf("failed signup email exists: %v", err)
		return nil, errors.New("email already exists")
	}

	panic("unimplemented")
}

// FindByID implements IService.
func (s *service) FindByID(ctx context.Context, accountID string) (*AccountResponse, error) {
	panic("unimplemented")
}

// SignIn implements IService.
func (s *service) SignIn(ctx context.Context, request SignInRequest) (*AccountWithTokenResponse, error) {
	panic("unimplemented")
}
