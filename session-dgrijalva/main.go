package main

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// var secretKey = "secretKey"

// Create the JWT key used to create the signature
var secretKey = []byte("secret-key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

// Create a struct to read the username and password from the request body
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func login(c echo.Context) error {
	var creds Credentials
	// Get the JSON body and decode into credentials
	err := c.Bind(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		return c.String(http.StatusUnauthorized, "Unauthorized")
	}

	// Declare the expiration time of the token
	// here, we have kept it as 5 minutes
	expirationTime := time.Now().Add(5 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare token with custom mapclaims
	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"foo": "bar",
	// 	"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	// })

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	c.SetCookie(&http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

	return c.String(http.StatusOK, "OK")
}

func accessible(c echo.Context) error {
	return c.String(http.StatusOK, "Accessible")
}

func restricted(c echo.Context) error {
	// user := c.Get("token").(*jwt.Token)
	// claims := user.Claims.(jwt.MapClaims)
	// name := claims["name"].(string)
	// admin := claims["admin"].(bool)
	// exp := claims["exp"].(float64)
	// msg := fmt.Sprintf("%s:%v:%f", name, admin, exp)

	cook, err := c.Cookie("token")
	if err != nil {
		return c.String(http.StatusUnauthorized, err.Error())
	}

	tokenVal := cook.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tokenVal, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return c.String(http.StatusUnauthorized, "Unauth")
		}
		return c.String(http.StatusBadRequest, "Bad")
	}

	if !tkn.Valid {
		return c.String(http.StatusUnauthorized, "Unauth")
	}

	return c.JSON(http.StatusOK, claims)
}

// MiddlewareJWT for middleware token
// var MiddlewareJWT = middleware.JWT([]byte(secretKey))

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Unauthenticated route
	e.GET("/", accessible)

	// Login route
	e.POST("/login", login)

	// Restricted group
	r := e.Group("/restricted")
	// r.Use(MiddlewareJWT)
	r.GET("", restricted)

	e.Logger.Fatal(e.Start(":1323"))
}
