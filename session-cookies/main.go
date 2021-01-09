package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const SESSION_ID = "sessionID"

// Note: Don't store your key in your source code. Pass it via an
// environmental variable, or flag (or both), and don't accidentally commit it
// alongside your code. Ensure your key is sufficiently random - i.e. use Go's
// crypto/rand or securecookie.GenerateRandomKey(32) and persist the result.
var store = sessions.NewCookieStore([]byte(securecookie.GenerateRandomKey(32)))

// session in cookies
func setSession(c echo.Context) error {

	store.MaxAge(10 * 24 * 3600)

	session, _ := store.Get(c.Request(), SESSION_ID)
	session.Values["message1"] = "hello"
	session.Values["message2"] = "world"

	err := session.Save(c.Request(), c.Response())
	if err != nil {
		fmt.Println(err.Error())
	}

	return c.Redirect(http.StatusTemporaryRedirect, "/get")
}

func getSession(c echo.Context) error {
	session, _ := store.Get(c.Request(), SESSION_ID)
	return c.String(http.StatusOK, fmt.Sprintf("%s %s", session.Values["message1"], session.Values["message2"]))
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/set", setSession)
	e.GET("/get", getSession)

	e.Logger.Fatal(e.Start(":1323"))
}
