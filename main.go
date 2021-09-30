package main

import (
	"fmt"
	"net"
	"net/http"

	"git.innovasive.co.th/backend/helper"
	helperMiddl "git.innovasive.co.th/backend/helper/middleware"
	helperRoute "git.innovasive.co.th/backend/helper/route"
	myMiddL "github.com/Blackmocca/opentracing-example/middleware"
	route "github.com/Blackmocca/opentracing-example/route"
	user_grpc_handler "github.com/Blackmocca/opentracing-example/service/user/grpc"
	user_handler "github.com/Blackmocca/opentracing-example/service/user/http"
	"github.com/Blackmocca/opentracing-example/service/user/repository"
	"github.com/Blackmocca/opentracing-example/service/user/usecase"
	user_validator "github.com/Blackmocca/opentracing-example/service/user/validator"
	_util_tracing "github.com/Blackmocca/opentracing-example/utils/opentracing"
	"github.com/Blackmocca/opentracing-example/utils/psql"
	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	echoMiddL "github.com/labstack/echo/v4/middleware"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

var (
	GRPC_PORT         = helper.GetENV("GRPC_PORT", "3100")
	PSQL_DATABASE_URL = helper.GetENV("PSQL_DATABASE_URL", "postgres://postgres:postgres@psql_db:5432/app_example?sslmode=disable")
)

func sqlDB(con string) *psql.Client {
	db, err := psql.NewPsqlConnection(con)
	if err != nil {
		panic(err)
	}
	return db
}

func sqlDBWithTracing(con string, tracer opentracing.Tracer) *psql.Client {
	db, err := psql.NewPsqlWithTracingConnection(con, tracer)
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	/* init tracing*/
	tracer, closer := _util_tracing.Init("opentracing-example")
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	psqlClient := sqlDBWithTracing(PSQL_DATABASE_URL, tracer)

	/* init grpc */
	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			otgrpc.OpenTracingServerInterceptor(tracer),
		),
		grpc.StreamInterceptor(
			otgrpc.OpenTracingStreamServerInterceptor(tracer),
		),
	)
	defer server.GracefulStop()

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
	psqlRepo := repository.NewPsqlUserRepository(psqlClient)

	userUs := usecase.NewUserUsecase(userRepo, psqlRepo)

	userHandler := user_handler.NewUserHandler(userUs)

	userValidator := user_validator.Validation{}

	grpcUserHandler := user_grpc_handler.NewGRPCHandler(userUs)

	r := route.NewRoute(e, middL)
	r.RegisterRouteUser(userHandler, userValidator)

	grpcR := route.NewGRPCRoute(server)
	grpcR.RegisterUserHandler(grpcUserHandler)

	/* serve gprc */
	go func() {
		if r := recover(); r != nil {
			fmt.Println(r.(error))
		}
		startGRPCServer(server)
	}()

	/* serve echo */
	port := fmt.Sprintf(":%s", "3000")
	e.Logger.Fatal(e.Start(port))
}

func startGRPCServer(server *grpc.Server) {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%s", GRPC_PORT))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	/* serve grpc */
	fmt.Println(fmt.Sprintf("Start grpc Server [::%s]", GRPC_PORT))
	if err := server.Serve(listen); err != nil {
		panic(err)
	}
}
