package api

import (
	"github.com/labstack/echo/v4"
)

func RequestLoader(c echo.Context, v interface{}) bool {
	if err := c.Bind(v); err != nil {
		ErrorResponseWriter(c, 400, err)
		return false
	}
	return true
}
