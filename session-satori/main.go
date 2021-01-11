package main

import (
	"net/http"
	
	"github.com/satori/go.uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type User struct {
	UserID int64
	Email  string
}

type SessionData struct {
	SessionID string
	UserID    int64
}

var sessionDB []SessionData
var userDB []User

func generateSession() string {
	id := uuid.NewV4()
	sessionDB = append(sessionDB, SessionData{SessionID: id.String()})
	return id.String()
}

func getSession(c echo.Context) error {
	cook, err := c.Cookie("sessionID")

	if err != nil {
		id := generateSession()

		c.SetCookie(&http.Cookie{
			Name:     "sessionID",
			Value:    id,
			HttpOnly: true,
		})

		return c.NoContent(http.StatusOK)
	}

	var userID int64
	for _, ses := range sessionDB {
		if ses.SessionID == cook.Value {
			userID = ses.UserID
		}
	}

	for _, us := range userDB {
		if us.UserID == userID {
			return c.JSON(http.StatusOK, us)
		}
	}

	return c.NoContent(http.StatusOK)
}

func setSession(c echo.Context) error {
	var user User
	err := c.Bind(&user)
	if err != nil {
		c.Logger().Error(err.Error())
	}

	cook, err := c.Cookie("sessionID")

	var id string
	if err != nil {
		id = generateSession()

		c.SetCookie(&http.Cookie{
			Name:     "sessionID",
			Value:    id,
			HttpOnly: true,
		})

	} else {
		id = cook.Value
	}

	for _, ses := range sessionDB {
		if ses.SessionID == id {
			ses.UserID = user.UserID
		}
	}

	return c.NoContent(http.StatusOK)
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	e.GET("/get", getSession)
	e.POST("/set", setSession)

	e.Logger.Fatal(e.Start(":1000"))
}
