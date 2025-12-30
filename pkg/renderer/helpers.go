package renderer

import (
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/freitasmatheusrn/social-fit/pkg/rest"
	"github.com/freitasmatheusrn/social-fit/pkg/toast"
	"github.com/labstack/echo/v4"
)

func Render(c echo.Context, component templ.Component, status int) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTMLCharsetUTF8)
	return component.Render(c.Request().Context(), c.Response())
}
func Respond(c echo.Context, component templ.Component, jsonData any, errObject *rest.ApiErr) error {
	if errObject.Code < 400 {
		t := toast.Danger(errObject.Message)
		t.SetHXTriggerHeader(c)
		return errObject
	}
	accept := c.Request().Header.Get("Accept")
	isHtmx := c.Request().Header.Get("Hx-Request") == "true"
	if isHtmx || strings.Contains(accept, "text/html") {
		return Render(c, component, http.StatusOK)
	}

	return c.JSON(errObject.Code, jsonData)
}
