package service

import (
	"context"
	"errors"
	"time"

	"example.com.backend/internal/config"
	"example.com.backend/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type IJwtService interface {
	Generate(ctx context.Context, userID string) (string, error)
	Validate(ctx context.Context, tokenString string) (string, error)
}

type jwtService struct {
	cfg *config.Config
}

func NewJwtService(cfg *config.Config) IJwtService {
	return &jwtService{cfg: cfg}
}

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func (s *jwtService) Generate(ctx context.Context, userID string) (string, error) {
	log := logger.GetLoggerFromContext(ctx)

	claims := &jwtCustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.cfg.Backend.Token.ExpiresAt)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    s.cfg.Backend.Name,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.cfg.Backend.Token.SecretKey))
	if err != nil {
		log.Error("failed to sign jwt",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return "", err
	}

	log.Info("jwt token generated successfully",
		zap.String("user_id", userID),
	)
	return signedToken, nil
}

func (s *jwtService) Validate(ctx context.Context, tokenString string) (string, error) {
	log := logger.GetLoggerFromContext(ctx)

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwtCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(s.cfg.Backend.Token.SecretKey), nil
		},
	)
	if err != nil {
		log.Warn("invalid jwt token",
			zap.Error(err),
		)
		return "", err
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok || !token.Valid {
		log.Warn("invalid token claims",
			zap.String("token", tokenString),
		)
		return "", errors.New("invalid token claims")
	}

	log.Info("jwt token validated successfully",
		zap.String("user_id", claims.UserID),
	)
	return claims.UserID, nil
}
