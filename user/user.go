package user

import "time"

// TODO: Eseguire un refactor della struct "User"
// user.UserUID 	-> user.UID
// user.UserRole 	-> user.Role
// user.User 		-> user.Repr (?)

// UserUID corrisponde all'uid dell'utente su ldap ed è una stringa unica che identifica l'utente.
type UserUID string

// UserRole corrisponde alla "descrizione" dell'utente di ldap
//  TODO: Verrà rinominato a user.Role
type UserRole string

// Descrizioni scoperte per ora
const (
	Studente   UserRole = "studente"
	Esterno    UserRole = "esterno"
	Dottorando UserRole = "dottorando"
)

// Token è un alias che rappresenta un token di accesso collegato ad una sessione
type Token string

// User ...
//  TODO: Verrà rinominato a ???, proposta: user.Repr (da Representation, in quanto
//  è la rappresentazione di un utente di LDAP in memoria), user.Authenticated, user.Info,
//  user.Ref, user.Instance, ...
type User struct {
	Username UserUID `ldap:"uid" json:"username"`

	ID          int      `ldap:"uidNumber"   json:"id"`
	Name        string   `ldap:"givenName"   json:"name"`
	Surname     string   `ldap:"sn"          json:"surname"`
	Email       string   `ldap:"mail"        json:"email"`
	Description UserRole `ldap:"description" json:"description"`
	FullName    string   `ldap:"gecos"       json:"fullname"`
}

// Session rappresenta una sessione di un utente associata al servizio di autenticazione, è salvata
// in memoria quindi l'utente dovrà ottenere un nuovo token in caso di riavvio.
type Session struct {
	Token Token

	// TODO: Per ora pare che Ldap non supporti direttamente l'accesso attraverso digest md5,
	// bisogna vedere meglio come funge l'accesso con SASL con DIGEST-MD5
	Password string
	Username UserUID

	CreatedOn time.Time
}
