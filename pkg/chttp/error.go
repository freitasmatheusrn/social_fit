package chttp

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/freitasmatheusrn/social-fit/internal/components"
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func CustomHTTPErrorHandler(defaultHandler echo.HTTPErrorHandler) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		code := http.StatusInternalServerError
		message := http.StatusText(code)

		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
			if he.Message != nil {
				message = fmt.Sprintf("%v", he.Message)
			}
		}

		c.Logger().Error(err)

		accept := c.Request().Header.Get(echo.HeaderAccept)
		wantsHTML := strings.Contains(accept, echo.MIMETextHTML)
		wantsJSON := strings.Contains(accept, echo.MIMEApplicationJSON)
		if wantsJSON {
			defaultHandler(err, c)
			return
		}
		if wantsHTML {
			if code == http.StatusNotFound || code == http.StatusInternalServerError {
				errorPage := fmt.Sprintf("%d.html", code)
				if err := c.File(errorPage); err != nil {
					c.Logger().Error(err)
				}
				return
			}
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
			c.Response().WriteHeader(code)

			if err := components.ToastError(code, message).Render(c.Request().Context(), c.Response().Writer); err != nil {
				c.Logger().Error(err)
			}
			return
		}
		defaultHandler(err, c)
	}
}
