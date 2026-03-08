package sessions

import (
	"finly-go/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func SetProfile(c *fiber.Ctx, userID string) error {
	sess, err := GetSession(c)
	if err != nil {
		logger.Log.Error("Error retrieving session", zap.Error(err))
		return err
	}

	sess.Set("profile", userID) // Store the user ID in session
	if err := sess.Save(); err != nil {
		logger.Log.Error("Error saving session", zap.Error(err))
		return err
	}
	return nil
}

func GetProfile(c *fiber.Ctx) (string, error) {
	sess, err := GetSession(c)
	if err != nil {
		logger.Log.Error("Error retrieving session", zap.Error(err))
		return "", err
	}

	profileID := sess.Get("profile")
	if profileID == nil {
		logger.Log.Info("Profile not found in session")
		return "", nil
	}

	return profileID.(string), nil
}

func DestroyAuthSession(c *fiber.Ctx) error {
	sess, err := GetSession(c)
	if err != nil {
		logger.Log.Error("Error retrieving session", zap.Error(err))
		return err
	}
	sess.Delete("profile") // Remove the user profile from the session
	if err := sess.Save(); err != nil {
		logger.Log.Error("Error saving session after deleting profile", zap.Error(err))
		return err
	}
	return nil
}
