package main

import (
	"log"
	"net/http"
	"time"

	"github.com/phc-dm/auth-poisson/model"
)

// Token è un alias che rappresenta un token di accesso collegato ad una sessione
type UserSessionToken string

// Session rappresenta una sessione di un utente associata al servizio di autenticazione, è salvata
// in memoria quindi l'utente dovrà ottenere un nuovo token in caso di riavvio.
type Session struct {
	Token UserSessionToken

	// TODO: Per ora pare che LDAP non supporti direttamente l'accesso attraverso digest md5,
	// bisogna vedere meglio come funge l'accesso con SASL con DIGEST-MD5
	Password string
	Username model.UserUID

	CreatedOn time.Time
}

// Service è l'intero servizio di autenticazione,
// contiene i dati per connettersi con ldap ed il server http
type Service struct {
	// LdapURL è l'url per il server di ldap
	LdapURL string
	// LdapBaseDN è il domino base di Ldap. Sotto questo dominio ci sono tutte le persone
	LdapBaseDN string

	server *http.Server

	sessionFromUsername map[model.UserUID]*Session
	sessionFromToken    map[UserSessionToken]*Session
}

// CreateSession ...
func (service *Service) CreateSession(username model.UserUID, password string) UserSessionToken {

	if session, ok := service.sessionFromUsername[username]; ok {
		log.Printf("Sending old token \"%s\" for user \"%s\"\n", session.Token, username)
		return session.Token
	}

	token := UserSessionToken(GenerateRandomString(16))
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
	mux.HandleFunc("/users", service.queryHandler)

	mux.HandleFunc("/debug", service.debugHandler)

	return mux
}

// NewAuthenticationService ...
func NewAuthenticationService(addr, url, baseDN string) *Service {
	service := &Service{}
	service.LdapURL = url
	service.LdapBaseDN = baseDN

	service.sessionFromUsername = make(map[model.UserUID]*Session)
	service.sessionFromToken = make(map[UserSessionToken]*Session)

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
