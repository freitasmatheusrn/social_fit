package renderer

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component, status int) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return component.Render(c.Request().Context(), c.Response())
}
func Respond(c echo.Context, component templ.Component, jsonData any, status int) error {
	accept := c.Request().Header.Get("Accept")
	isHtmx := c.Request().Header.Get("Hx-Request") == "true"
	if isHtmx || strings.Contains(accept, "text/html") {
		return Render(c, component, status)
	}

	return c.JSON(status, jsonData)
}
