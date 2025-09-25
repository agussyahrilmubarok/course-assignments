package sessions

import (
	"encoding/json"
	"finly-go/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func SetInfo(c *fiber.Ctx, message string, msgType string) error {
	sess, err := GetSession(c)
	if err != nil {
		logger.Log.Error("Error retrieving session", zap.Error(err))
		return err
	}

	info := fiber.Map{
		"Message": message,
		"Type":    msgType,
		"Time":    time.Now().Unix(), // Store timestamp when info is set
	}

	infoJSON, err := json.Marshal(info)
	if err != nil {
		logger.Log.Error("Error marshalling info map", zap.Error(err))
		return err
	}

	sess.Set("info", string(infoJSON))
	if err := sess.Save(); err != nil {
		logger.Log.Error("Error saving session", zap.Error(err))
		return err
	}

	return nil
}

func GetInfo(c *fiber.Ctx) (fiber.Map, error) {
	sess, err := GetSession(c)
	if err != nil {
		logger.Log.Error("Error retrieving session", zap.Error(err))
		return nil, err
	}

	info := sess.Get("info")
	if info == nil {
		logger.Log.Info("No info found in session")
		return nil, nil
	}

	infoJSON, ok := info.(string)
	if !ok {
		logger.Log.Error("Error: info in session is not a valid string")
		return nil, nil
	}

	var infoMap fiber.Map
	err = json.Unmarshal([]byte(infoJSON), &infoMap)
	if err != nil {
		logger.Log.Error("Error unmarshalling info", zap.Error(err))
		return nil, err
	}

	setTime := int64(infoMap["Time"].(float64))
	if time.Now().Unix()-setTime > 3 { // Info expires after 3 seconds
		sess.Delete("info") // Remove the info from the session
		if err := sess.Save(); err != nil {
			logger.Log.Error("Error saving session after deleting info", zap.Error(err))
			return nil, err
		}
		return nil, nil
	}

	return infoMap, nil
}

func DestroyInfoSession(c *fiber.Ctx) error {
	sess, err := GetSession(c)
	if err != nil {
		logger.Log.Error("Error retrieving session", zap.Error(err))
		return err
	}

	sess.Delete("info")
	if err := sess.Save(); err != nil {
		logger.Log.Error("Error saving session after deleting info", zap.Error(err))
		return err
	}

	return nil
}
