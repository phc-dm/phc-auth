package main

import (
	"log"
	"net/http"
	"time"
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

// UserUID corrisponde all'uid dell'utente su ldap ed è una stringa unica che identifica l'utente
type UserUID string

// Token è un alias che rappresenta un token di accesso collegato ad una sessione
type Token string

// Session ...
type Session struct {
	Token Token

	// Per ora pare che Ldap non supporti direttamente l'accesso attraverso digest md5,
	// bisogna vedere meglio come funge l'accesso con SASL con DIGEST-MD5
	Password string
	Username UserUID
}

// CreateSession ...
func (service *Service) CreateSession(username UserUID, password string) Token {
	token := Token(GenerateRandomString(16))
	session := &Session{
		Username: username,
		Password: password,
		Token:    token,
	}

	service.sessionFromUsername[username] = session
	service.sessionFromToken[token] = session

	return token
}

// DestroySession ...
func (service *Service) DestroySession(session *Session) {
	delete(service.sessionFromUsername, session.Username)
	delete(service.sessionFromToken, session.Token)
}

func (service *Service) newMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", service.statusHandler)
	mux.HandleFunc("/login", service.loginHandler)
	mux.HandleFunc("/logout", service.logoutHandler)
	mux.HandleFunc("/q", service.queryHandler)
	mux.HandleFunc("/token", service.tokenHandler)

	mux.HandleFunc("/debug", service.debugHandler)

	return mux
}

// NewAuthenticationService ...
func NewAuthenticationService(addr, url, baseDN string) *Service {
	service := &Service{}
	service.LdapURL = url
	service.LdapBaseDN = baseDN

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
