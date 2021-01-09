package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/boj/redistore.v1"
)

func fetchStore() *redistore.RediStore {
	// Fetch new store.
	store, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte("secret-key"))
	if err != nil {
		panic(err)
	}

	return store
}

func fetchSession(c echo.Context, store *redistore.RediStore) *sessions.Session {
	// Get a session.
	session, err := store.Get(c.Request(), "sessionID")
	if err != nil {
		fmt.Println(err.Error())
	}

	return session
}

func getSession(c echo.Context) error {
	// Fetch new store.
	store := fetchStore()

	// Get a session.
	session := fetchSession(c, store)

	// Save session
	if err := sessions.Save(c.Request(), c.Response()); err != nil {
		log.Fatalf("Error saving session: %v", err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("sessionID: %s \n userAgent: %s \n ip: %s",
		session.ID,
		session.Values["userAgent"],
		session.Values["ip"]))
}

// session in redis
func updateSession(c echo.Context) error {
	// Fetch new store.
	store := fetchStore()

	// Get a session.
	session := fetchSession(c, store)

	// Change session storage configuration for MaxAge = 10 days.
	store.SetMaxAge(10 * 24 * 3600)

	// data := LogFields{
	// 	RequestID: c.Request().Header.Get(echo.HeaderXRequestID),
	// 	Method:    c.Request().Method,
	// 	URI:       c.Request().RequestURI,
	// 	IP:        c.Request().RemoteAddr,
	// 	RemoteIP:  c.RealIP(),
	// 	Host:      c.Request().Host,
	// 	Status:    c.Response().Status,
	// 	Size:      c.Response().Size,
	// 	UserAgent: c.Request().UserAgent(),
	// }

	// Add a value.
	session.Values["userAgent"] = c.Request().UserAgent()
	session.Values["ip"] = c.Request().RemoteAddr

	// Save session
	if err := sessions.Save(c.Request(), c.Response()); err != nil {
		log.Fatalf("Error saving session: %v", err)
	}

	return c.String(http.StatusOK, fmt.Sprintf("sessionID: %s \n userAgent: %s \n ip: %s",
		session.ID,
		session.Values["userAgent"],
		session.Values["ip"]))
}

func deleteSession(c echo.Context) error {
	// Fetch new store.
	store := fetchStore()

	// Get a session.
	session := fetchSession(c, store)

	// Delete session.
	session.Options.MaxAge = -1
	if err := sessions.Save(c.Request(), c.Response()); err != nil {
		log.Fatalf("Error saving session: %v", err)
	}

	// Get a session.
	session = fetchSession(c, store)

	return c.String(http.StatusOK, fmt.Sprintf("sessionID: %s was deleted", session.ID))
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set Session redis
	e.GET("/session/get", getSession)

	// Add value
	e.GET("/session/put", updateSession)

	// Set Session redis
	e.GET("/session/delete", deleteSession)

	e.Logger.Fatal(e.Start(":1323"))
}
