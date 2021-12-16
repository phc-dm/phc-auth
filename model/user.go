package model

// UserUID corrisponde all'uid dell'utente su ldap ed è una stringa unica che identifica l'utente.
type UserUID string

// UserRole corrisponde alla "descrizione" dell'utente di ldap
type UserRole string

// Possibili ruoli scoperti fino ad ora
const (
	Studente   UserRole = "studente"
	Esterno    UserRole = "esterno"
	Dottorando UserRole = "dottorando"
)

// User is the a representation of an LDAP user with the most important properties
type User struct {
	Username UserUID `ldap:"uid" json:"username"`

	ID          int      `ldap:"uidNumber"   json:"id"`
	Name        string   `ldap:"givenName"   json:"name"`
	Surname     string   `ldap:"sn"          json:"surname"`
	Email       string   `ldap:"mail"        json:"email"`
	Description UserRole `ldap:"description" json:"description"`
	FullName    string   `ldap:"gecos"       json:"fullname"`
}
