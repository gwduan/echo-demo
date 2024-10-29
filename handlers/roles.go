package handlers

import (
	"echo-demo/db"
	"echo-demo/users"
	"net/http"
	"strconv"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

func Login(c echo.Context) error {
	aIn := new(users.AuthInput)
	if err := c.Bind(aIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Data Invalid")
	}
	if err := c.Validate(aIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Validation Faild")
	}

	uOut, err := users.Auth(aIn.Name, aIn.Password)
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == db.ErrNotFound {
			return UnauthorizedErr("Name|Password Incorrect")
		}
		return err
	}

	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}
	sess.Values["role_id"] = uOut.ID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, uOut)
}

func GetOneRole(c echo.Context) error {
	sess, err := session.Get("session", c)
	if err != nil {
		return err
	}
	if _, ok := sess.Values["role_id"]; !ok {
		return UnauthorizedErr("Please login")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Id(%s) Invalid", c.Param("id"))
	}

	uOut, err := users.GetOneByID(int64(id))
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == db.ErrNotFound {
			return NotFoundErr("User(id:%d) Not Found", id)
		}
		return err
	}

	return c.JSON(http.StatusOK, uOut)
}
