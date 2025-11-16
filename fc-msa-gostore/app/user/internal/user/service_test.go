package user_test

import (
	"context"
	"errors"
	"testing"

	"example.com/api/users/internal/user"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	mockIStore *MockIStore
	service    user.IService
}

func (ts *ServiceTestSuite) SetupTest() {
	ts.mockIStore = NewMockIStore(ts.T())
	logger := zerolog.Nop()

	ts.service = user.NewServie(ts.mockIStore, logger)
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (ts *ServiceTestSuite) TestSignUpSuccess() {
	param := user.SignUpParam{
		Name:     "John Doe",
		Email:    "johndoe@test.com",
		Password: "P@ssw0rd!",
	}

	ts.mockIStore.On("ExistsUserEmailByIgnoreCase", mock.Anything, param.Email).
		Return(nil).Once()

	ts.mockIStore.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.User")).
		Return(nil).Once()

	res, err := ts.service.SignUp(context.Background(), param)

	ts.NoError(err)
	ts.Equal(param.Email, res.Email)
	ts.Equal(param.Name, res.Name)
	ts.mockIStore.AssertExpectations(ts.T())
}

func (ts *ServiceTestSuite) TestSignUpFailEmailAlreadyExists() {
	param := user.SignUpParam{
		Name:     "John Doe",
		Email:    "johndoe@test.com",
		Password: "P@ssw0rd!",
	}

	ts.mockIStore.On("ExistsUserEmailByIgnoreCase", mock.Anything, param.Email).
		Return(errors.New("email already exists")).Once()

	_, err := ts.service.SignUp(context.Background(), param)

	ts.Error(err)
}

func (ts *ServiceTestSuite) TestSignUpFailCreateUser() {
	param := user.SignUpParam{
		Name:     "John Doe",
		Email:    "johndoe@test.com",
		Password: "P@ssw0rd!",
	}

	ts.mockIStore.On("ExistsUserEmailByIgnoreCase", mock.Anything, param.Email).
		Return(nil).Once()

	ts.mockIStore.On("CreateUser", mock.Anything, mock.AnythingOfType("*user.User")).
		Return(errors.New("failed to create user")).Once()

	_, err := ts.service.SignUp(context.Background(), param)

	ts.Error(err)
}
