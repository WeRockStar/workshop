package main

import (
	"net/http"

	"github.com/globalsign/mgo/bson"

	"github.com/globalsign/mgo"
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
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic" bson:"topic"`
	Done  bool          `json:"done" bson:"done"`
}

func create(c echo.Context) error {
	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}
	session, err := mgo.Dial("root:example@13.250.119.252")
	if err != nil {
		return err
	}
	t.ID = bson.NewObjectId()
	col := session.DB("workshop").C("todos")
	if err := col.Insert(t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}
