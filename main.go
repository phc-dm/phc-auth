package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// UserType corrisponde alla "descrizione" dell'utente di ldap
type UserType int

// UserUID corrisponde all'uid dell'utente su ldap ed è una stringa unica che identifica l'utente
type UserUID string

// Token è un alias che rappresenta un token di accesso collegato ad una sessione
type Token string

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
	Token Token

	// Per ora pare che Ldap non supporti direttamente l'accesso attraverso digest md5,
	// bisogna vedere meglio come funge l'accesso con SASL con DIGEST-MD5
	Password string
	Username UserUID
}

// Service è l'intero servizio di autenticazione,
// contiene i dati per connettersi con ldap ed il server http
type Service struct {
	// LdapURL è l'url per il server di ldap
	LdapURL string
	// LdapBaseDN è il domino base di Ldap. Sotto questo dominio ci sono tutte le persone
	LdapBaseDN string

	server *http.Server

	sessionFromUsername map[UserUID]*UserSession
	sessionFromToken    map[Token]*UserSession
}

func (service *Service) newMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/status", service.statusHandler)
	mux.HandleFunc("/login", service.loginHandler)
	mux.HandleFunc("/q", service.queryHandler)
	mux.HandleFunc("/token", service.tokenHandler)

	return mux
}

func newAuthenticationService(addr, url, baseDN string) *Service {

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
	return service.server.ListenAndServe()
}

func (service *Service) statusHandler(res http.ResponseWriter, req *http.Request) {
	conn, err := service.NewLdapConnection()
	if err != nil {
		fmt.Fprint(res, false)
		return
	}
	defer conn.Close()

	fmt.Fprint(res, true)
}

func main() {
	service := newAuthenticationService(":5353", "ldaps://service", "ou=People,dc=phc,dc=unipi,dc=it")
	log.Fatal(service.ListenAndServe())
}
