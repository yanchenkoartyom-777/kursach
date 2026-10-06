package middleware

import (
	"car-dealership-system/pkg/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	UserIDKey    = "user_id"
	UserEmailKey = "user_email"
	UserRoleKey  = "user_role"
)

// AuthMiddleware creates a middleware that requires authentication
func AuthMiddleware(sessionManager *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(session.CookieName)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		sessionData, err := sessionManager.ValidateSession(cookie)
		if err != nil {
			// Clear invalid cookie
			c.SetCookie(session.CookieName, "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		// Store user data in context
		c.Set(UserIDKey, sessionData.UserID)
		c.Set(UserEmailKey, sessionData.Email)
		c.Set(UserRoleKey, sessionData.Role)

		c.Next()
	}
}

// RequireAdmin creates a middleware that requires admin role
func RequireAdmin(sessionManager *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(session.CookieName)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		sessionData, err := sessionManager.ValidateSession(cookie)
		if err != nil {
			c.SetCookie(session.CookieName, "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		if sessionData.Role != "admin" {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"error":   "Доступ запрещен",
				"message": "У вас нет прав для доступа к этой странице",
			})
			c.Abort()
			return
		}

		// Store user data in context
		c.Set(UserIDKey, sessionData.UserID)
		c.Set(UserEmailKey, sessionData.Email)
		c.Set(UserRoleKey, sessionData.Role)

		c.Next()
	}
}

// RequireClient creates a middleware that requires client role
func RequireClient(sessionManager *session.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		cookie, err := c.Cookie(session.CookieName)
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		sessionData, err := sessionManager.ValidateSession(cookie)
		if err != nil {
			c.SetCookie(session.CookieName, "", -1, "/", "", false, true)
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		if sessionData.Role != "client" && sessionData.Role != "admin" {
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"error":   "Доступ запрещен",
				"message": "У вас нет прав для доступа к этой странице",
			})
			c.Abort()
			return
		}

		// Store user data in context
		c.Set(UserIDKey, sessionData.UserID)
		c.Set(UserEmailKey, sessionData.Email)
		c.Set(UserRoleKey, sessionData.Role)

		c.Next()
	}
}

// GetUserID retrieves user ID from context
func GetUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, false
	}
	id, ok := userID.(uint)
	return id, ok
}

// GetUserEmail retrieves user email from context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(UserEmailKey)
	if !exists {
		return "", false
	}
	emailStr, ok := email.(string)
	return emailStr, ok
}

// GetUserRole retrieves user role from context
func GetUserRole(c *gin.Context) (string, bool) {
	role, exists := c.Get(UserRoleKey)
	if !exists {
		return "", false
	}
	roleStr, ok := role.(string)
	return roleStr, ok
}
