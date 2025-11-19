package account

import (
	"example.com/account/internal/config"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	_ "example.com/account/cmd/server/docs"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Router struct {
	config    *config.Config
	logger    zerolog.Logger
	db        *gorm.DB
	echoGroup *echo.Group
}

func NewRouter(
	config *config.Config,
	logger zerolog.Logger,
	db *gorm.DB,
	echoGroup *echo.Group,
) *Router {
	return &Router{
		config:    config,
		logger:    logger,
		db:        db,
		echoGroup: echoGroup,
	}
}

func (r *Router) MapAPIV1() {
	store := NewStore(r.db, r.logger)
	service := NewService(r.config, r.logger)
	handler := NewHandler(store, service, r.logger)
	customMiddleware := NewCustomMiddleware(service, r.logger)

	v1 := r.echoGroup.Group("/api/v1/accounts")
	{
		v1.GET("/swagger/*", echoSwagger.WrapHandler)

		v1.POST("/sign-up", handler.SignUp)
		v1.POST("/sign-in", handler.SignIn, middleware.RateLimiterWithConfig(customMiddleware.RateLimiterConfig()))
		v1.POST("/validate", handler.Validate)
		v1.GET("/me", handler.GetMe, customMiddleware.XUserID())
	}
}
