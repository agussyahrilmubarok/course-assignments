package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/api/users/internal/user"
	"github.com/agussyahrilmubarok/gox/pkg/xexception"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite
	mockService *MockIService
	handler     *user.Handler
	app         *fiber.App
}

func (ts *HandlerTestSuite) SetupTest() {
	ts.mockService = NewMockIService(ts.T())
	logger := zerolog.Nop()

	ts.handler = user.NewHandler(ts.mockService, logger)
	ts.app = fiber.New()
	ts.app.Post("/api/v1/users/auth/signup", ts.handler.SignUp)
	ts.app.Post("/api/v1/users/auth/signup", ts.handler.SignIn)
}

func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (ts *HandlerTestSuite) TestSignUp_Success() {
	param := user.SignUpParam{
		Name:     "John Doe",
		Email:    "john@test.com",
		Password: "secret",
	}
	respUser := user.UserResponse{
		Name:  param.Name,
		Email: param.Email,
	}

	ts.mockService.On("SignUp", mock.Anything, param).
		Return(respUser, nil).Once()

	body, _ := json.Marshal(param)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := ts.app.Test(req)
	ts.Equal(fiber.StatusCreated, resp.StatusCode)

	ts.mockService.AssertExpectations(ts.T())
}

func (ts *HandlerTestSuite) TestSignUp_InvalidBody() {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/auth/signup", bytes.NewBuffer([]byte(`invalid-json`)))
	req.Header.Set("Content-Type", "application/json")
	resp, _ := ts.app.Test(req)

	ts.Equal(fiber.StatusBadRequest, resp.StatusCode)
}

func (ts *HandlerTestSuite) TestSignUp_ServiceReturnsHTTPError() {
	param := user.SignUpParam{
		Name:     "John",
		Email:    "john@test.com",
		Password: "123456",
	}

	body, _ := json.Marshal(param)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	httpErr := xexception.NewHTTPBadRequest("email already exists", nil)
	ts.mockService.On("SignUp", mock.Anything, param).
		Return(user.UserResponse{}, httpErr).Once()

	resp, _ := ts.app.Test(req)
	ts.Equal(fiber.StatusBadRequest, resp.StatusCode)

	ts.mockService.AssertExpectations(ts.T())
}

func (ts *HandlerTestSuite) TestSignUp_ServiceReturnsGenericError() {
	param := user.SignUpParam{
		Name:     "Jane",
		Email:    "jane@test.com",
		Password: "123456",
	}

	body, _ := json.Marshal(param)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/users/auth/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	ts.mockService.On("SignUp", mock.Anything, param).
		Return(user.UserResponse{}, errors.New("db connection failed")).Once()

	resp, _ := ts.app.Test(req)
	ts.Equal(fiber.StatusInternalServerError, resp.StatusCode)

	ts.mockService.AssertExpectations(ts.T())
}
