package users

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetOne(c echo.Context) error {
	return c.String(http.StatusOK, "GetUser")
}

func GetAll(c echo.Context) error {
	return c.String(http.StatusOK, "GetAllUsers")
}

func Create(c echo.Context) error {
	return c.String(http.StatusOK, "CreateUser")
}

func Update(c echo.Context) error {
	return c.String(http.StatusOK, "UpdateUser")
}

func Delete(c echo.Context) error {
	return c.String(http.StatusOK, "DeleteUser")
}
