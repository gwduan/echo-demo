package users

import (
	"echo-demo/config"
	"echo-demo/db"
	"fmt"
	"net/http"
	"strconv"
	"time"

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
		return BadRequestErr("Data Invalid")
	}
	if err := c.Validate(aIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Validation Faild")
	}

	uOut, err := auth(aIn.Name, aIn.Password)
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == db.ErrNotFound {
			return UnauthorizedErr("Name|Password Incorrect")
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
	if authID := claims(c).ID; authID != 1 {
		return UnauthorizedErr("Admin Required")
	}

	uIn := new(UserInput)
	if err := c.Bind(uIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Data Invalid")
	}
	if err := c.Validate(uIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Validation Faild")
	}

	uOut, err := newOneOut(uIn.Name, uIn.Password, uIn.Age, time.Now())
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == db.ErrDupRows {
			return BadRequestErr("User(%s) Duplicate", uIn.Name)
		}
		return err
	}

	return c.JSON(http.StatusCreated, uOut)
}

func GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Id(%s) Invalid", c.Param("id"))
	}

	uOut, err := getOneOutByID(int64(id))
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == db.ErrNotFound {
			return NotFoundErr("User(id:%d) Not Found", id)
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
		return BadRequestErr("Id(%s) Invalid", c.Param("id"))
	}

	if authID := claims(c).ID; authID != 1 && authID != int64(id) {
		return UnauthorizedErr("Admin or User(id:%d) Required", id)
	}

	uIn := new(UserInput)
	if err := c.Bind(uIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Data Invalid")
	}
	if err := c.Validate(uIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Validation Faild")
	}

	uOut, err := updateOneOut(int64(id), uIn.Name, uIn.Password, uIn.Age)
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == db.ErrNotFound {
			return NotFoundErr("User(id:%d) Not Found", id)
		} else if err == db.ErrDupRows {
			return BadRequestErr("User(%s) Duplicate", uIn.Name)
		}
		return err
	}

	return c.JSON(http.StatusOK, uOut)
}

func Delete(c echo.Context) error {
	if authID := claims(c).ID; authID != 1 {
		return UnauthorizedErr("Admin Required")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Id(%s) Invalid", c.Param("id"))
	}

	err = deleteOne(int64(id))
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == db.ErrNotFound {
			return NotFoundErr("User(id:%d) Not Found", id)
		}
		return err
	}

	return c.NoContent(http.StatusNoContent)
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
