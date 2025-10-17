package account

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=IService
type IService interface {
	GenerateJwt(userID string) (string, error)
	ValidateJwt(tokenString string) (string, error)
}

type service struct {
	cfg *Config
	log zerolog.Logger
}

func NewService(
	cfg *Config,
	log zerolog.Logger,
) IService {
	return &service{
		cfg: cfg,
		log: log,
	}
}

func (s *service) GenerateJwt(userID string) (string, error) {
	secretKey := s.cfg.Jwt.SecretKey
	expiresIn, err := time.ParseDuration(s.cfg.Jwt.ExpiresIn)
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
		s.log.Error().Err(err).Str("user_id", userID).Msg("Failed to sign JWT token")
		return "", err
	}

	s.log.Info().Str("user_id", userID).Msg("Generate jwt for user successfully")
	return signedToken, nil
}

func (s *service) ValidateJwt(tokenString string) (string, error) {
	secretKey := s.cfg.Jwt.SecretKey

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			s.log.Error().Msg("Unexpected signing method")
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		s.log.Error().Err(err).Msg("Failed to parse JWT")
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		s.log.Error().Err(err).Msg("Invalid JWT claims or token")
		return "", jwt.ErrTokenInvalidClaims
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		s.log.Error().Msg("user_id claim is missing or invalid")
		return "", jwt.ErrInvalidKeyType
	}

	s.log.Info().Str("user_id", userID).Msg("Validate jwt user successfully")
	return userID, nil
}
