package api

import (
	"errors"
	"github.com/labstack/echo/v4"
)

func ErrorResponseWriter(c echo.Context, code int, err error) {
	jsonResponseWriter(c, code, Response{Status: "Error", Error: err.Error()})
}

func ResponseWriter(c echo.Context, payload interface{}) {
	jsonResponseWriter(c, 200, Response{Status: "OK", Data: payload})
}

func jsonResponseWriter(c echo.Context, code int, raw interface{}) {
	if err := c.JSON(code, raw); err != nil {
		ErrorResponseWriter(c, 500, errors.New("JSON Marshal error"))
	}
}
