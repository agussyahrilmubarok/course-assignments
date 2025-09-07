package usecase_test

import (
	"context"
	"ecommerce/user-service/model"
	"ecommerce/user-service/service/mocks"
	"ecommerce/user-service/store"
	storeMock "ecommerce/user-service/store/mocks"
	"ecommerce/user-service/usecase"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthUseCase_SignUp(t *testing.T) {
	ctx := context.TODO()
	mockStore := new(storeMock.IUserMongoStore)
	mockJWT := new(mocks.IJWTService)
	log := zerolog.Nop()

	uc := usecase.NewAuthUseCase(mockStore, mockJWT, log)

	t.Run("success", func(t *testing.T) {
		req := model.SignUpRequest{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "secret123",
		}

		mockStore.On("FindByEmail", ctx, "john@example.com").
			Return(nil, mongo.ErrNoDocuments).Once()

		userID := primitive.NewObjectID()
		mockStore.On("Create", ctx, mock.AnythingOfType("*store.User")).
			Return(&store.User{
				ID:    userID,
				Name:  "John Doe",
				Email: "john@example.com",
			}, nil).Once()

		res, err := uc.SignUp(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, userID.Hex(), res.ID)
		assert.Equal(t, "john@example.com", res.Email)

		mockStore.AssertExpectations(t)
	})

	t.Run("email already exists", func(t *testing.T) {
		req := model.SignUpRequest{
			Name:     "Jane Doe",
			Email:    "jane@example.com",
			Password: "secret123",
		}

		mockStore.On("FindByEmail", ctx, "jane@example.com").
			Return(&store.User{}, nil).Once()

		res, err := uc.SignUp(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, usecase.ErrEmailAlreadyExists, err)

		mockStore.AssertExpectations(t)
	})

	t.Run("create user error", func(t *testing.T) {
		req := model.SignUpRequest{
			Name:     "Rick Doe",
			Email:    "rick@example.com",
			Password: "secret123",
		}

		mockStore.On("FindByEmail", ctx, "rick@example.com").
			Return(nil, mongo.ErrNoDocuments).Once()

		mockStore.On("Create", ctx, mock.AnythingOfType("*store.User")).
			Return(nil, errors.New("insert failed")).Once()

		res, err := uc.SignUp(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, res)

		mockStore.AssertExpectations(t)
	})
}

func TestAuthUseCase_SignIn(t *testing.T) {
	ctx := context.TODO()
	mockStore := new(storeMock.IUserMongoStore)
	mockJWT := new(mocks.IJWTService)
	log := zerolog.Nop()

	uc := usecase.NewAuthUseCase(mockStore, mockJWT, log)

	t.Run("success", func(t *testing.T) {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("mypassword"), bcrypt.DefaultCost)
		user := &store.User{
			ID:       primitive.NewObjectID(),
			Email:    "user@example.com",
			Password: string(hashed),
		}

		mockStore.On("FindByEmail", ctx, "user@example.com").
			Return(user, nil).Once()

		mockJWT.On("GenerateToken", user.ID.Hex()).
			Return("token123", nil).Once()

		res, err := uc.SignIn(ctx, model.SignInRequest{
			Email:    "user@example.com",
			Password: "mypassword",
		})

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, "token123", res.Token)

		mockStore.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
		user := &store.User{
			ID:       primitive.NewObjectID(),
			Email:    "wrong@example.com",
			Password: string(hashed),
		}

		mockStore.On("FindByEmail", ctx, "wrong@example.com").
			Return(user, nil).Once()

		res, err := uc.SignIn(ctx, model.SignInRequest{
			Email:    "wrong@example.com",
			Password: "wrongpass",
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, usecase.ErrInvalidCredentials, err)

		mockStore.AssertExpectations(t)
	})

	t.Run("email not found", func(t *testing.T) {
		mockStore.On("FindByEmail", ctx, "missing@example.com").
			Return(nil, mongo.ErrNoDocuments).Once()

		res, err := uc.SignIn(ctx, model.SignInRequest{
			Email:    "missing@example.com",
			Password: "whatever",
		})

		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Equal(t, usecase.ErrInvalidCredentials, err)

		mockStore.AssertExpectations(t)
	})

	t.Run("jwt generation error", func(t *testing.T) {
		hashed, _ := bcrypt.GenerateFromPassword([]byte("secret123"), bcrypt.DefaultCost)
		user := &store.User{
			ID:       primitive.NewObjectID(),
			Email:    "jwtfail@example.com",
			Password: string(hashed),
		}

		mockStore.On("FindByEmail", ctx, "jwtfail@example.com").
			Return(user, nil).Once()

		mockJWT.On("GenerateToken", user.ID.Hex()).
			Return("", errors.New("jwt error")).Once()

		res, err := uc.SignIn(ctx, model.SignInRequest{
			Email:    "jwtfail@example.com",
			Password: "secret123",
		})

		assert.Error(t, err)
		assert.Nil(t, res)

		mockStore.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
}
