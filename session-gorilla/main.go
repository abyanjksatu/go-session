package main

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var SessionID = "sessionID"
var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(32)))

type CartData struct {
	ProductID   int64
	ProductName string
	Qty         int64
	Price       float64
}

type SessionData struct {
	ID      int64
	Name    string
	Email   string
	Address string
	Cart    CartData
}

func getSession(c echo.Context) error {
	sess, err := store.Get(c.Request(), SessionID)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	_, err = c.Cookie(SessionID)
	if err != nil {
		err = sess.Save(c.Request(), c.Response())
		if err != nil {
			c.Logger().Error(err.Error())
		}
	}

	return c.JSON(http.StatusOK, sess.Values["sessionData"])
}

func setSession(c echo.Context) error {
	var sessionData SessionData

	err := c.Bind(&sessionData)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	store.MaxAge(10 * 24 * 3600)
	
	sess, err := store.Get(c.Request(), SessionID)
	if err != nil {
		sess, err = store.New(c.Request(), SessionID)
		if err != nil {
			c.Logger().Error(err.Error())
		}
	}

	sess.Values["sessionData"] = sessionData

	// err = store.Save(c.Request(), c.Response(), sess)
	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		fmt.Println(err.Error())
	}

	return c.JSON(http.StatusOK, sess.Values["sessionData"])
}

func flushSession(c echo.Context) error {
	sess, err := store.Get(c.Request(), SessionID)
	if err != nil {
		c.Logger().Error(err.Error())
		return c.String(http.StatusNotFound, err.Error())
	}

	sess.Values["sessionData"] = nil

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		fmt.Println(err.Error())
	}

	return c.JSON(http.StatusOK, sess.Values["sessionData"])
}

func refreshSession(c echo.Context) error {
	sess, err := store.Get(c.Request(), SessionID)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	sess.Options.MaxAge = -1

	err = sess.Save(c.Request(), c.Response())
	if err != nil {
		c.Logger().Error(err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/session/get")
}

func main() {
	gob.Register(SessionData{})
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Get Session
	e.GET("/session/get", getSession)
	e.POST("/session/set", setSession)
	e.GET("/session/flush", flushSession)
	e.GET("/session/refresh", refreshSession)

	e.Logger.Fatal(e.Start(":1323"))
}
