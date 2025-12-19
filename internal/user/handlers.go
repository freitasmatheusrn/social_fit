package user

import (
	"log"
	"net/http"

	"github.com/freitasmatheusrn/social-fit/internal/user/userpgs"
	"github.com/freitasmatheusrn/social-fit/pkg/auth"
	"github.com/freitasmatheusrn/social-fit/pkg/renderer"
	"github.com/freitasmatheusrn/social-fit/pkg/rest"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	Service *Service
	secret  string
}

func NewHandler(s *Service, secret string) *Handler {
	return &Handler{
		Service: s,
		secret:  secret,
	}
}

func (h *Handler) SignUpPage(c echo.Context) error {
	return renderer.Render(c, userpgs.Signup(nil), 200)
}
func (h *Handler) LoginPage(c echo.Context) error {
	return renderer.Render(c, userpgs.Signin(nil), 200)
}

func (h *Handler) CreateUser(c echo.Context) error {
	var s SignupRequest
	err := c.Bind(&s)
	if err != nil {
		return err
	}
	response, restErr := h.Service.Signup(c.Request().Context(), s)
	if restErr != nil {
		data := userpgs.SignupFormData{
			Name:      s.Name,
			Email:     s.Email,
			CPF:       s.Cpf,
			Phone:     s.Phone,
			BirthDate: s.BirthDate,
		}
		return renderer.Respond(c, userpgs.SignupForm(restErr, data), s, http.StatusOK)
	}
	claims := auth.NewClaims(response.Name, response.Email, response.Admin)

	token, err := auth.GenerateJWT(claims, h.secret)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to generate token",
		})
	}
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)
	c.Response().Header().Set("HX-Redirect", "/home")
	return c.NoContent(http.StatusCreated)
}


func (h *Handler) Signin(c echo.Context) error{
	var s SigninRequest
	err := c.Bind(&s)
	if err != nil {
		return err
	}
	response, restErr := h.Service.Login(c.Request().Context(), s)
	if restErr != nil {
		data := userpgs.SigninFormData{
			Email:     s.Email,
		}
		log.Println(restErr)
		return renderer.Respond(c, userpgs.SigninForm(restErr, data), s, http.StatusOK)
	}
	claims := auth.NewClaims(response.Name, response.Email, response.Admin)

	token, err := auth.GenerateJWT(claims, h.secret)
	if err != nil {
		return rest.NewInternalServerError("erro interno")
	}
	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}

	c.SetCookie(cookie)
	c.Response().Header().Set("HX-Redirect", "/dashboard/home")
	return c.NoContent(http.StatusCreated)
}

func (h *Handler) Home(c echo.Context) error{

	return renderer.Render(c, userpgs.Home(), http.StatusOK)
}