package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/gorilla/mux"
)

// AuthenticationService è l'intero servizio di autenticazione,
// contiene la connessione con ldap ed il server http
type AuthenticationService struct {
	server *http.Server
	conn   *ldap.Conn
}

// AuthRequest è una richiesta di autenticazione
type AuthRequest struct {
	Username, Password string
}

func (service *AuthenticationService) statusHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "online\n")
}

func (service *AuthenticationService) authHandler(w http.ResponseWriter, req *http.Request) {

	var authRequest AuthRequest

	body, _ := ioutil.ReadAll(req.Body)

	err := json.Unmarshal([]byte(body), &authRequest)

	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, false)

}

func main() {

	service := &AuthenticationService{}

	r := mux.NewRouter()
	r.HandleFunc("/status", service.statusHandler)
	r.HandleFunc("/auth", service.authHandler).Methods("POST")

	// conn, err := ldap.DialURL("ldap://blabla.ldap.it")

	// if err != nil {
	// 	log.Fatal("Ldap connection error", err)
	// }

	service.server = &http.Server{
		Handler:      r,
		Addr:         "localhost:5353",
		WriteTimeout: 1 * time.Second,
		ReadTimeout:  1 * time.Second,
	}
	// service.conn = conn

	log.Fatal(service.server.ListenAndServe())
}
