package api

import (
	"github.com/ftCommunity/roboheart/internal/services/core/acm"
	"github.com/labstack/echo/v4"
)

func CheckACMAPICall(c echo.Context, acm acm.ACM, perm string, f func()) {
	req := TokenRequest{}
	if !RequestLoader(c, req) {
		return
	}
	if err, uae := acm.CheckTokenPermission(req.Token, perm); err != nil {
		if uae {
			ErrorResponseWriter(c, 403, err)
		} else {
			ErrorResponseWriter(c, 500, err)
		}
	} else {
		f()
	}
}
