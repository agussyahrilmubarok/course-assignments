package middleware

import (
	"example.com.backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(RequestIDHeader)
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)

		ctx := c.Request.Context()
		ctx = logger.WithRequestID(ctx, requestID)
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set(RequestIDHeader, requestID)

		c.Next()
	}
}
