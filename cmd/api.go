package main

import (
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/freitasmatheusrn/social-fit/assets"
	"github.com/freitasmatheusrn/social-fit/config"
	"github.com/freitasmatheusrn/social-fit/internal/events"
	"github.com/freitasmatheusrn/social-fit/internal/types"
	"github.com/freitasmatheusrn/social-fit/internal/user"
	"github.com/freitasmatheusrn/social-fit/pkg/auth"
	"github.com/freitasmatheusrn/social-fit/pkg/handlers"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func (app *application) mount() http.Handler {
	e := echo.New()
	e.HTTPErrorHandler = handlers.CustomErrorHandler
	e.StaticFS("/assets", assets.Files)
	cfg := zap.NewDevelopmentConfig()

	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logger, _ := cfg.Build()
	defer logger.Sync()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.Info("request",
					zap.String("URI", v.URI),
					zap.Int("status", v.Status),
				)

			} else {
				logger.Error("ERROR LOG:",
					zap.String("URI", v.URI),
					zap.String("erro:", v.Error.Error()),
				)
			}

			return nil
		},
	}))
	e.Use(middleware.Recover())
	evtService := events.NewService(events.NewRepo(app.db))
	eventHandler := events.NewHandler(evtService)
	usrService := user.NewService(user.NewRepo(app.db), evtService)
	userHandler := user.NewHandler(usrService, app.config.JWTSecret)
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(auth.JWTCustomClaims)
		},
		SigningKey:  []byte(app.config.JWTSecret),
		TokenLookup: "header:Authorization,cookie:access_token",
		SuccessHandler: func(c echo.Context) {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*auth.JWTCustomClaims)
			currentUser := types.CurrentUser{
				ID:    claims.UserID,
				Name:  claims.Name,
				Email: claims.Email,
			}

			types.SetCurrentUser(c, currentUser)
		},
	}
	e.GET("/signup", userHandler.SignUpPage)
	e.GET("/login", userHandler.LoginPage)
	e.POST("/create_user", userHandler.CreateUser)
	e.POST("/sign_in", userHandler.Signin)
	authenticated := e.Group("/dashboard")
	authenticated.Use(echojwt.WithConfig(config))
	authenticated.GET("/home", userHandler.Home)
	authenticated.GET("/events/form", eventHandler.Form)
	authenticated.POST("/events/create_event", eventHandler.Create)
	authenticated.GET("/events/:event_id", eventHandler.Show)

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
