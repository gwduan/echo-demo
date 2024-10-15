package users

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

func Create(c echo.Context) error {
	uIn := new(UserInput)
	if err := c.Bind(uIn); err != nil {
		c.Echo().Logger.Debug(err)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Please provide valid data")
	}
	if err := c.Validate(uIn); err != nil {
		c.Echo().Logger.Debug(err)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Validation error, please verify your data")
	}

	uOut, err := newOneOut(uIn.Name, uIn.Password, uIn.Age, time.Now())
	if err != nil {
		c.Echo().Logger.Debug(err)
		if e, ok := err.(*mysql.MySQLError); ok {
			//Duplicate
			if e.Number == 1062 {
				return echo.NewHTTPError(http.StatusBadRequest,
					"User(name:"+uIn.Name+") Duplicate")
			}
		}
		return err
	}

	return c.JSON(http.StatusCreated, uOut)
}

func GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Please provide valid id (number)")
	}

	uOut, err := getOneOutByID(int64(id))
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound,
				"User(id:"+c.Param("id")+") Not Found")
		}
		return err
	}

	return c.JSON(http.StatusOK, uOut)
}

func GetAll(c echo.Context) error {
	return c.String(http.StatusOK, "GetAllUsers")
}

func Update(c echo.Context) error {
	return c.String(http.StatusOK, "UpdateUser")
}

func Delete(c echo.Context) error {
	return c.String(http.StatusOK, "DeleteUser")
}
