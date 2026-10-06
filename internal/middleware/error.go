package middleware

import (
	"car-dealership-system/pkg/session"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorHandler handles errors and logs them
func ErrorHandler(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			logger.WithFields(logrus.Fields{
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
				"error":  err.Error(),
			}).Error("Request error")
		}
	}
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// NoCache middleware prevents browser caching of HTML pages
func NoCache() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Next()
	}
}

// OptionalAuth middleware tries to authenticate user but doesn't require it
func OptionalAuth(sessionManager *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(session.CookieName)
		if err == nil {
			sessionData, err := sessionManager.ValidateSession(cookie)
			if err == nil {
				// Store user data in context
				c.Set(UserIDKey, sessionData.UserID)
				c.Set(UserEmailKey, sessionData.Email)
				c.Set(UserRoleKey, sessionData.Role)

				// Create a simple user object for templates
				c.Set("user", map[string]interface{}{
					"ID":        sessionData.UserID,
					"Email":     sessionData.Email,
					"FirstName": sessionData.FirstName,
					"Role":      sessionData.Role,
				})
			}
		}

		c.Next()
	}
}
