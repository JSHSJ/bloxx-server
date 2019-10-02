package users

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo"
)

type UserCreation struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Cart      []int  `json:"cart"`
}

var counter int

// Dir: Root directory for all data
var Dir string = "."

func SetupDatabase() {
	Dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadFile(Dir + "/counter")
	if err != nil {
		counter = 0
	} else {
		counter, err = strconv.Atoi(string(content))
		if err != nil {
			counter = 0
		}
	}

	os.MkdirAll(Dir+"/tokens", os.ModePerm)
	os.MkdirAll(Dir+"/users", os.ModePerm)
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func HandleCreateUser(c echo.Context) error {
	u := new(UserCreation)
	if err := c.Bind(u); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	token, err := createUser(*u, counter)

	if err != nil {
		if err.Error() == "email in use" {
			res := ErrorResponse{Message: "email in use"}
			return c.JSON(http.StatusBadRequest, res)
		}
		return c.NoContent(http.StatusInternalServerError)
	}

	increaseCounter()

	return c.JSON(http.StatusCreated, token)
}

func HandleUpdateUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.NoContent(http.StatusBadRequest)
	}

	userUpdate := new(User)
	if err := c.Bind(userUpdate); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	_, err = getUser(id)

	if err != nil {
		c.NoContent(http.StatusBadRequest)
	}

	err = saveUser(*userUpdate)

	if err != nil {
		c.NoContent(http.StatusBadRequest)
	}
	return c.NoContent(http.StatusOK)
}

type UserLogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func HandleLoginByEmail(c echo.Context) error {
	u := new(UserLogin)
	if err := c.Bind(u); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	user, err := loginUserByEmail(u.Email, u.Password)

	if err != nil {
		if err.Error() == "user doesn't exist" {
			return c.JSON(http.StatusNotFound, ErrorResponse{Message: err.Error()})
		}

		if err.Error() == "wrong password" {
			return c.JSON(http.StatusBadRequest, ErrorResponse{Message: err.Error()})

		}

		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, user)
}

func HandleLookupToken(c echo.Context) error {
	tokenString := c.Param("token")
	token, err := getToken(tokenString)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "Invalid Token"})
	}
	user, err := getUser(token.UserID)

	if err != nil {
		c.NoContent(http.StatusInternalServerError)
	}

	deleteToken(token)
	user.Active = true

	return c.JSON(http.StatusOK, user)
}

func CheckoutWithPaypal(c echo.Context) error {
	user := new(User)

	if err := c.Bind(user); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !userEmailExists(user.Email) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "user not in database"})
	}

	checkoutUser(*user)

	return c.JSON(http.StatusOK, user)
}

func CheckoutWithSofortUeberweisung(c echo.Context) error {
	user := new(User)

	if err := c.Bind(user); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !userEmailExists(user.Email) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "user not in database"})
	}

	checkoutUser(*user)

	return c.JSON(http.StatusOK, user)
}

func CheckoutWithBitcoin(c echo.Context) error {
	user := new(User)

	if err := c.Bind(user); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if !userEmailExists(user.Email) {
		return c.JSON(http.StatusBadRequest, ErrorResponse{Message: "user not in database"})
	}

	checkoutUser(*user)

	return c.JSON(http.StatusOK, user)
}

//Helpers

func increaseCounter() {
	counter++
	err := ioutil.WriteFile(Dir+"/counter", []byte(strconv.Itoa(counter)), 0644)
	if err != nil {
		fmt.Println("Can't update counter")
	}
}
