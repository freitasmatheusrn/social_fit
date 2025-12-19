package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/freitasmatheusrn/social-fit/internal/views"
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

	if isHTMX {
		handleHTMXError(c, apiErr)
	} else {
		c.JSON(apiErr.Code, apiErr)
	}
}

func handleHTMXError(c echo.Context, apiErr *rest.ApiErr) {
	c.Response().Status = apiErr.Code

	var component templ.Component

	switch apiErr.Code {
	case http.StatusBadRequest:
		component = views.BadRequest(apiErr)
	case http.StatusUnauthorized:
		component = views.Unauthorized(apiErr)
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
