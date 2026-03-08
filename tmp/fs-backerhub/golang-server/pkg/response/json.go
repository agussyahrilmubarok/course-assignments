package response

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(ctx *gin.Context, status int, message string, data interface{}) {
	ctx.JSON(status, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

type ErrorResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Error(ctx *gin.Context, status int, message string, err interface{}) {
	ctx.AbortWithStatusJSON(status, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  err,
	})
}
