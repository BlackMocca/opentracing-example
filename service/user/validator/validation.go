package validator

import "github.com/labstack/echo/v4"

type Validation struct {
}

func (v Validation) ValidateFetchUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		return next(c)
	}
}
