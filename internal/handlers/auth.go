package handlers

import (
	"car-dealership-system/internal/models"
	"car-dealership-system/internal/services"
	"car-dealership-system/pkg/session"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService    *services.AuthService
	sessionManager *session.Manager
}

func NewAuthHandler(authService *services.AuthService, sessionManager *session.Manager) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		sessionManager: sessionManager,
	}
}

func (h *AuthHandler) ShowLogin(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "auth/login.html", gin.H{
		"Title": "Вход",
		"User":  user,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "auth/login.html", gin.H{
			"Title": "Вход",
			"Error": "Неверные данные",
		})
		return
	}

	resp, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "auth/login.html", gin.H{
			"Title": "Вход",
			"Error": err.Error(),
		})
		return
	}

	// If 2FA is disabled, create session immediately
	if resp.Skip2FA && resp.User != nil {
		sessionValue, err := h.sessionManager.CreateSession(resp.User.ID, resp.User.Email, resp.User.FirstName, resp.User.Role)
		if err != nil {
			c.HTML(http.StatusInternalServerError, "auth/login.html", gin.H{
				"Title": "Вход",
				"Error": "Ошибка создания сессии",
			})
			return
		}

		c.SetCookie(
			session.CookieName,
			sessionValue,
			h.sessionManager.GetMaxAge(),
			"/",
			"",
			false,
			true,
		)

		if resp.User.Role == "admin" {
			c.Redirect(http.StatusFound, "/admin/dashboard")
		} else {
			c.Redirect(http.StatusFound, "/cars")
		}
		return
	}

	// Show 2FA verification page
	c.HTML(http.StatusOK, "auth/verify.html", gin.H{
		"Title":   "Подтверждение",
		"Email":   resp.Email,
		"Success": resp.Message,
	})
}

func (h *AuthHandler) ShowVerify(c *gin.Context) {
	email := c.Query("email")
	c.HTML(http.StatusOK, "auth/verify.html", gin.H{
		"Title": "Подтверждение",
		"Email": email,
	})
}

func (h *AuthHandler) Verify(c *gin.Context) {
	var req models.VerifyRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "auth/verify.html", gin.H{
			"Title": "Подтверждение",
			"Email": req.Email,
			"Error": "Неверные данные",
		})
		return
	}

	user, err := h.authService.VerifyTwoFactor(req.Email, req.Code)
	if err != nil {
		c.HTML(http.StatusUnauthorized, "auth/verify.html", gin.H{
			"Title": "Подтверждение",
			"Email": req.Email,
			"Error": err.Error(),
		})
		return
	}

	sessionValue, err := h.sessionManager.CreateSession(user.ID, user.Email, user.FirstName, user.Role)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "auth/verify.html", gin.H{
			"Title": "Подтверждение",
			"Email": req.Email,
			"Error": "Ошибка создания сессии",
		})
		return
	}

	c.SetCookie(
		session.CookieName,
		sessionValue,
		h.sessionManager.GetMaxAge(),
		"/",
		"",
		false,
		true,
	)

	if user.Role == "admin" {
		c.Redirect(http.StatusFound, "/admin/dashboard")
	} else {
		c.Redirect(http.StatusFound, "/cars")
	}
}

func (h *AuthHandler) ShowRegister(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/register.html", gin.H{
		"Title": "Регистрация",
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "auth/register.html", gin.H{
			"Title": "Регистрация",
			"Error": "Неверные данные",
		})
		return
	}

	if err := h.authService.Register(&req); err != nil {
		c.HTML(http.StatusBadRequest, "auth/register.html", gin.H{
			"Title": "Регистрация",
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "auth/login.html", gin.H{
		"Title":   "Вход",
		"Success": "Регистрация успешна! Войдите в систему",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie(session.CookieName, "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}

func (h *AuthHandler) ShowProfile(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "auth/profile.html", gin.H{
		"Title": "Профиль",
		"User":  user,
	})
}

func (h *AuthHandler) ShowForgotPassword(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/forgot-password.html", gin.H{
		"Title": "Забыли пароль?",
	})
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "auth/forgot-password.html", gin.H{
			"Title": "Забыли пароль?",
			"Error": "Неверные данные",
		})
		return
	}

	if err := h.authService.RequestPasswordReset(req.Email); err != nil {
		c.HTML(http.StatusInternalServerError, "auth/forgot-password.html", gin.H{
			"Title": "Забыли пароль?",
			"Error": "Ошибка отправки письма",
		})
		return
	}

	c.HTML(http.StatusOK, "auth/forgot-password.html", gin.H{
		"Title":   "Забыли пароль?",
		"Success": "Если email существует в системе, на него будет отправлено письмо с инструкциями",
	})
}

func (h *AuthHandler) ShowResetPassword(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"Error": "Неверная ссылка для сброса пароля",
		})
		return
	}

	c.HTML(http.StatusOK, "auth/reset-password.html", gin.H{
		"Title": "Сброс пароля",
		"Token": token,
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBind(&req); err != nil {
		c.HTML(http.StatusBadRequest, "auth/reset-password.html", gin.H{
			"Title": "Сброс пароля",
			"Token": req.Token,
			"Error": "Неверные данные",
		})
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.Password); err != nil {
		c.HTML(http.StatusBadRequest, "auth/reset-password.html", gin.H{
			"Title": "Сброс пароля",
			"Token": req.Token,
			"Error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "auth/login.html", gin.H{
		"Title":   "Вход",
		"Success": "Пароль успешно изменён! Войдите с новым паролем",
	})
}
