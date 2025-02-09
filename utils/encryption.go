package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/zalando/go-keyring"
	"golang.org/x/crypto/argon2"
)

const service = "keeper_passwordmanager"

// Generate a random encryption key (only run once if no key exists)
func GenerateEncryptionKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", fmt.Errorf("failed to generate key: %w", err)
	}
	return base64.StdEncoding.EncodeToString(key), nil
}

// Store the encryption key in the OS keyring
func StoreKey() error {
	// Check if key already exists
	_, err := keyring.Get(service, "encryption_key")
	if err == nil {
		return nil // Key already stored
	}

	// Generate a new encryption key
	key, err := GenerateEncryptionKey()
	if err != nil {
		return err
	}

	// Store the key in the OS keyring
	return keyring.Set(service, "encryption_key", key)
}

// Retrieve encryption key from the OS keyring
func GetStoredKey() ([]byte, error) {
	keyStr, err := keyring.Get(service, "encryption_key")
	if err != nil {
		return nil, errors.New("encryption key not found, run setup first")
	}

	return base64.StdEncoding.DecodeString(keyStr)
}

// Encrypt password using AES-GCM.
func EncryptAES(plaintext string) (string, []byte, error) {
	// Get the stored encryption key
	key, err := GetStoredKey()
	if err != nil {
		return "", nil, err
	}

	// Generate a random salt
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", nil, err
	}

	// Derive encryption key
	derivedKey := argon2.IDKey(key, salt, 1, 64*1024, 4, 32)

	// Create AES cipher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", nil, err
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", nil, err
	}

	// Generate a random nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", nil, err
	}

	// Encrypt the password
	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)

	// Return encrypted data + salt (Base64 encoded)
	return base64.StdEncoding.EncodeToString(ciphertext), salt, nil
}

// DecryptAES decrypts the stored password.
func DecryptAES(ciphertextBase64 string, salt []byte) (string, error) {
	// Decode the base64 values
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	// Get the stored encryption key
	key, err := GetStoredKey()
	if err != nil {
		return "", err
	}

	// Derive encryption key
	derivedKey := argon2.IDKey(key, salt, 1, 64*1024, 4, 32)

	// Create AES cipher
	block, err := aes.NewCipher(derivedKey)
	if err != nil {
		return "", err
	}

	// Create GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Extract nonce
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the password
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// GenerateRandomPassphrase generates a secure random passphrase of the given length.
func GenerateRandomPassphrase(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'()*+,-./:;<=>?@[]^_`{|}~"
	if length < 12 {
		return "", fmt.Errorf("passphrase length should be at least 12 characters for security")
	}

	passphrase := make([]byte, length)
	for i := range passphrase {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		passphrase[i] = charset[randomIndex.Int64()]
	}
	log.Printf("Generated random passkey of length %v", length)
	return string(passphrase), nil
}

// func GenerateEncryptionKey() ([]byte, error) {
// 	log.Println("Generating encryption key...")
// 	key := make([]byte, 32) // AES-256 requires 32 bytes
// 	_, err := rand.Read(key)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return key, nil
// }

// func Encrypt(text string, key []byte) (string, error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return "", err
// 	}

// 	nonce := make([]byte, gcm.NonceSize())
// 	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
// 		return "", err
// 	}

// 	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
// 	return base64.StdEncoding.EncodeToString(ciphertext), nil
// }

// func Decrypt(encrypted string, key []byte) (string, error) {
// 	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
// 	if err != nil {
// 		return "", err
// 	}

// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		return "", err
// 	}

// 	gcm, err := cipher.NewGCM(block)
// 	if err != nil {
// 		return "", err
// 	}

// 	nonceSize := gcm.NonceSize()
// 	if len(ciphertext) < nonceSize {
// 		return "", errors.New("ciphertext too short")
// 	}

// 	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
// 	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
// 	if err != nil {
// 		return "", err
// 	}

// 	return string(plaintext), nil
// }

// func EncryptPassword(password string) (string, []byte, error) {
// 	// Encrypt new password and generate new key
// 	encryptionKey, err := GenerateEncryptionKey()
// 	if err != nil {
// 		return "", nil, fmt.Errorf("failed to generate encryption key: %w", err)
// 	}

// 	// Encrypt password with key
// 	encryptedPassword, err := Encrypt(password, encryptionKey)
// 	if err != nil {
// 		return "", nil, fmt.Errorf("failed to encrypt passphrase: %w", err)
// 	}

// 	return encryptedPassword, encryptionKey, nil
// }
