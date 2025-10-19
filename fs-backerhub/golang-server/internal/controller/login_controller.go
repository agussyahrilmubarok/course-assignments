package controller

import (
	"net/http"

	"example.com/backend/internal/domain"
	"example.com/backend/internal/model"
	"example.com/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type loginController struct {
	baseController
	userService service.IUserService
	log         zerolog.Logger
}

func NewLoginController(
	userService service.IUserService,
	log zerolog.Logger,
) *loginController {
	return &loginController{
		userService: userService,
		log:         log,
	}
}

func (h *loginController) Index(c *gin.Context) {
	data := gin.H{
		"title": "Login",
	}

	c.HTML(http.StatusOK, "login.html", data)
}

func (h *loginController) Login(c *gin.Context) {
	data := gin.H{
		"title": "Login",
	}

	var input struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&input); err != nil {
		h.log.Warn().Err(err).Msg("failed to bind login")

		data["form"] = input
		data["error"] = "Your email or password is wrong"
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}

	ctx := c.Request.Context()

	user, err := h.userService.FindByEmail(ctx, input.Email)
	if err != nil || user == nil {
		h.log.Error().Err(err).Msgf("user not found with email %v", input.Email)

		data["form"] = input
		data["error"] = "Your email is not registered."
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}

	if !user.ComparePassword(input.Password) {
		h.log.Error().Msgf("incorrect password for email %s", input.Email)

		data["form"] = input
		data["error"] = "Your password is wrong."
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}

	if user.Role != domain.RoleAdmin {
		h.log.Error().Msgf("unauthorized login attempt by user %s with role %s", input.Email, user.Role)

		data["form"] = input
		data["error"] = "You do not have permission."
		c.HTML(http.StatusBadRequest, "login.html", data)
		return
	}

	h.saveUserSession(c, model.UserDTO{
		ID:         user.ID,
		Name:       user.Name,
		Email:      user.Email,
		Occupation: user.Occupation,
		ImageName:  user.ImageName,
	})

	c.Redirect(http.StatusFound, "/dashboard")
}
