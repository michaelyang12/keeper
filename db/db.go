package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/michaelyang12/keeper/models"
	"github.com/michaelyang12/keeper/utils"
	_ "github.com/mutecomm/go-sqlcipher/v4"
)

var SqlDb *sql.DB

func GenerateEncryptionKey() ([]byte, error) {
	key, err := utils.GenerateRandomKey()
	if err != nil {
		return nil, fmt.Errorf("error generating encryption key: %w", err)
	}
	log.Printf("Encryption Key: %v\n", key)

	return key, nil
}

// TODO: Come back to this and implement passphrase
func InitializeLocalDatabase(passphrase string) error {

	// if passphrase == "" {
	// 	var err error
	// 	passphrase, err = utils.GenerateRandomPassphrase(12)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to generate random passphrase: %w", err)
	// 	}
	// }
	// os.Remove("credentials.db")

	// OPEN DATABASE
	// connectionString := fmt.Sprintf("file:credentials.db?_pragma_key=test&_pragma_cipher=AES-256-CBC&_pragma_cipher_page_size=4096", passphrase)
	connectionString := "file:credentials.db?_pragma_key=&_pragma_cipher=AES-256-CBC&_pragma_cipher_page_size=4096"
	// db, err := sql.Open("sqlite3", fmt.Sprintf("file:credentials.db?_pragma_key=%s", passphrase))
	db, err := sql.Open("sqlite3", connectionString)
	if err != nil {
		return fmt.Errorf("open error: %v", err)
	}
	// defer db.Close()
	// Create credentials table
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS credentials (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        tag TEXT NOT NULL,
        username TEXT NOT NULL,
        password TEXT NOT NULL,
		key_data BLOB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
	if _, err := db.Exec(createTableQuery); err != nil {
		return fmt.Errorf("failed to create credentials table: %w", err)
	}
	SqlDb = db
	return nil
}

func InsertNewCredential(tag string, user string, password string) error {
	// Generate encryption key to encrypt password
	encryptionKey, err := GenerateEncryptionKey()
	if err != nil {
		return fmt.Errorf("failed to generate encryption key: %w", err)
	}

	// Encrypt password with key
	encryptedPassword, err := utils.Encrypt(password, encryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt passphrase: %w", err)
	}

	// Insert initial record, with encrypted password and encryption key used
	insertQuery := `INSERT OR REPLACE INTO credentials (tag, username, password, key_data) VALUES (?, ?, ?, ?)`
	if _, err := SqlDb.Exec(insertQuery, tag, user, encryptedPassword, encryptionKey); err != nil {
		return fmt.Errorf("failed to insert initial record: %w", err)
	}

	return nil
}

func FetchExistingCredential(tag string) (*models.Credentials, error) {
	fetchQuery := `SELECT username, password, key_data FROM credentials WHERE tag = ?`
	var entity models.CredentialsEntity
	row := SqlDb.QueryRow(fetchQuery, tag)

	err := row.Scan(&entity.User, &entity.EncryptedPassword, &entity.EncryptionKey)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no credential found with tag: %s", tag)
		}
		return nil, fmt.Errorf("error retrieving credential: %w", err)
	}

	password, err := utils.Decrypt(entity.EncryptedPassword, entity.EncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error decrypting password: %w", err)
	}

	cred := models.Credentials{
		User:     entity.User,
		Password: password,
	}

	return &cred, nil
}
