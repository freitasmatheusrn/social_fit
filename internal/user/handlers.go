package user

import (
	"net/http"

	"github.com/freitasmatheusrn/social-fit/internal/user/userpgs"
	"github.com/freitasmatheusrn/social-fit/pkg/renderer"
	"github.com/labstack/echo/v4"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SignUpPage(c echo.Context) error {
	return c.Render(http.StatusOK, "", userpgs.Signup())
}

func (h *Handler) CreateUser(c echo.Context) error {
	var s SignupRequest
	err := c.Bind(s)
	if err != nil{
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "dados inválidos", err)
	}
	return renderer.Respond(c, userpgs.Home(), s, http.StatusCreated)
}


