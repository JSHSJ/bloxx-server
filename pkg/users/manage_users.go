package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type User struct {
	ID               int    `json:"id"`
	FirstName        string `json:"firstname"`
	LastName         string `json:"lastname"`
	Email            string `json:"email"`
	Password         string `json:"password"`
	Active           bool   `json:"active"`
	Cart             []int  `json:"cart"`
	Street           string `json:"street"`
	City             string `json:"city"`
	ZipCode          string `json:"zip_code"`
	Country          string `json:"counter"`
	PreferredPayment string `json:"preferred_payment"`
}

func createUser(uc UserCreation, counter int) (Token, error) {
	if userEmailExists(uc.Email) {
		err := errors.New("email in use")
		return Token{}, err
	}

	u := User{
		ID:               counter,
		FirstName:        uc.FirstName,
		LastName:         uc.LastName,
		Email:            uc.Email,
		Password:         uc.Password,
		Active:           false,
		Cart:             uc.Cart,
		Street:           "",
		City:             "",
		ZipCode:          "",
		Country:          "",
		PreferredPayment: "",
	}

	err := saveUser(u)
	if err != nil {
		return Token{}, err
	}

	return createToken(u)
}

func saveUser(user User) error {
	data, _ := json.MarshalIndent(user, "", " ")
	filename := fmt.Sprintf("%d.txt", user.ID)

	return ioutil.WriteFile(Dir+"/users/"+filename, data, 0644)
}

func getUser(id int) (User, error) {

	filename := fmt.Sprintf("%d.txt", id)

	content, err := ioutil.ReadFile("./users/" + filename)
	if err != nil {
		return User{}, err
	}

	var user User

	err = json.Unmarshal(content, &user)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func loginUserByEmail(email, password string) (User, error) {
	user, err := getUserByMail(email)

	if err != nil {
		err = errors.New("user doesn't exist")
		return User{}, err
	}

	if user.Password == password {
		if user.Active == false {
			return user, errors.New("email not activated")
		}
		return user, nil
	}
	return User{}, errors.New("wrong password")
}

func checkoutUser(user User) {
	user.Cart = user.Cart[:0]
}

func userEmailExists(email string) bool {
	_, err := getUserByMail(email)
	if err != nil {
		return false
	}

	return true
}

func getUserByMail(email string) (User, error) {
	files, err := ioutil.ReadDir(Dir + "/users")
	if err != nil {
		return User{}, err
	}

	for _, file := range files {
		content, err := ioutil.ReadFile(Dir + "/users/" + file.Name())
		if err != nil {
			return User{}, err
		}

		var user User

		err = json.Unmarshal(content, &user)
		if err != nil {
			return User{}, err
		}

		if user.Email == email {
			return user, nil
		}
	}

	return User{}, errors.New("user doesn't exist")
}
