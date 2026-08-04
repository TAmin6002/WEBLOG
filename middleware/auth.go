package middleware

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("session")

		if err != nil || cookie.Value == "" {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		userID, err := strconv.Atoi(cookie.Value)
		
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		c.Set("user_id", userID)

		return next(c)
	}
}