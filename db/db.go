package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/michaelyang12/keeper/models"
	"github.com/michaelyang12/keeper/utils"
	_ "github.com/mutecomm/go-sqlcipher/v4"
)

var SqlDb *sql.DB

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// TODO: Come back to this and implement passphrase
func InitializeLocalDatabase() error {

	// Check if database exists
	isNewDB := !fileExists("credentials.db")
	// var passphrase string
	// if isNewDB {
	// 	var err error
	// 	passphrase, err = utils.GenerateRandomPassphrase(16)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to generate random passphrase: %w", err)
	// 	}
	// } else {
	// 	// Get passphrase from local file storage
	// }
	// // Open database with encryption settings in connection string
	// connectionString := fmt.Sprintf("file:credentials.db?_pragma_key='%s'&_pragma_cipher_page_size=4096", passphrase)
	// db, err := sql.Open("sqlite3", connectionString)
	// if err != nil {
	// 	return fmt.Errorf("open error: %v", err)
	// }

	// // For new database, create the table
	// if isNewDB {
	// 	createTableQuery := `
	//     CREATE TABLE IF NOT EXISTS credentials (
	//         id INTEGER PRIMARY KEY AUTOINCREMENT,
	//         tag TEXT NOT NULL,
	//         username TEXT NOT NULL,
	//         password TEXT NOT NULL,
	//         key_data BLOB NOT NULL,
	//         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	//     )`
	// 	if _, err := db.Exec(createTableQuery); err != nil {
	// 		return fmt.Errorf("failed to create credentials table: %w", err)
	// 	}
	// }

	// // Set other PRAGMA statements
	// pragmas := []string{
	// 	"PRAGMA cipher_page_size = 4096",
	// 	"PRAGMA kdf_iter = 64000",
	// 	"PRAGMA cipher_memory_security = ON",
	// 	"PRAGMA cipher_default_kdf_iter = 64000",
	// 	"PRAGMA cipher_use_hmac = ON",
	// }

	// for _, pragma := range pragmas {
	// 	if _, err := db.Exec(pragma); err != nil {
	// 		return fmt.Errorf("failed to set pragma: %w", err)
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
	if isNewDB {
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
		log.Println("Initialized local sqlcipher database with table 'credentials'")
	}

	SqlDb = db
	return nil
}

func InsertNewCredentials(tag string, user string, password string) error {
	// Generate encrypted password and salt
	encryptedPassword, salt, err := utils.EncryptAES(password)
	if err != nil {
		return fmt.Errorf("failed to encrypt password: %w", err)
	}

	// Insert encrypted credentials into db
	query := `INSERT OR REPLACE INTO credentials (tag, username, password, salt) VALUES (?, ?, ?, ?)`
	if _, err := SqlDb.Exec(query, tag, user, encryptedPassword, salt); err != nil {
		return fmt.Errorf("failed to insert initial record: %w", err)
	}
	return nil
}

func FetchExistingCredentials(tag string) (*models.Credentials, error) {
	// Get credentials from db
	query := `SELECT tag, username, password, salt FROM credentials WHERE tag = ?`
	var entity models.CredentialsEntity
	err := SqlDb.QueryRow(query, tag).Scan(&entity.Tag, &entity.Username, &entity.Password, &entity.Salt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no credentials found with tag: %s", tag)
		}
		return nil, fmt.Errorf("error retrieving credentials: %w", err)
	}

	// Decrypt password
	password, err := utils.DecryptAES(entity.Password, entity.Salt)
	if err != nil {
		return nil, fmt.Errorf("error decrypting password: %w", err)
	}

	cred := models.Credentials{
		Tag:      tag,
		Username: entity.Username,
		Password: password,
	}

	return &cred, nil
}

// func InsertNewCredential(tag string, user string, password string) error {
// 	// Generate encryption key to encrypt password
// 	encryptedPassword, encryptionKey, err := utils.EncryptPassword(password)
// 	if err != nil {
// 		return err
// 	}

// 	// Insert initial record, with encrypted password and encryption key used
// 	insertQuery := `INSERT OR REPLACE INTO credentials (tag, username, password, key_data) VALUES (?, ?, ?, ?)`
// 	if _, err := SqlDb.Exec(insertQuery, tag, user, encryptedPassword, encryptionKey); err != nil {
// 		return fmt.Errorf("failed to insert initial record: %w", err)
// 	}

// 	return nil
// }

// func FetchExistingCredential(tag string) (*models.Credentials, error) {
// 	fetchQuery := `SELECT tag, username, password, key_data FROM credentials WHERE tag = ?`
// 	var entity models.CredentialsEntity
// 	row := SqlDb.QueryRow(fetchQuery, tag)

// 	err := row.Scan(&entity.Tag, &entity.Username, &entity.Password, &entity.Key_Data)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			return nil, fmt.Errorf("no credentials found with tag: %s", tag)
// 		}
// 		return nil, fmt.Errorf("error retrieving credentials: %w", err)
// 	}

// 	password, err := utils.Decrypt(entity.Password, entity.Key_Data)
// 	if err != nil {
// 		return nil, fmt.Errorf("error decrypting password: %w", err)
// 	}

// 	cred := models.Credentials{
// 		Tag:      tag,
// 		Username: entity.Username,
// 		Password: password,
// 	}

// 	return &cred, nil
// }

// func DeleteExistingCredential(tag string) error {
// 	deleteQuery := `DELETE FROM credentials WHERE tag = ?`
// 	result, err := SqlDb.Exec(deleteQuery, tag)
// 	if err != nil {
// 		return fmt.Errorf("failed to delete record: %w", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("error getting affected rows: %w", err)
// 	}

// 	if rowsAffected > 0 {
// 		logging.Success("Removed credentials with tag: %v\n", tag)
// 		return nil
// 	}
// 	logging.Warn("Credentials with specified tag doesn't exist. Nothing to delete.\n")
// 	return nil
// }

// func UpdateExistingCredential(tag string, username string, password string) error {
// 	updateQuery := `UPDATE credentials
//     SET username = ?,
//         password = ?,
//         key_data = ?
//     WHERE tag = ?`

// 	encryptedPassword, encryptionKey, err := utils.EncryptPassword(password)
// 	if err != nil {
// 		return err
// 	}

// 	result, err := SqlDb.Exec(updateQuery, username, encryptedPassword, encryptionKey, tag)
// 	if err != nil {
// 		return fmt.Errorf("failed to update record: %w", err)
// 	}

// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return fmt.Errorf("error getting affected rows: %w", err)
// 	}

// 	if rowsAffected > 0 {
// 		logging.Success("Updated credentials with tag: %v\n", tag)
// 		return nil
// 	}
// 	logging.Warn("Credentials with specified tag doesn't exist. Nothing to update.\n")
// 	return nil
// }

func FetchAllExistingCredentials() ([]models.Credentials, error) {
	fetchQuery := `SELECT tag, username, password, key_data FROM credentials`
	var creds []models.Credentials

	rows, err := SqlDb.Query(fetchQuery)
	if err != nil {
		return nil, fmt.Errorf("error fetching rows from db: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entity models.CredentialsEntity

		// Scan the columns into the entity struct
		if err := rows.Scan(&entity.Tag, &entity.Username, &entity.Password, &entity.Key_Data); err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}

		// Get decrypted password
		password, err := utils.Decrypt(entity.Password, entity.Key_Data)
		if err != nil {
			return nil, fmt.Errorf("error decrypting password: %w", err)
		}

		// Get credentials and append to list
		cred := models.Credentials{
			Tag:      entity.Tag,
			Username: entity.Username,
			Password: password,
		}

		creds = append(creds, cred)
	}

	return creds, nil
}
