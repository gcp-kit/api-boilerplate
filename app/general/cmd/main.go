package main

import (
	"api-boilerplate/app/general/internal/config"
	"api-boilerplate/app/general/internal/initialize"
	"api-boilerplate/app/general/internal/interfaces/openapi"
	"api-boilerplate/app/general/internal/interfaces/props"
	"api-boilerplate/app/general/internal/server"
	"api-boilerplate/app/internal/middleware/validator"
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"

	"github.com/brpaz/echozap"
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

var (
	internalPort = func() string {
		p := os.Getenv("PORT")
		if p == "" {
			p = "1234"
		}
		return p
	}()
)

func main() {
	ctx := context.Background()
	e := echo.New()
	e.HTTPErrorHandler = initialize.ErrorHandler()

	e.Use(echoMiddleware.Recover())

	// log settings
	zapLoggerConfig := zap.NewProductionConfig()
	zapLoggerConfig.EncoderConfig.LevelKey = "severity"
	zapLogger, err := zapLoggerConfig.Build()
	if err != nil {
		panic(err)
	}
	defer zapLogger.Sync()
	e.Use(echozap.ZapLogger(zapLogger))
	e.Logger.SetLevel(log.INFO)

	useCORSMiddleware(e, "*")
	cfg := readConfig(ctx)
	ed := newExternalDependencies()
	useRequestValidatorMiddleware(ed, e)
	cp := newControllerProps(cfg)
	setHandlers(e, cp, "")

	// wait for shutdown signal
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt)
		<-sig
		e.Logger.Info("shutting down server...")
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	}()

	// start server
	e.Logger.Fatal(e.Start(":" + internalPort))
}

func useCORSMiddleware(e *echo.Echo, corsAllowOrigins string) {
	e.Use(
		echoMiddleware.CORSWithConfig(
			echoMiddleware.CORSConfig{
				AllowOrigins: lo.Uniq(strings.Split(corsAllowOrigins, ",")),
				AllowHeaders: []string{
					echo.HeaderOrigin,
					echo.HeaderXRequestedWith,
					echo.HeaderContentType,
					echo.HeaderAccept,
					echo.HeaderCookie,
					echo.HeaderAuthorization,
				},
				AllowMethods: []string{
					http.MethodGet,
					http.MethodPut,
					http.MethodPost,
					http.MethodDelete,
					http.MethodPatch,
					http.MethodOptions,
				},
				AllowCredentials: true,
			},
		),
	)
}

func readConfig(ctx context.Context) *config.Config {
	cfg, err := config.ReadConfig(ctx)
	if err != nil {
		panic(err)
	}
	return cfg
}

func newExternalDependencies() *initialize.ExternalDependencies {
	ed, err := initialize.NewExternalDependencies()
	if err != nil {
		panic(err)
	}
	return ed
}

func useRequestValidatorMiddleware(ed *initialize.ExternalDependencies, e *echo.Echo) {
	ovm, err := validator.NewMiddleware(ed.OpenAPISpec())
	if err != nil {
		e.Logger.Fatalf("failed to create validator middleware: %+v", err)
	}
	e.Use(ovm.Middleware)
}

func newControllerProps(cfg *config.Config) *props.ControllerProps {
	usecases := initialize.NewUsecases()
	cp := initialize.NewControllerProps(cfg, usecases)
	return cp
}

func setHandlers(e *echo.Echo, cp *props.ControllerProps, basePath string) {
	strictServer := server.NewServer(cp)
	si := openapi.NewStrictHandler(strictServer, []openapi.StrictMiddlewareFunc{})
	openapi.RegisterHandlersWithBaseURL(e, si, basePath)
}
