package helper

import (
	"encoding/json"
	"net/http"

	"example.com/backend/internal/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func AdminShouldLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		val := session.Get("profile")
		if val == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		pJSON, ok := val.([]byte)
		if !ok {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		var profile model.UserDTO
		if err := json.Unmarshal(pJSON, &profile); err != nil || profile.ID == "" {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		c.Set("profile", profile)
		c.Next()
	}
}

func AdminHasLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		val := session.Get("profile")
		if val == nil {
			c.Next()
			return
		}

		pJSON, ok := val.([]byte)
		if !ok {
			c.Next()
			return
		}

		var profile model.UserDTO
		if err := json.Unmarshal(pJSON, &profile); err != nil || profile.ID == "" {
			c.Next()
			return
		}

		c.Redirect(http.StatusFound, "/dashboard")
		c.Abort()
	}
}
