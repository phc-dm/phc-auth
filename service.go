package main

import (
	"log"
	"net/http"
	"time"

	// TODO: Per ora c'è questo perché "user.User" non è il top, pianificare un eventuale refactor
	. "github.com/phc-dm/auth-poisson/user"
)

// Service è l'intero servizio di autenticazione,
// contiene i dati per connettersi con ldap ed il server http
type Service struct {
	// LdapURL è l'url per il server di ldap
	LdapURL string
	// LdapBaseDN è il domino base di Ldap. Sotto questo dominio ci sono tutte le persone
	LdapBaseDN string

	server *http.Server

	sessionFromUsername map[UserUID]*Session
	sessionFromToken    map[Token]*Session
}

// CreateSession ...
func (service *Service) CreateSession(username UserUID, password string) Token {

	if session, ok := service.sessionFromUsername[username]; ok {
		log.Printf("Sending old token \"%s\" for user \"%s\"\n", session.Token, username)
		return session.Token
	}

	token := Token(GenerateRandomString(16))
	session := &Session{
		Username:  username,
		Password:  password,
		Token:     token,
		CreatedOn: time.Now(),
	}

	service.sessionFromUsername[username] = session
	service.sessionFromToken[token] = session

	log.Printf("Generated new token \"%s\" for user \"%s\"\n", session.Token, username)

	return token
}

// DestroySession ...
func (service *Service) DestroySession(session *Session) {
	delete(service.sessionFromUsername, session.Username)
	delete(service.sessionFromToken, session.Token)

	log.Printf("Destroied session for user \"%s\" with token \"%s\"\n", session.Username, session.Token)
}

func (service *Service) newMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", service.statusHandler)
	mux.HandleFunc("/login", service.loginHandler)
	mux.HandleFunc("/logout", service.logoutHandler)
	mux.HandleFunc("/q", service.queryHandler)

	mux.HandleFunc("/debug", service.debugHandler)

	return mux
}

// NewAuthenticationService ...
func NewAuthenticationService(addr, url, baseDN string) *Service {
	service := &Service{}
	service.LdapURL = url
	service.LdapBaseDN = baseDN

	service.sessionFromUsername = make(map[UserUID]*Session)
	service.sessionFromToken = make(map[Token]*Session)

	mux := service.newMux()
	service.server = &http.Server{
		Handler:      mux,
		Addr:         addr,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	return service
}

// ListenAndServe starts the server and returns if there are errors
func (service *Service) ListenAndServe() error {
	log.Printf("Starting server on address %s\n", service.server.Addr)
	return service.server.ListenAndServe()
}
