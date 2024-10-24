package main

import (
	"context"
	"echo-demo/config"
	"echo-demo/db"
	"echo-demo/users"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
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

	if len(os.Args) > 1 {
		if err := config.Etcd(os.Args[1]); err != nil {
			e.Logger.Fatal("Config from Etcd: ", err)
		}
	} else {
		if err := config.Init(); err != nil {
			e.Logger.Fatal("Config from File: ", err)
		}
	}

	if err := db.ConnInit(); err != nil {
		e.Logger.Fatal("Database: ", err)
	}

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	gv := e.Group("/v1")
	gv.POST("/auth", users.Auth)

	gu := gv.Group("/users")
	gu.Use(echojwt.WithConfig(echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(users.JwtCustomClaims)
		},
		SigningKey: config.VerifyKey(),
	}))
	gu.GET("", users.GetAll)
	gu.GET("/:id", users.GetOne)
	gu.POST("", users.Create)
	gu.PUT("/:id", users.Update)
	gu.DELETE("/:id", users.Delete)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		err := e.Start(config.ServerAddr())
		if err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("Server close: ", err)
		}
	}()

	//wait for signals to gracefully shutdown the server.
	<-ctx.Done()

	fmt.Print("Shutting down the server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	fmt.Println("done.")
}
