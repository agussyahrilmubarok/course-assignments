package service

import (
	"errors"
	"time"

	"example.com/backend/pkg/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

//go:generate mockery --name=IJwtService
type IJwtService interface {
	Generate(userID string) (string, error)
	Validate(tokenString string) (string, error)
}

type jwtService struct {
	app config.App
	log zerolog.Logger
}

func NewJwtService(app config.App, log zerolog.Logger) IJwtService {
	return &jwtService{
		app: app,
		log: log,
	}
}

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *jwtService) Generate(userID string) (string, error) {
	claims := &jwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.app.Token.ExpiresAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.app.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.app.Token.SecretKey))
	if err != nil {
		s.log.Error().Err(err).Msg("failed to sign jwt")
		return "", err
	}

	return signedToken, nil
}

func (s *jwtService) Validate(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(s.app.Token.SecretKey), nil
		},
	)
	if err != nil {
		s.log.Warn().Err(err).Msg("invalid jwt token")
		return "", err
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		return "", errors.New("invalid token claims")
	}

	return claims.UserID, nil
}
