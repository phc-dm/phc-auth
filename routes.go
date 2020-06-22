package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

	token := service.CreateSession(username, loginRequest.Password)

	log.Printf("Created new token \"%s\" for user @%s\n", token, username)

	fmt.Fprint(res, token)
}

// LogoutRequest ...
type LogoutRequest struct {
	Token *string `json:"token,omitempty"`
}

func (service *Service) logoutHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		httpError(res, errors.New("Only POST requests allowed"))
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		httpError(res, err)
		return
	}

	var logoutRequest LogoutRequest

	if err := json.Unmarshal([]byte(body), &logoutRequest); err != nil {
		httpError(res, err)
		return
	}

	if logoutRequest.Token != nil {
		token := Token(*logoutRequest.Token)
		service.DestroySession(service.sessionFromToken[token])
	} else {
		httpError(res, errors.New("Invalid logout request, at least a token or a username field should be provided"))
		return
	}

	fmt.Fprint(res, true)
}

func (service *Service) queryHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		httpError(res, errors.New("Only GET requests allowed"))
		return
	}

	if req.FormValue("username") != "" {
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
	} else {
		users, err := service.GetUsers()

		if err != nil {
			http.Error(res, "User not found", http.StatusNotFound)
			return
		}

		usersJSON, err := json.Marshal(users)
		if err != nil {
			httpError(res, err)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.Write(usersJSON)
	}
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

func (service *Service) statusHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		httpError(res, errors.New("Only GET requests allowed"))
		return
	}

	conn, err := service.NewLdapConnection()
	if err != nil {
		fmt.Fprint(res, false)
		return
	}
	defer conn.Close()

	fmt.Fprint(res, true)
}

func (service *Service) debugHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodGet {
		httpError(res, errors.New("Only GET requests allowed"))
		return
	}

	log.Printf("Currently %d stored sessions\n", len(service.sessionFromUsername))
	for username, session := range service.sessionFromUsername {

		h := md5.New()
		io.WriteString(h, session.Password)

		logSession := *session
		logSession.Password = fmt.Sprintf("%x", h.Sum(nil))

		log.Printf("Session for @%s %+v\n", username, logSession)
	}

	fmt.Fprint(res, "Logged service information")
}
