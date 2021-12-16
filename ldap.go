package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-ldap/ldap/v3"
	"github.com/phc-dm/auth-poisson/model"
)

var attributesToRetrive = []string{
	"dn",
	"cn",

	"uid",
	"uidNumber",
	"givenName",
	"sn",
	"gecos",
	"mail",
	"description",
	"homeDirectory",
	"loginShell",
}

// NewLdapConnection ...
func (service *Service) NewLdapConnection() (*ldap.Conn, error) {
	return ldap.DialURL(service.LdapURL)
}

func (service *Service) loginWithConn(conn *ldap.Conn, username model.UserUID, password string) error {

	usernameDN := fmt.Sprintf("uid=%s,%s", username, service.LdapBaseDN)

	bindErr := conn.Bind(usernameDN, password)
	if bindErr != nil {
		return bindErr
	}

	return nil

}

// CheckPassword takes a username and a password and does a ldap "bind" to check that the username and password are correct.
func (service *Service) CheckPassword(username model.UserUID, password string) error {
	conn, err := service.NewLdapConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	return service.loginWithConn(conn, username, password)
}

// GetUser by creating a new connection and quering ldap
func (service *Service) GetUser(username model.UserUID) (*model.User, error) {
	conn, err := service.NewLdapConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return service.getUserWithConn(conn, username)
}

// getUserWithConn gets informations about the asked user using a given ldap connection
func (service *Service) getUserWithConn(conn *ldap.Conn, username model.UserUID) (*model.User, error) {

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
		return nil, fmt.Errorf("invalid number of entries from LDAP, got %d", len(searchResult.Entries))
	}

	entry := searchResult.Entries[0]

	uidNumber, err := strconv.ParseInt(entry.GetAttributeValue("uidNumber"), 10, 32)
	if err != nil {
		return nil, err
	}

	description := model.UserRole(entry.GetAttributeValue("description"))

	return &model.User{
		Username: model.UserUID(entry.GetAttributeValue("uid")),

		ID:          int(uidNumber),
		Name:        entry.GetAttributeValue("givenName"),
		Surname:     entry.GetAttributeValue("sn"),
		FullName:    entry.GetAttributeValue("gecos"),
		Email:       entry.GetAttributeValue("mail"),
		Description: description,
	}, nil
}

// GetUsers by creating a new connection and quering ldap
func (service *Service) GetUsers() ([]*model.User, error) {
	conn, err := service.NewLdapConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return service.getUsersWithConn(conn)
}

// getUsersWithConn gets informations about the asked user using a given ldap connection
func (service *Service) getUsersWithConn(conn *ldap.Conn) ([]*model.User, error) {

	searchRequest := ldap.NewSearchRequest(
		service.LdapBaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(uid=*)",
		attributesToRetrive,
		nil,
	)

	searchResult, err := conn.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	users := []*model.User{}

	for _, entry := range searchResult.Entries {

		uidNumber, err := strconv.ParseUint(entry.GetAttributeValue("uidNumber"), 10, 32)
		if err != nil {
			return nil, err
		}

		description := model.UserRole(entry.GetAttributeValue("description"))

		users = append(users, &model.User{
			Username: model.UserUID(entry.GetAttributeValue("uid")),

			ID:          int(uidNumber),
			Name:        entry.GetAttributeValue("givenName"),
			Surname:     entry.GetAttributeValue("sn"),
			FullName:    entry.GetAttributeValue("gecos"),
			Email:       entry.GetAttributeValue("mail"),
			Description: description,
		})
	}

	return users, nil
}

// UpdateUserProperty ...
func (service *Service) UpdateUserProperty(username model.UserUID, password, property, newValue string) error {
	conn, err := service.NewLdapConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	return service.updateUserPropertyWithConn(conn, username, password, property, newValue)
}

func (service *Service) updateUserPropertyWithConn(conn *ldap.Conn, username model.UserUID, password, property, newValue string) error {

	loginErr := service.loginWithConn(conn, username, password)
	if loginErr != nil {
		return loginErr
	}

	// TODO: Do ldap update property
	log.Fatal("Not implemented")

	return nil
}

// UpdateUserPassword ...
func (service *Service) UpdateUserPassword(username model.UserUID, password, newPassword string) error {
	conn, err := service.NewLdapConnection()
	if err != nil {
		return err
	}
	defer conn.Close()

	return service.updateUserPasswordWithConn(conn, username, password, newPassword)
}

func (service *Service) updateUserPasswordWithConn(conn *ldap.Conn, username model.UserUID, password, newPassword string) error {

	loginErr := service.loginWithConn(conn, username, password)
	if loginErr != nil {
		return loginErr
	}

	// TODO: Do ldap update password
	log.Fatal("Not implemented")

	return nil
}
