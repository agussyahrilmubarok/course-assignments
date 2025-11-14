package server

import (
	"context"
	"fmt"
	"net/http"

	"example.com/account/internal/account"
	"example.com/account/internal/config"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type HttpServer struct {
	config *config.Config
	logger zerolog.Logger
	db     *gorm.DB
	echo   *echo.Echo
}

func NewHttpServer(
	config *config.Config,
	logger zerolog.Logger,
	db *gorm.DB,
) *HttpServer {
	return &HttpServer{
		config: config,
		logger: logger,
		db:     db,
		echo:   echo.New(),
	}
}

func (s *HttpServer) Start() error {
	s.echo.HideBanner = false
	s.echo.Use(middleware.Recover())
	s.echo.Validator = &account.CustomValidator{Validator: validator.New()}

	s.echo.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"status":  "ok",
			"service": "account-service",
		})
	})

	accountRouter := account.NewRouter(s.config, s.logger, s.db, s.echo.Group(""))
	accountRouter.MapAPIV1()

	addr := fmt.Sprintf(":%d", s.config.Http.Port)
	s.logger.Info().Msgf("server running at %s", addr)
	if err := s.echo.Start(addr); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal().Err(err).Msg("failed to start server")
		return err
	}

	return nil
}

func (s *HttpServer) Stop(ctx context.Context) error {
	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Fatal().Err(err).Msg("graceful shutdown failed")
		return err
	}

	return nil
}
