package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/big"

	// "fmt"
	"io"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!#$%&'()*+,-./:;<=>?@[]^_`{|}~"

// GenerateRandomPassphrase generates a secure random passphrase of the given length.
func GenerateRandomPassphrase(length int) (string, error) {
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
	log.Printf("Generated random passkey of length %v: %v", length, string(passphrase))
	return string(passphrase), nil
}

func GenerateEncryptionKey() ([]byte, error) {
	log.Println("Generating encryption key...")
	key := make([]byte, 32) // AES-256 requires 32 bytes
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func Encrypt(text string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(encrypted string, key []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
