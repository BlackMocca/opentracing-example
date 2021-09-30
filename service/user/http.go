package user

import (
	"github.com/labstack/echo/v4"
)

type UserHandler interface {
	FetchAll(c echo.Context) error
	GetCover(c echo.Context) error
	InternalError(c echo.Context) error
	Conflict(c echo.Context) error
}
