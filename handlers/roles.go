package handlers

import (
	"echo-demo/config"
	"echo-demo/db"
	"echo-demo/users"
	"net/http"
	"strconv"
	"time"

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

func CreateRole(c echo.Context) error {
	logID, err := loginID(c)
	if err != nil {
		return err
	}
	if logID != 1 {
		return UnauthorizedErr("Admin Required")
	}

	uIn := new(users.Input)
	if err := c.Bind(uIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Data Invalid")
	}
	if err := c.Validate(uIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Validation Faild")
	}

	uOut, err := users.NewOne(uIn.Name, uIn.Password, uIn.Age, time.Now())
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == db.ErrDupRows {
			return BadRequestErr("User(%s) Duplicate", uIn.Name)
		}
		return err
	}

	return c.JSON(http.StatusCreated, uOut)
}

func GetOneRole(c echo.Context) error {
	if _, err := loginID(c); err != nil {
		return err
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

func GetAllRoles(c echo.Context) error {
	if _, err := loginID(c); err != nil {
		return err
	}

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit <= 0 {
		limit = config.RecordLimit()
	}
	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = config.RecordOffset()
	}

	uOuts, err := users.GetAll(int64(limit), int64(offset))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return err
	}

	return c.JSON(http.StatusOK, uOuts)
}

func UpdateRole(c echo.Context) error {
	logID, err := loginID(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Id(%s) Invalid", c.Param("id"))
	}

	if logID != 1 && logID != int64(id) {
		return UnauthorizedErr("Admin or User(id:%d) Required", id)
	}

	uIn := new(users.Input)
	if err := c.Bind(uIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Data Invalid")
	}
	if err := c.Validate(uIn); err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Validation Faild")
	}

	uOut, err := users.UpdateOne(int64(id), uIn.Name, uIn.Password, uIn.Age)
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

func DeleteRole(c echo.Context) error {
	logID, err := loginID(c)
	if err != nil {
		return err
	}
	if logID != 1 {
		return UnauthorizedErr("Admin Required")
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Echo().Logger.Debug(err)
		return BadRequestErr("Id(%s) Invalid", c.Param("id"))
	}

	err = users.DeleteOne(int64(id))
	if err != nil {
		c.Echo().Logger.Debug(err)
		if err == db.ErrNotFound {
			return NotFoundErr("User(id:%d) Not Found", id)
		}
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func loginID(c echo.Context) (int64, error) {
	sess, err := session.Get("session", c)
	if err != nil {
		return 0, err
	}

	v, ok := sess.Values["role_id"]
	if !ok {
		return 0, UnauthorizedErr("Please login")
	}

	id, ok := v.(int64)
	if !ok {
		return 0, BadRequestErr("Session Invalid")
	}

	return id, nil
}
