package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func httpError(res http.ResponseWriter, err error) {
	http.Error(res, err.Error(), http.StatusInternalServerError)
	log.Println(err)
}

// LoginRequest è una richiesta di autenticazione
type LoginRequest struct {
	Username, Password string
}

func (service *AuthenticationService) loginHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		httpError(res, errors.New("Only POST requests allowed"))
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		httpError(res, err)
		return
	}

	var loginRequest LoginRequest

	if err := json.Unmarshal([]byte(body), &loginRequest); err != nil {
		httpError(res, err)
		return
	}

	user, err := service.Login(loginRequest.Username, loginRequest.Password)
	if err != nil {
		httpError(res, err)
		return
	}

	token := GenerateRandomString(16)
	service.sessions[user.Username] = &UserSession{
		User:  *user,
		Token: token,
	}

	log.Printf("Created new token \"%s\" for user @%s", token, user.Username)

	fmt.Fprint(res, token)
}

func (service *AuthenticationService) queryHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		httpError(res, errors.New("Only GET requests allowed"))
		return
	}

	username := req.FormValue("username")

	user, err := service.GetUser(username)

	if err != nil {
		http.Error(res, "User not found", http.StatusNotFound)
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		httpError(res, err)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.Write(userJSON)
}

// UpdateUserPropertyRequest ...
type UpdateUserPropertyRequest struct {
	Username, Token, Property, Value string
}

func (service *AuthenticationService) updateHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		httpError(res, errors.New("Only POST requests allowed"))
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		httpError(res, err)
		return
	}

	var updateRequest UpdateUserPropertyRequest

	if err := json.Unmarshal([]byte(body), &updateRequest); err != nil {
		httpError(res, err)
		return
	}

	session, ok := service.sessions[UserUID(updateRequest.Username)]
	if ok && session.Token != updateRequest.Token {
		http.Error(res, "Invalid token", http.StatusUnauthorized)
		return
	}

	// TODO: Do actual change of the prop

	fmt.Fprint(res, true)

}
