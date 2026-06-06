package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	api "github.com/masraga/golang-echo-boilerplate/generated/api"
	"github.com/masraga/golang-echo-boilerplate/internal/app/backend/server"
	"github.com/masraga/golang-echo-boilerplate/internal/echo/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"
)

func InitApi() {
	e := echo.New()
	ctx := context.Background()

	config, server := Initialize(ctx)

	api.RegisterHandlers(e, server)
	RegisterMiddlewares(e, server)

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		sig := <-sigChan

		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal("Error shutting down server: ", err)
		}

		fmt.Printf("Received signal: %s. Shutting down server...\n", sig)
	}()

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.AppPort)))
}

func RegisterMiddlewares(e *echo.Echo, server *server.Server) {
	e.Use(echomiddleware.RequestID())
	e.Use(echomiddleware.RequestLogger())
	e.Use(echomiddleware.Recover())
	e.Use(middleware.AuthValidationFilterMiddleware(server.AuthService))
	e.HTTPErrorHandler = middleware.HTTPErrorHandler
	e.Use(oapimiddleware.OapiRequestValidatorWithOptions(middleware.OapiGetSwagger(), middleware.OapiValidatorOpt()))
}
