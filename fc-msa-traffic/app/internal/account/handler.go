package account

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	store   IStore
	service IService
	log     zerolog.Logger
}

func NewHandler(
	store IStore,
	service IService,
	log zerolog.Logger,
) *Handler {
	return &Handler{
		store:   store,
		service: service,
		log:     log,
	}
}

// SignUp godoc
// @Summary Register a new user
// @Description Register a new account with email and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param SignUpRequest body SignUpRequest true "Sign up payload"
// @Success 200 {object} UserResponse "Registered user data"
// @Failure 400 {object} map[string]interface{} "Validation or bad request error"
// @Failure 500 {object} map[string]interface{} "Server error"
// @Router /sign-up [post]
func (h *Handler) SignUp(c echo.Context) error {
	ctx := c.Request().Context()

	var req SignUpRequest
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Failed to bind request")
		return c.JSON(400, echo.Map{"error": "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		h.log.Warn().Err(err).Msg("Validation failed")
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation error", "errors": errors})
	}

	if h.store.ExistsUserByEmailIgnoreCase(ctx, req.Email) {
		h.log.Warn().Msg("Failed to sign up, email already in use")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email already in use"})
	}

	var user *User
	user = req.ToUser()
	if err := h.store.SaveUser(ctx, user); err != nil {
		h.log.Error().Err(err).Msg("Failed to save user")
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to sign up"})
	}

	var res UserResponse
	res.FromUser(user)

	return c.JSON(http.StatusOK, res)
}

// SignIn godoc
// @Summary User login
// @Description Authenticate user and return JWT token
// @Tags Authentication
// @Accept json
// @Produce json
// @Param SignInRequest body SignInRequest true "Sign in payload"
// @Success 200 {object} AccountResponse "JWT token and user data"
// @Failure 400 {object} map[string]interface{} "Validation or authentication error"
// @Router /sign-in [post]
func (h *Handler) SignIn(c echo.Context) error {
	ctx := c.Request().Context()

	var req SignInRequest
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Failed to bind request")
		return c.JSON(400, echo.Map{"error": "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		h.log.Warn().Err(err).Msg("Validation failed")
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, fieldErr := range validationErrors {
			errors[fieldErr.Field()] = fieldErr.Tag()
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation error", "errors": errors})
	}

	user, err := h.store.FindUserByEmail(ctx, req.Email)
	if err != nil || user == nil {
		h.log.Warn().Err(err).Msg("Failed to find email when sign in")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Email is not registered"})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		h.log.Warn().Err(err).Msg("Failed to compare password")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Wrong password"})
	}

	tokenString, err := h.service.GenerateJwt(user.ID)
	if err != nil {
		h.log.Warn().Err(err).Msg("Failed to generate jwt")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to sign in"})
	}

	var userRes UserResponse
	userRes.FromUser(user)

	return c.JSON(http.StatusOK, AccountResponse{
		Token: tokenString,
		User:  userRes,
	})
}

// Validate godoc
// @Summary Validate JWT token
// @Description Validate a JWT token and return user info if valid
// @Tags Authentication
// @Accept json
// @Produce json
// @Param ValidateRequest body ValidateRequest true "Token to validate"
// @Success 200 {object} AccountResponse "Validated user info and token"
// @Failure 400 {object} map[string]interface{} "Invalid token or validation failed"
// @Router /validate [post]
func (h *Handler) Validate(c echo.Context) error {
	ctx := c.Request().Context()

	var req ValidateRequest
	if err := c.Bind(&req); err != nil {
		h.log.Warn().Err(err).Msg("Failed to bind request")
		return c.JSON(400, echo.Map{"error": "Invalid request format"})
	}

	userID, err := h.service.ValidateJwt(req.Token)
	if err != nil {
		h.log.Warn().Err(err).Msg("Failed to validate token")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Failed to validate account"})
	}

	user, err := h.store.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		h.log.Warn().Err(err).Msg("Failed to find user by id")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Account is not registered"})
	}

	var userRes UserResponse
	userRes.FromUser(user)

	return c.JSON(http.StatusOK, AccountResponse{
		Token: req.Token,
		User:  userRes,
	})
}

// GetMe godoc
// @Summary Get current authenticated user info
// @Description Retrieve info about the authenticated user from JWT token
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} UserResponse "Current user info"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 400 {object} map[string]interface{} "User not found"
// @Router /me [get]
func (h *Handler) GetMe(c echo.Context) error {
	ctx := c.Request().Context()

	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		h.log.Warn().Msg("user_id not found in context")
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized"})
	}

	user, err := h.store.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		h.log.Warn().Err(err).Msg("Failed to find user by id")
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Account is not registered"})
	}

	var userRes UserResponse
	userRes.FromUser(user)

	return c.JSON(http.StatusOK, userRes)
}
