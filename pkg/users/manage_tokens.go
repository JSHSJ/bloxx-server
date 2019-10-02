package users

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
)

type Token struct {
	Token  string
	UserID int
}

func createToken(user User) (Token, error) {
	tokenString := ""

	for len(tokenString) < 10 {
		next := rand.Intn(9)
		tokenString += strconv.Itoa(next)
	}

	token := Token{
		Token:  tokenString,
		UserID: user.ID,
	}

	data, _ := json.MarshalIndent(token, "", " ")
	filename := fmt.Sprintf("%s.txt", token.Token)

	err := ioutil.WriteFile(Dir+"/tokens/"+filename, data, 0644)

	if err != nil {
		return Token{}, err
	}

	return token, nil
}

func getToken(tokenString string) (Token, error) {

	filename := fmt.Sprintf("%s.txt", tokenString)

	content, err := ioutil.ReadFile(Dir + "/tokens/" + filename)
	if err != nil {
		return Token{}, err
	}

	var token Token

	err = json.Unmarshal(content, &token)
	if err != nil {
		return Token{}, err
	}
	return token, nil
}

func deleteToken(token Token) {
	filename := fmt.Sprintf("%s.txt", token.Token)
	os.Remove(Dir + "/tokens/" + filename)
}
