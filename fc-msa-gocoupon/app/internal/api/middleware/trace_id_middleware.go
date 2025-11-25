package middleware

import (
	"example.com/coupon-service/internal/instrument"
	"example.com/coupon-service/internal/logger"
	"github.com/labstack/echo/v4"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

func TraceIDMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request().Header))

			ctx, span := instrument.StartSpan(ctx, c.Request().Method+" "+c.Path())
			defer span.End()

			traceID := span.SpanContext().TraceID().String()
			ctx = logger.WithTraceID(ctx, traceID)

			c.SetRequest(c.Request().WithContext(ctx))
			return next(c)
		}
	}
}
