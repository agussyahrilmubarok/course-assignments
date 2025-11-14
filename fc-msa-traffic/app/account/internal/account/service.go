package account

import (
	"time"

	"example.com/account/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=IService
type IService interface {
	GenerateJwt(userID string) (string, error)
	ValidateJwt(tokenString string) (string, error)
}

type service struct {
	config *config.Config
	logger zerolog.Logger
}

func NewService(config *config.Config, logger zerolog.Logger) IService {
	return &service{
		config: config,
		logger: logger,
	}
}

func (s *service) GenerateJwt(userID string) (string, error) {
	secretKey := s.config.Jwt.SecretKey
	expiresIn, err := time.ParseDuration(s.config.Jwt.ExpiresIn)
	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(expiresIn)).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		s.logger.Error().Err(err).Str("user_id", userID).Msg("failed to sign jwt")
		return "", err
	}

	s.logger.Info().Str("user_id", userID).Msg("generate jwt for user successfully")
	return signedToken, nil
}

func (s *service) ValidateJwt(tokenString string) (string, error) {
	secretKey := s.config.Jwt.SecretKey

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.logger.Error().Msg("unexpected signing method")
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to parse JWT")
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		s.logger.Error().Err(err).Msg("invalid JWT claims or token")
		return "", jwt.ErrTokenInvalidClaims
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		s.logger.Error().Msg("user_id claim is missing or invalid")
		return "", jwt.ErrInvalidKeyType
	}

	s.logger.Info().Str("user_id", userID).Msg("validate jwt user successfully")
	return userID, nil
}
