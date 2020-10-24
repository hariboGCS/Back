package rest

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

//RunAPI used in main.go
func RunAPI(address string) {
	e := echo.New()
	h, _ := newHandler()

	e.Use(middleware.Logger())
	e.POST("/signup", h.Signup)
	e.POST("/signin", h.Signin)
	e.PUT("/users/:id", h.updateUser)
	e.DELETE("/users/:id", h.deleteUser)

	e.Logger.Fatal(e.Start(address))
}
