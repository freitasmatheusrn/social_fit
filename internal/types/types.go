package types

import (
	"github.com/labstack/echo/v4"
)

type CurrentUser struct {
	ID    string
	Name  string
	Email string
}

func SetCurrentUser(c echo.Context, user CurrentUser) {
	c.Set("current_user", user)
}

func GetCurrentUser(c echo.Context) CurrentUser {
	return c.Get("current_user").(CurrentUser)
}