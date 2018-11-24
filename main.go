package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())

	e.POST("/todos", create)
	e.Logger.Fatal(e.Start(":1323"))
}

type todo struct {
	ID    string `json:id`
	Topic string `json:topic`
	Done  bool   `json:done`
}

func create(c echo.Context) error {
	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}
