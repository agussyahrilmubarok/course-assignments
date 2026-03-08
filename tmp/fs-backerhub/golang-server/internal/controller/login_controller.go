package controller

import (
	"net/http"

	"example.com.backend/internal/domain"
	"example.com.backend/internal/model"
	"example.com.backend/internal/service"
	"example.com.backend/pkg/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type loginController struct {
	baseController
	userService service.IUserService
}

func NewLoginController(userService service.IUserService) *loginController {
	return &loginController{userService: userService}
}

func (h *loginController) Index(c *gin.Context) {
	data := gin.H{"title": "Login"}

	c.HTML(http.StatusOK, "login.html", data)
}

func (h *loginController) Login(c *gin.Context) {
	ctx := c.Request.Context()
	log := logger.GetLoggerFromContext(c)
	data := gin.H{"title": "Login"}

	var input struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		log.Warn("failed to bind login input", zap.String("user_email", input.Email), zap.Error(err))
		data["form"] = input
		data["error"] = "Your email or password is wrong"
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}

	user, err := h.userService.FindByEmail(ctx, input.Email)
	if err != nil || user == nil {
		log.Error("user not found", zap.String("user_email", input.Email), zap.Error(err))
		data["form"] = input
		data["error"] = "Your email is not registered."
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}

	if !user.ComparePassword(input.Password) {
		log.Error("incorrect password attempt", zap.String("user_email", input.Email))
		data["form"] = input
		data["error"] = "Your password is wrong."
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}

	if user.Role != domain.RoleAdmin {
		log.Error("unauthorized login attempt", zap.String("user_email", input.Email))
		data["form"] = input
		data["error"] = "You do not have permission."
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}

	// Save session
	h.saveUserSession(c, model.UserDTO{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Occupation: user.Occupation,
		ImageName:  user.ImageName,
	})

	log.Info("admin login successful",
		zap.String("user_email", input.Email),
		zap.String("user_id", user.ID),
	)
	c.Redirect(http.StatusFound, "/dashboard")
}
