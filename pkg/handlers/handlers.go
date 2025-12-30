package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/freitasmatheusrn/social-fit/internal/views"
	"github.com/freitasmatheusrn/social-fit/pkg/renderer"
	"github.com/freitasmatheusrn/social-fit/pkg/rest"
	"github.com/labstack/echo/v4"
)

func CustomErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	apiErr, ok := err.(*rest.ApiErr)
	if !ok {
		if he, ok := err.(*echo.HTTPError); ok {
			apiErr = &rest.ApiErr{
				Message: he.Message.(string),
				Err:     http.StatusText(he.Code),
				Code:    he.Code,
			}
		} else {
			apiErr = rest.NewInternalServerError(err.Error())
		}
	}

	isHTMX := c.Request().Header.Get("HX-Request") == "true"
	isJson := c.Request().Header.Get("Accept") == "application/json"
	switch {
	case isHTMX:
		handleHTMXError(c, apiErr)
		return
	case isJson:
		c.JSON(apiErr.Code, apiErr)
	default:
		handleHTMLRedirect(c, apiErr)
	}

}

func handleHTMXError(c echo.Context, apiErr *rest.ApiErr) {
	c.Response().Status = apiErr.Code

	var component templ.Component

	switch apiErr.Code {
	case http.StatusBadRequest:
		component = views.BadRequest(apiErr)
	case http.StatusUnauthorized:
		c.Response().Header().Set("HX-Redirect", "/sign_in")
		return
	case http.StatusForbidden:
		component = views.Forbidden(apiErr)
	case http.StatusNotFound:
		component = views.NotFound(apiErr)
	case http.StatusUnprocessableEntity:
		component = views.UnprocessableEntity(apiErr)
	default:
		component = views.InternalServerError(apiErr)
	}

	component.Render(c.Request().Context(), c.Response().Writer)
}

func handleHTMLRedirect(c echo.Context, apiErr *rest.ApiErr) {
	switch apiErr.Code {
	case http.StatusUnauthorized:
		_ = c.Redirect(http.StatusSeeOther, "/login")

	case http.StatusForbidden:
		renderer.Render(c, views.Forbidden(apiErr), 403)

	case http.StatusNotFound:
		_ = renderer.Render(c, views.NotFound(apiErr), 404)

	case http.StatusUnprocessableEntity:
		renderer.Render(c, views.UnprocessableEntity(apiErr), 422)

	case http.StatusBadRequest:
		_ = renderer.Render(c, views.BadRequest(apiErr), 400)

	default:
		_ = renderer.Render(c, views.InternalServerError(apiErr), 500)
	}
}
