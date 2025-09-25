package service_test

import (
	"testing"
	"time"

	"ecommerce/user-service/service"

	"github.com/stretchr/testify/assert"
)

func TestJWTService(t *testing.T) {
	secret := "my-secret"
	duration := time.Second * 1
	jwtSvc := service.NewJWTService(secret, duration)

	t.Run("generate and validate success", func(t *testing.T) {
		token, err := jwtSvc.GenerateToken("user123")
		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		userID, err := jwtSvc.ValidateToken(token)
		assert.NoError(t, err)
		assert.Equal(t, "user123", userID)
	})

	t.Run("expired token", func(t *testing.T) {
		shortJWT := service.NewJWTService(secret, time.Millisecond*100)
		token, err := shortJWT.GenerateToken("user123")
		assert.NoError(t, err)

		// tunggu sampai expired
		time.Sleep(time.Millisecond * 200)

		userID, err := shortJWT.ValidateToken(token)
		assert.Error(t, err)
		assert.Empty(t, userID)
	})

	t.Run("invalid secret", func(t *testing.T) {
		token, err := jwtSvc.GenerateToken("user123")
		assert.NoError(t, err)

		otherJWT := service.NewJWTService("wrong-secret", duration)
		userID, err := otherJWT.ValidateToken(token)

		assert.Error(t, err)
		assert.Empty(t, userID)
	})

	t.Run("malformed token", func(t *testing.T) {
		userID, err := jwtSvc.ValidateToken("not-a-valid-token")

		assert.Error(t, err)
		assert.Empty(t, userID)
	})
}
