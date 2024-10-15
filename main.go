package main

import (
	"echo-demo/db"
	"echo-demo/users"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}

	return nil
}

func main() {
	e := echo.New()
	e.Validator = &CustomValidator{validator: validator.New()}
	e.Debug = true

	if err := db.ConnInit(); err != nil {
		e.Logger.Fatal("Database: ", err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	gv := e.Group("/v1")
	gu := gv.Group("/users")
	gu.GET("", users.GetAll)
	gu.GET("/:id", users.GetOne)
	gu.POST("", users.Create)
	gu.PUT("/:id", users.Update)
	gu.DELETE("/:id", users.Delete)

	e.Logger.Fatal(e.Start(":8080"))
}
