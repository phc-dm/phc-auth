package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-ldap/ldap/v3"
)

// AuthenticationService è l'intero servizio di autenticazione,
// contiene la connessione con ldap ed il server http
type AuthenticationService struct {
	server *http.Server
	conn   *ldap.Conn
}

func newLdapConnection() *ldap.Conn {
	conn, err := ldap.DialURL("ldap://blabla.ldap.it")

	if err != nil {
		log.Fatal("Ldap connection error", err)
	}

	return conn
}

func newAuthenticationService(Addr string) *AuthenticationService {

	service := &AuthenticationService{}

	mux := http.NewServeMux()

	mux.HandleFunc("/status", service.statusHandler)
	mux.HandleFunc("/auth", service.authHandler)

	service.server = &http.Server{
		Handler:      mux,
		Addr:         Addr,
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}

	service.conn = newLdapConnection()

	return service
}

// ListenAndServe starts the server and returns if there are errors
func (service *AuthenticationService) ListenAndServe() error {
	return service.server.ListenAndServe()
}

// AuthRequest è una richiesta di autenticazione
type AuthRequest struct {
	Username, Password string
}

func (service *AuthenticationService) statusHandler(res http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(res, "online\n")
}

func (service *AuthenticationService) authHandler(res http.ResponseWriter, req *http.Request) {

	if req.Method != http.MethodPost {
		http.Error(res, "Only POST requests allowed", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		log.Fatal(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	var authRequest AuthRequest

	if err := json.Unmarshal([]byte(body), &authRequest); err != nil {
		log.Fatal(err)
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(res, false)

}

func main() {
	service := newAuthenticationService(":5353")
	log.Fatal(service.ListenAndServe())
}
