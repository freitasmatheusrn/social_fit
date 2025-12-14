package main

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/freitasmatheusrn/social-fit/config"
	"github.com/freitasmatheusrn/social-fit/internal/user"
	"github.com/freitasmatheusrn/social-fit/pkg/auth"
	"github.com/freitasmatheusrn/social-fit/pkg/chttp"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (app *application) mount() http.Handler {
	e := echo.New()
	defaultErrorHandler := e.HTTPErrorHandler
	e.HTTPErrorHandler = chttp.CustomHTTPErrorHandler(defaultErrorHandler)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	userHandler := user.NewHandler()
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTCustomClaims)
		},
		SigningKey: app.config.JWTSecret,
	}
	e.GET("/signup", userHandler.SignUpPage)
	e.POST("/create_user", userHandler.CreateUser)
	r := e.Group("/api")
	r.Use(echojwt.WithConfig(config))

	e.Logger.Fatal(e.Start(":8080"))
	return e
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.WebServerPort,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.WebServerPort)

	return srv.ListenAndServeTLS("server.crt", "server.key")
}

type application struct {
	config config.Config
	logger *slog.Logger
	db     *pgx.Conn
}
