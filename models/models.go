package models

type Credentials struct {
	User     string
	Password string
}

type CredentialsEntity struct {
	User              string
	EncryptedPassword string
	EncryptionKey     []byte
}
