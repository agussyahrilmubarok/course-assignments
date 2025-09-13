package middleware

import "github.com/gin-gonic/gin"

func GetUserID(c *gin.Context) (uint, bool) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return 0, false
	}
	userID, ok := userIDVal.(uint)
	return userID, ok
}

func GetUserRole(c *gin.Context) (string, bool) {
	roleVal, exists := c.Get("userRole")
	if !exists {
		return "", false
	}
	role, ok := roleVal.(string)
	return role, ok
}
