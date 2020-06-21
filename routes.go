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

// LoginRequest rappresenta una richiesta di autenticazione in JSON
type LoginRequest struct {
	Username, Password string
}

func (service *Service) createSession(username UserUID, password string) Token {
	token := Token(GenerateRandomString(16))
	session := &UserSession{
		Username: username,
		Password: password,
		Token:    token,
	}

	service.sessionFromUsername[username] = session
	service.sessionFromToken[token] = session

	return token
}

func (service *Service) loginHandler(res http.ResponseWriter, req *http.Request) {

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

	username := UserUID(loginRequest.Username)

	loginErr := service.CheckPassword(username, loginRequest.Password)
	if loginErr != nil {
		httpError(res, loginErr)
		return
	}

	token := service.createSession(username, loginRequest.Password)

	log.Printf("Created new token \"%s\" for user @%s", token, username)

	fmt.Fprint(res, token)
}

func (service *Service) queryHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		httpError(res, errors.New("Only GET requests allowed"))
		return
	}

	username := UserUID(req.FormValue("username"))

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

func (service *Service) tokenHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		httpError(res, errors.New("Only GET requests allowed"))
		return
	}

	username := UserUID(req.FormValue("username"))

	session, ok := service.sessionFromUsername[username]

	if !ok {
		http.Error(res, "User not found", http.StatusNotFound)
		return
	}

	fmt.Fprint(res, session.Token)
}

// UpdateUserPropertyRequest rappresenta una richiesta di cambio di attributo su ldap
type UpdateUserPropertyRequest struct {
	Token, Property, Value string
}

func (service *Service) updateHandler(res http.ResponseWriter, req *http.Request) {

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

	token := Token(updateRequest.Token)
	session, ok := service.sessionFromToken[token]
	if !ok {
		http.Error(res, "Invalid token", http.StatusUnauthorized)
		return
	}

	service.UpdateUserProperty(session.Username, session.Password, updateRequest.Property, updateRequest.Value)

	fmt.Fprint(res, true)

}
