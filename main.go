package main

import (
	"bloxxter/mockup-server/pkg/users"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {

	users.SetupDatabase()

	e := echo.New()
	// e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS default
	// Allows requests from any origin wth GET, HEAD, PUT, POST or DELETE method.
	e.Use(middleware.CORS())

	e.POST("/users", users.HandleCreateUser)
	e.PUT("/users/:id", users.HandleUpdateUser)
	e.GET("/tokens/:token", users.HandleLookupToken)

	e.POST("/login", users.HandleLoginByEmail)
	e.POST("/checkout/paypal", users.CheckoutWithPaypal)
	e.POST("/checkout/sofortueberweisung", users.CheckoutWithSofortUeberweisung)
	e.POST("/checkout/bitcoin", users.CheckoutWithBitcoin)

	e.Logger.Fatal(e.Start(":8000"))
}

//TOKEN
