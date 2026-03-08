package controller

import (
	"finly-go/internal/forms"
	"finly-go/internal/model"
	"finly-go/internal/sessions"
	"finly-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type AuthController struct {
	userCollection *mongo.Collection
}

func NewAuthController(
	userCollection *mongo.Collection,
) *AuthController {
	return &AuthController{
		userCollection: userCollection,
	}
}

func (ah *AuthController) GetSignUp(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Sign Up",
	}

	info, err := sessions.GetInfo(c)
	if err == nil && info != nil {
		data["Info"] = info
	}

	return c.Render("signup", data, "layouts/main")
}

func (ah *AuthController) PostSignUp(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Sign Up",
	}

	var form forms.SignupForm
	if err := c.BodyParser(&form); err != nil {
		data["Error"] = "Error parsing form data"
		data["User"] = form

		logger.Log.Error("Error parsing form data", zap.Error(err))
		return c.Render("signup", data, "layouts/main")
	}

	if err := form.Validate(); err != nil {
		data["Error"] = err.Error()
		data["User"] = form

		logger.Log.Error("Form validation failed", zap.String("error", err.Error()))
		return c.Render("signup", data, "layouts/main")
	}

	var existingUser model.User
	err := ah.userCollection.FindOne(c.Context(), bson.M{"email": form.Email}).Decode(&existingUser)
	if err == nil {
		data["Error"] = "Email is already registered"
		data["User"] = form

		logger.Log.Warn("Email already registered", zap.String("email", form.Email))
		return c.Render("signup", data, "layouts/main")
	} else if err != mongo.ErrNoDocuments {
		data["Error"] = "Error checking email"
		data["User"] = form

		logger.Log.Error("Error checking email", zap.Error(err))
		return c.Render("signup", data, "layouts/main")
	}

	hashedPassword, err := form.HashPassword()
	if err != nil {
		data["Error"] = "Error hashing password"
		data["User"] = form

		logger.Log.Error("Error hashing password", zap.Error(err))
		return c.Render("signup", data, "layouts/main")
	}

	user := model.User{
		Email:    form.Email,
		Password: hashedPassword,
	}

	_, err = ah.userCollection.InsertOne(c.Context(), user)
	if err != nil {
		data["Error"] = "Error saving user"
		data["User"] = form

		logger.Log.Error("Error saving user", zap.Error(err))
		return c.Render("signup", data, "layouts/main")
	}

	if err := sessions.SetInfo(c, "Registration successful! You can now log in.", "success"); err != nil {
		logger.Log.Error("Error setting session info", zap.Error(err))
	}

	return c.Redirect("/login")
}

func (ah *AuthController) GetLogin(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Login",
	}

	info, err := sessions.GetInfo(c)
	if err == nil && info != nil {
		data["Info"] = info
	}

	return c.Render("login", data, "layouts/main")
}

func (ah *AuthController) PostLogin(c *fiber.Ctx) error {
	data := fiber.Map{
		"Title": "Login",
	}

	var form forms.LoginForm
	if err := c.BodyParser(&form); err != nil {
		data["Error"] = "Error parsing form data"
		data["User"] = form

		logger.Log.Error("Error parsing form data", zap.Error(err))
		return c.Render("login", data, "layouts/main")
	}

	if err := form.Validate(); err != nil {
		data["Error"] = err.Error()
		data["User"] = form

		logger.Log.Error("Form validation failed", zap.String("error", err.Error()))
		return c.Render("login", data, "layouts/main")
	}

	var user model.User
	err := ah.userCollection.FindOne(c.Context(), bson.M{"email": form.Email}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			data["Error"] = "Invalid email or password"
			data["User"] = form
		} else {
			data["Error"] = "Error finding user"
			data["User"] = form
			logger.Log.Error("Error finding user", zap.Error(err))
		}
		return c.Render("login", data, "layouts/main")
	}

	if !form.ComparePassword(user.Password) {
		data["Error"] = "Invalid email or password"
		data["User"] = form
		return c.Render("login", data, "layouts/main")
	}

	if err := sessions.SetProfile(c, user.ID); err != nil {
		data["Error"] = "Error saving session"
		data["User"] = form

		logger.Log.Error("Error setting session", zap.Error(err))
		return c.Render("login", data, "layouts/main")
	}

	if err := sessions.SetInfo(c, "Login successful!", "success"); err != nil {
		logger.Log.Error("Error setting session info", zap.Error(err))
	}

	return c.Redirect("/dashboard")
}

func (ah *AuthController) GetLogout(c *fiber.Ctx) error {
	if err := sessions.DestroyAuthSession(c); err != nil {
		logger.Log.Error("Error destroying auth session", zap.Error(err))
		return err
	}

	if err := sessions.DestroyInfoSession(c); err != nil {
		logger.Log.Error("Error destroying info session", zap.Error(err))
		return err
	}

	if err := sessions.SetInfo(c, "You have logged out successfully.", "success"); err != nil {
		logger.Log.Error("Error setting session info", zap.Error(err))
	}

	return c.Redirect("/login")
}
