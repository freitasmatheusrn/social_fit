package main

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/freitasmatheusrn/social-fit/assets"
	"github.com/freitasmatheusrn/social-fit/config"
	"github.com/freitasmatheusrn/social-fit/internal/user"
	"github.com/freitasmatheusrn/social-fit/pkg/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/freitasmatheusrn/social-fit/pkg/handlers"
)

func (app *application) mount() http.Handler {
	e := echo.New()
	e.HTTPErrorHandler = handlers.CustomErrorHandler
	e.StaticFS("/assets", assets.Files)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	usrService := user.NewService(user.NewRepo(app.db))
	userHandler := user.NewHandler(usrService, app.config.JWTSecret)
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTCustomClaims)
		},
		SigningKey: []byte(app.config.JWTSecret),
		TokenLookup: "header:Authorization,cookie:access_token",
	}
	e.GET("/signup", userHandler.SignUpPage)
	e.GET("/login", userHandler.LoginPage)
	e.POST("/create_user", userHandler.CreateUser)
	e.POST("/sign_in", userHandler.Signin)
	authenticated := e.Group("/dashboard")
	authenticated.Use(echojwt.WithConfig(config))
	authenticated.GET("/home", userHandler.Home)

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
