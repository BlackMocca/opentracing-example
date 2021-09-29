package http

import (
	"net/http"
	"sync"
	"time"

	"github.com/Blackmocca/opentracing-example/service/user"
	"github.com/labstack/echo/v4"
)

type userHandler struct {
	userUs user.UserUsecase
}

func NewUserHandler(userUs user.UserUsecase) user.UserHandler {
	return &userHandler{
		userUs: userUs,
	}
}

func (u userHandler) FetchAll(c echo.Context) error {
	var ctx = c.Request().Context()
	var args = new(sync.Map)

	time.Sleep(time.Duration(1 * time.Second))

	users, err := u.userUs.FetchAll(ctx, args)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp := map[string]interface{}{
		"users": users,
	}
	return c.JSON(http.StatusOK, resp)
}
