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

	session, err := mgo.Dial("root:example@13.250.119.252")
	if err != nil {
		e.Logger.Fatal(err)
		return
	}

	h := &handler{
		m: session,
	}

	e.Use(middleware.Logger())
	e.GET("/todos", h.list)
	e.POST("/todos", h.create)
	e.PUT("/todos/:id", h.done)
	e.GET("/todos/:id", h.view)
	e.Logger.Fatal(e.Start(":1323"))
}

type handler struct {
	m *mgo.Session
}

type todo struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Topic string        `json:"topic" bson:"topic"`
	Done  bool          `json:"done" bson:"done"`
}

func (h *handler) done(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	var t todo
	col := session.DB("workshop").C("todos")
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}
	t.Done = true
	if err := col.UpdateId(id, t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (h *handler) view(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	id := bson.ObjectIdHex(c.Param("id"))
	var t todo
	col := session.DB("workshop").C("todos")
	if err := col.FindId(id).One(&t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func (h *handler) list(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()
	var ts []todo
	col := session.DB("workshop").C("todos")
	if err := col.Find(nil).All(&ts); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ts)
}

func (h *handler) create(c echo.Context) error {
	session := h.m.Copy()
	defer session.Close()

	var t todo
	if err := c.Bind(&t); err != nil {
		return err
	}
	t.ID = bson.NewObjectId()

	col := session.DB("workshop").C("todos")
	if err := col.Insert(t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}
