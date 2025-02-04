package models

// General credentials struct for cli use
type Credentials struct {
	Tag      string
	Username string
	Password string
}

// Struct for credentials entity in database
type CredentialsEntity struct {
	Tag      string
	Username string
	Password string
	Key_Data []byte
}
