package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// UserType corrisponde alla "descrizione" dell'utente di ldap
type UserType int

// UserUID corrisponde all'uid dell'utente su ldap ed è una stringa unica che identifica l'utente
type UserUID string

const (
	_ UserType = iota
	// Studente su ldap è `studente`
	Studente
	// Esterno su ldap è `esterno`
	Esterno
	// Dottorando su ldap è `dottorando`
	Dottorando

	// Unknown per quando il campo è assente o non riconosciuto
	Unknown
)

var userTypeDescriptionMap = map[string]UserType{
	"studente":   Studente,
	"esterno":    Esterno,
	"dottornado": Dottorando,
}

// User ...
type User struct {
	Username UserUID

	ID          int
	Name        string
	Surname     string
	Email       string
	Description UserType

	// On ldap this is gecos
	FullName string
}

// UserSession ...
type UserSession struct {
	User  User
	Token string
}

// AuthenticationService è l'intero servizio di autenticazione,
// contiene la connessione con ldap ed il server http
type AuthenticationService struct {
	// LdapURL è l'url per il server di ldap
	LdapURL string
	// LdapBaseDN è il domino base di Ldap. Sotto questo dominio ci sono tutte le persone
	LdapBaseDN string

	server   *http.Server
	sessions map[UserUID]*UserSession
}

// NewLdapConnection ...
func (service *AuthenticationService) NewLdapConnection() (*ldap.Conn, error) {
	return ldap.DialURL(service.LdapURL)
}

func (service *AuthenticationService) newMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", service.statusHandler)
	mux.HandleFunc("/login", service.loginHandler)
	mux.HandleFunc("/q", service.queryHandler)

	return mux
}

func newAuthenticationService(addr, url, baseDN string) *AuthenticationService {

	service := &AuthenticationService{}
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
func (service *AuthenticationService) ListenAndServe() error {
	return service.server.ListenAndServe()
}

func (service *AuthenticationService) statusHandler(res http.ResponseWriter, req *http.Request) {

	fmt.Fprint(res, true)
}

func main() {
	service := newAuthenticationService(":5353", "ldaps://service", "ou=People,dc=phc,dc=unipi,dc=it")
	log.Fatal(service.ListenAndServe())
}
