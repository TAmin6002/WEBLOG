package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"weblog/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) HandleLogin(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "login.html", nil)
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.authService.Login(username, password)
	if err != nil {
		return c.Render(http.StatusOK, "login.html", map[string]interface{}{
			"Error":    err.Error(),
			"Username": username,
		})
	}

	setSessionCookie(c, user.ID)

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthHandler) HandleSignup(c echo.Context) error {
	if c.Request().Method == http.MethodGet {
		return c.Render(http.StatusOK, "signup.html", nil)
	}

	username := c.FormValue("username")
	password := c.FormValue("password")

	if username == "" || password == "" {
		return c.Render(http.StatusOK, "signup.html", map[string]interface{}{
			"Error":    "username and password cannot be empty",
			"Username": username,
		})
	}

	user, err := h.authService.Signup(username, password)
	if err != nil {
		return c.Render(http.StatusOK, "signup.html", map[string]interface{}{
			"Error":    err.Error(),
			"Username": username,
		})
	}

	setSessionCookie(c, user.ID)

	return c.Redirect(http.StatusSeeOther, "/")
}

func (h *AuthHandler) HandleLogout(c echo.Context) error {
	c.SetCookie(&http.Cookie{
		Name:   "session",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
	return c.Redirect(http.StatusSeeOther, "/login")
}

func setSessionCookie(c echo.Context, userID int) {
	c.SetCookie(&http.Cookie{
		Name:     "session",
		Value:    strconv.Itoa(userID),
		Path:     "/",
		HttpOnly: true,
	})
}