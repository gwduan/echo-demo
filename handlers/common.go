package handlers

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func claims(c echo.Context) *JwtCustomClaims {
	token := c.Get("user").(*jwt.Token)
	return token.Claims.(*JwtCustomClaims)
}

func BadRequestErr(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return echo.NewHTTPError(http.StatusBadRequest, msg)
}

func NotFoundErr(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return echo.NewHTTPError(http.StatusNotFound, msg)
}

func UnauthorizedErr(format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	return echo.NewHTTPError(http.StatusUnauthorized, msg)
}
