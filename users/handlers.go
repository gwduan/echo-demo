package users

import (
	"database/sql"
	"echo-demo/config"
	"net/http"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func Auth(c echo.Context) error {
	aIn := new(AuthInput)
	if err := c.Bind(aIn); err != nil {
		c.Echo().Logger.Debug(err)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Please provide valid data")
	}
	if err := c.Validate(aIn); err != nil {
		c.Echo().Logger.Debug(err)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Validation error, please verify your data")
	}

	uOut, err := auth(aIn.Name, aIn.Password)
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusUnauthorized,
				"User(name or password) Incorrect")
		}
		return err
	}

	claims := &JwtCustomClaims{
		uOut.ID,
		uOut.Name,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(config.SignKey())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, AuthOutput{uOut, tokenStr})
}

func Create(c echo.Context) error {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)
	authID := claims.ID
	if authID != 1 {
		return echo.NewHTTPError(http.StatusUnauthorized,
			"Only administrator can process")
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
		limit = config.RecordLimit()
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = config.RecordOffset()
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

	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)
	authID := claims.ID
	if authID != 1 && authID != int64(id) {
		return echo.NewHTTPError(http.StatusUnauthorized,
			"Only administrator or user(id:"+c.Param("id")+
				") can process")
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
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(*JwtCustomClaims)
	authID := claims.ID
	if authID != 1 {
		return echo.NewHTTPError(http.StatusUnauthorized,
			"Only administrator can process")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return echo.NewHTTPError(http.StatusBadRequest,
			"Please provide valid id (number)")
	}

	err = deleteOne(int64(id))
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound,
				"User(id:"+c.Param("id")+") Not Found")
		}
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
