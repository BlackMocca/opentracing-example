package main

import (
	"fmt"
	"net/http"

	helperMiddl "git.innovasive.co.th/backend/helper/middleware"
	helperRoute "git.innovasive.co.th/backend/helper/route"
	myMiddL "github.com/Blackmocca/opentracing-example/middleware"
	route "github.com/Blackmocca/opentracing-example/route"
	user_handler "github.com/Blackmocca/opentracing-example/service/user/http"
	"github.com/Blackmocca/opentracing-example/service/user/repository"
	"github.com/Blackmocca/opentracing-example/service/user/usecase"
	user_validator "github.com/Blackmocca/opentracing-example/service/user/validator"
	_util_tracing "github.com/Blackmocca/opentracing-example/utils/opentracing"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	echoMiddL "github.com/labstack/echo/v4/middleware"
	"github.com/opentracing/opentracing-go"
)

func main() {
	/* init tracing*/
	tracer, closer := _util_tracing.Init("opentracing-example")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	e := echo.New()
	e.HTTPErrorHandler = helperMiddl.SentryCapture(e)
	helperRoute.RegisterVersion(e)

	e.Use(echoMiddL.Logger())
	e.Use(echoMiddL.Recover())
	e.Use(sentryecho.New(sentryecho.Options{Repanic: true}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	middL := myMiddL.InitMiddleware("test")
	e.Use(echoMiddL.Recover())
	e.Use(echoMiddL.CORSWithConfig(echoMiddL.CORSConfig{
		Skipper:      echoMiddL.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middL.InitContextIfNotExists)
	e.Use(middL.InputForm)
	e.Use(middL.SetTracer)

	userRepo := repository.NewUserRepository()
	userUs := usecase.NewUserUsecase(userRepo)
	userHandler := user_handler.NewUserHandler(userUs)
	userValidator := user_validator.Validation{}

	r := route.NewRoute(e, middL)
	r.RegisterRouteUser(userHandler, userValidator)

	/* serve echo */
	port := fmt.Sprintf(":%s", "3000")
	e.Logger.Fatal(e.Start(port))
}
