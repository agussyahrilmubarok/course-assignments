package sessions

import (
	"finly-go/pkg/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var store *session.Store

// AddSession initializes the session store.
func AddSession() {
	store = session.New(session.Config{
		Expiration:     24 * time.Hour, // Session expires after 24 hours
		CookieHTTPOnly: true,           // Make cookie HTTPOnly for security
		CookieSecure:   false,          // Set to true if using HTTPS
	})
}

// GetSession retrieves the current session for the request.
func GetSession(c *fiber.Ctx) (*session.Session, error) {
	sess, err := store.Get(c)
	if err != nil {
		logger.Log.Sugar().Errorf("Error getting session: %v", err)
		return nil, err
	}
	return sess, nil
}
