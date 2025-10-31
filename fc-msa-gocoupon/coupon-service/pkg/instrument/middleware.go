package instrument

import (
	"context"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type ctxKeyRequestID struct{}
type ctxKeyTraceID struct{}

var (
	RequestIDKey = ctxKeyRequestID{}
	TraceIDKey   = ctxKeyTraceID{}
)

// GetLogger membuat logger dari context
func GetLogger(ctx context.Context, baseLogger zerolog.Logger) zerolog.Logger {
	logger := baseLogger
	if rid, ok := ctx.Value(RequestIDKey).(string); ok {
		logger = logger.With().Str("request_id", rid).Logger()
	}
	if tid, ok := ctx.Value(TraceIDKey).(string); ok {
		logger = logger.With().Str("trace_id", tid).Logger()
	}
	return logger
}

func Middleware(tracer trace.Tracer, baseLogger zerolog.Logger) gin.HandlerFunc {

	return func(c *gin.Context) {
		start := time.Now()
		ctx := c.Request.Context()

		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = xid.New().String()
		}
		c.Writer.Header().Set("X-Request-ID", requestID)
		ctx = context.WithValue(ctx, RequestIDKey, requestID)

		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(c.Request.Header))
		ctx, span := tracer.Start(ctx, c.Request.Method+" "+c.FullPath())
		defer span.End()
		traceID := span.SpanContext().TraceID().String()
		ctx = context.WithValue(ctx, TraceIDKey, traceID)

		c.Request = c.Request.WithContext(ctx)

		defer func() {
			log := GetLogger(c.Request.Context(), baseLogger)
			duration := time.Since(start)
			status := c.Writer.Status()

			if r := recover(); r != nil {
				log.Error().
					Interface("panic", r).
					Str("stack", string(debug.Stack())).
					Msg("Recovered from panic")
				c.AbortWithStatus(http.StatusInternalServerError)
				panic(r)
			}

			span.SetAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.path", c.FullPath()),
				attribute.Int("http.status_code", status),
				attribute.String("client.ip", c.ClientIP()),
				attribute.Float64("http.duration_ms", float64(duration.Milliseconds())),
			)

			if len(c.Errors) > 0 {
				for _, e := range c.Errors {
					span.RecordError(e.Err)
				}
			}

			log.Info().
				Str("method", c.Request.Method).
				Str("path", c.FullPath()).
				Int("status_code", status).
				Dur("duration", duration).
				Msg("handled HTTP request")
		}()

		c.Next()
	}
}
