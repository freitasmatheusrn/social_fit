package renderer

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component) error {
	return component.Render(c.Request().Context(), c.Response())
}

func Respond(c echo.Context, component templ.Component, jsonData any, status int) error {
	accept := c.Request().Header.Get("Accept")
	if strings.Contains(accept, "text/html") {
		return Render(c, component)
	}

	return c.JSON(status, jsonData)
}
