package main

import (
	"log"
)

func main() {
	service := NewAuthenticationService(":5353", "ldaps://service", "ou=People,dc=phc,dc=unipi,dc=it")
	log.Fatal(service.ListenAndServe())
}
