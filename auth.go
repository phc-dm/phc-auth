package main

import (
	"fmt"
	"log"

	"github.com/go-ldap/ldap/v3"
)

var attributesToRetrive = []string{
	"dn",
	"cn",
	"sn",
	"givenName",
	"mail",
	"uid",
	"homeDirectory",
	"loginShell",
	"gecos",
}

// NewLdapConnection ...
func (service *Service) NewLdapConnection() (*ldap.Conn, error) {
	return ldap.DialURL(service.LdapURL)
}

func (service *Service) loginWithConn(conn *ldap.Conn, username UserUID, password string) error {

	usernameDN := fmt.Sprintf("uid=%s,%s", username, service.LdapBaseDN)

	bindErr := conn.Bind(usernameDN, password)
	if bindErr != nil {
		return bindErr
	}

	return nil

}

// CheckPassword takes a username and a password and does a ldap "bind" to check that the username and password are correct.
func (service *Service) CheckPassword(username UserUID, password string) error {
	conn, err := service.NewLdapConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	return service.loginWithConn(conn, username, password)
}

// GetUser by creating a new connection and quering ldap
func (service *Service) GetUser(username UserUID) (*User, error) {
	conn, err := service.NewLdapConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return service.getUserWithConn(conn, username)
}

// getUserWithConn gets informations about the asked user using a given ldap connection
func (service *Service) getUserWithConn(conn *ldap.Conn, username UserUID) (*User, error) {

	ldapFilter := fmt.Sprintf("(uid=%s)", username)
	searchRequest := ldap.NewSearchRequest(
		service.LdapBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		ldapFilter,
		attributesToRetrive,
		nil,
	)

	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(searchResult.Entries) != 1 {
		return nil, fmt.Errorf("Invalid number of entries from LDAP, got %d", len(searchResult.Entries))
	}

	entry := searchResult.Entries[0]

	description := userTypeDescriptionMap[entry.GetAttributeValue("description")]
	return &User{
		Name:        entry.GetAttributeValue("givenName"),
		Surname:     entry.GetAttributeValue("sn"),
		FullName:    entry.GetAttributeValue("gecos"),
		Email:       entry.GetAttributeValue("mail"),
		Description: description,
	}, nil
}

// UpdateUserProperty ...
func (service *Service) UpdateUserProperty(username UserUID, password, property, newValue string) error {
	conn, err := service.NewLdapConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	return service.updateUserPropertyWithConn(conn, username, password, property, newValue)
}

func (service *Service) updateUserPropertyWithConn(conn *ldap.Conn, username UserUID, password, property, newValue string) error {

	loginErr := service.loginWithConn(conn, username, password)
	if loginErr != nil {
		return loginErr
	}

	// TODO: Do ldap update property
	log.Fatal("Not implemented")

	return nil
}

// UpdateUserPassword ...
func (service *Service) UpdateUserPassword(username UserUID, password, newPassword string) error {
	conn, err := service.NewLdapConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	return service.updateUserPasswordWithConn(conn, username, password, newPassword)
}

func (service *Service) updateUserPasswordWithConn(conn *ldap.Conn, username UserUID, password, newPassword string) error {

	loginErr := service.loginWithConn(conn, username, password)
	if loginErr != nil {
		return loginErr
	}

	// TODO: Do ldap update password
	log.Fatal("Not implemented")

	return nil
}
