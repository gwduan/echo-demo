package users

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Please provide valid id (number)")
	}

	u, err := getOneOutByID(int64(id))
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound,
				"User(id:"+c.Param("id")+") Not Found")
		}
		return err
	}

	return c.JSON(http.StatusOK, u)
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
