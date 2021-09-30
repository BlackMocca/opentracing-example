package route

import (
	"github.com/Blackmocca/opentracing-example/middleware"
	"github.com/Blackmocca/opentracing-example/service/user"
	_user_validator "github.com/Blackmocca/opentracing-example/service/user/validator"
	"github.com/labstack/echo/v4"
)

type Route struct {
	e     *echo.Echo
	middl middleware.GoMiddlewareInf
}

func NewRoute(e *echo.Echo, middl middleware.GoMiddlewareInf) *Route {
	return &Route{e: e, middl: middl}
}

func (r *Route) RegisterRouteUser(handler user.UserHandler, validation _user_validator.Validation) {
	r.e.GET("/users", handler.FetchAll)
	r.e.GET("/users/:user_id/cover", handler.GetCover)
	r.e.GET("/internal-error", handler.InternalError)
	r.e.GET("/conflict", handler.Conflict)
}
