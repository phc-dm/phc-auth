package main

import (
	"fmt"

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

// Login takes a username and a password and does a ldap "bind" to check that the username and password are correct.
func (service *AuthenticationService) Login(username, password string) (*User, error) {

	usernameDN := fmt.Sprintf("uid=%s,%s", username, service.LdapBaseDN)

	conn, err := service.NewLdapConnection()

	if err != nil {
		return nil, err
	}

	bindErr := conn.Bind(usernameDN, password)
	if bindErr != nil {
		return nil, bindErr
	}

	return service.GetUser(username)
}

// GetUser gets informations about the asked user
func (service *AuthenticationService) GetUser(username string) (*User, error) {

	conn, err := service.NewLdapConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

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
