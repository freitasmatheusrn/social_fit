package user

import (
	"net/http"

	"github.com/freitasmatheusrn/social-fit/internal/user/userpgs"
	"github.com/freitasmatheusrn/social-fit/pkg/renderer"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) SignUpPage(c echo.Context) error {
	return renderer.Render(c, userpgs.Signup())
}

func (h *Handler) CreateUser(c echo.Context) error {
	var s SignupRequest
	err := c.Bind(s)
	if err != nil{
		return echo.NewHTTPError(http.StatusUnprocessableEntity, "dados inválidos", err)
	}
	response, err := h.Service.Signup(c.Request().Context(), s)
	if err != nil{
		return echo.NewHTTPError(http.StatusBadRequest, "Erro ao criar conta", err)
	}

	return renderer.Respond(c, userpgs.Home(), response, http.StatusCreated)
}


