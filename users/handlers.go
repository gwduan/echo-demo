package users

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
)

const (
	DEFAULT_LIMIT  = 5
	DEFAULT_OFFSET = 0
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
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = DEFAULT_LIMIT
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = DEFAULT_OFFSET
	}

	uOuts, err := getAllOut(int64(limit), int64(offset))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return err
	}

	return c.JSON(http.StatusOK, uOuts)
}

func Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Please provide valid id (number)")
	}

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

	uOut, err := updateOneOut(int64(id), uIn.Name, uIn.Password, uIn.Age)
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound,
				"User(id:"+c.Param("id")+") Not Found")
		}
		if e, ok := err.(*mysql.MySQLError); ok {
			//Duplicate
			if e.Number == 1062 {
				return echo.NewHTTPError(http.StatusBadRequest,
					"User(name:"+uIn.Name+") Duplicate")
			}
		}
		return err
	}

	return c.JSON(http.StatusOK, uOut)
}

func Delete(c echo.Context) error {
	return c.String(http.StatusOK, "DeleteUser")
}
