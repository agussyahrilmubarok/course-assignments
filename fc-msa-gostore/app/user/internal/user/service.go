package user

import (
	"context"
	"errors"

	"github.com/agussyahrilmubarok/gox/pkg/xexception"
	"github.com/agussyahrilmubarok/gox/pkg/xjwt"
	"github.com/agussyahrilmubarok/gox/pkg/xpassword"
	"github.com/rs/zerolog"
)

type IService interface {
	SignUp(ctx context.Context, param SignUpParam) (UserResponse, error)
	SignIn(ctx context.Context, param SignInParam) (UserWithTokenResponse, error)
}

type service struct {
	store  IStore
	logger zerolog.Logger
}

func NewServie(
	store IStore,
	logger zerolog.Logger,
) IService {
	return &service{
		store:  store,
		logger: logger,
	}
}

func (s *service) SignUp(ctx context.Context, param SignUpParam) (UserResponse, error) {
	err := s.store.ExistsUserEmailByIgnoreCase(ctx, param.Email)
	if err != nil {
		err := errors.New("email already exists")
		s.logger.Error().
			Err(err).
			Str("email", param.Email).
			Msg("user email already exists")
		return UserResponse{}, xexception.NewHTTPBadRequest("email already exits", nil)
	}

	hashPass, saltPass, err := xpassword.PBKDF2Hash(param.Password)
	if err != nil {
		s.logger.Error().
			Err(err).
			Str("email", param.Email).
			Msg("failed to hash password")
		return UserResponse{}, xexception.NewHTTPInternal("failed to sign up user", err)
	}

	user := &User{
		Name:         param.Name,
		Email:        param.Email,
		HashPassword: hashPass,
		SaltPassword: saltPass,
	}

	if err := s.store.CreateUser(ctx, user); err != nil {
		s.logger.Error().
			Err(err).
			Str("email", param.Email).
			Msg("failed to create user")
		return UserResponse{}, xexception.NewHTTPInternal("failed to sign up user", err)
	}

	var res UserResponse
	res.From(user)

	s.logger.Info().
		Str("email", param.Email).
		Msg("signup user successfully")
	return res, nil
}

func (s *service) SignIn(ctx context.Context, param SignInParam) (UserWithTokenResponse, error) {
	user, err := s.store.FindUserByEmail(ctx, param.Email)
	if err != nil || user == nil {
		s.logger.Error().
			Err(err).
			Str("email", param.Email).
			Msg("email is not found")
		return UserWithTokenResponse{}, xexception.NewHTTPBadRequest("email not found", nil)
	}

	if err := xpassword.PBKDF2Compare(param.Password, user.HashPassword, user.SaltPassword); err != nil {
		s.logger.Error().
			Err(err).
			Str("email", param.Email).
			Msg("password does not match")
		return UserWithTokenResponse{}, xexception.NewHTTPBadRequest("password does not match", nil)
	}

	claims := map[string]interface{}{
		"user_id": user.ID.Hex(),
	}
	xjwt.SecretKey = []byte("secret_secret_secret_secret")
	tokenString, err := xjwt.Generate(claims, 600)
	if err != nil {
		s.logger.Error().
			Err(err).
			Msg("failed to generate jwt")
		return UserWithTokenResponse{}, xexception.NewHTTPInternal("sign in failed", nil)
	}

	var res UserWithTokenResponse
	res.Token = tokenString
	res.User.From(user)

	s.logger.Info().
		Str("email", param.Email).
		Msg("signup user successfully")

	return res, nil
}
