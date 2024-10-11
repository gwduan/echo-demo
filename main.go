package main

import (
	"echo-demo/users"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

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
