package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"

	// "fmt"
	"io"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@%&?"

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

	return string(passphrase), nil
}

func GenerateRandomKey() ([]byte, error) {
	log.Println("Generating encryption key...")
	key := make([]byte, 32) // AES-256 requires 32 bytes
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// encrypt takes a plaintext string and encryption key, returns the encrypted string
// and any errors that might occur during encryption
func Encrypt(text string, key []byte) (string, error) {
	// Create a new AES cipher block using our encryption key
	// This is like creating our encryption machine with our specific key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Create a slice of bytes that will hold:
	// 1. The IV (Initialization Vector) at the beginning
	// 2. The encrypted text after that
	// aes.BlockSize is 16 bytes - that's the size of our IV
	ciphertext := make([]byte, aes.BlockSize+len(text))

	// Get a slice that points to just the IV portion (first 16 bytes)
	iv := ciphertext[:aes.BlockSize]

	// Fill the IV with random bytes
	// This makes sure each encryption is unique, even of the same text
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Create an encryption stream using CFB mode
	// CFB lets us encrypt data of any length, not just 16-byte blocks
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt our text and store it right after the IV
	// We're converting the string to []byte to encrypt it
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(text))

	// Convert everything to base64 so it can be safely stored/transmitted as text
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decrypt takes a base64 encoded encrypted string and the key used to encrypt it
// returns the original text and any errors that occur during decryption
func Decrypt(encrypted string, key []byte) (string, error) {
	// Convert the base64 encoded string back to bytes
	ciphertext, _ := base64.StdEncoding.DecodeString(encrypted)

	// Create the same type of AES cipher as we used for encryption
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Get the IV from the start of the ciphertext
	iv := ciphertext[:aes.BlockSize]

	// Get the actual encrypted data (everything after the IV)
	ciphertext = ciphertext[aes.BlockSize:]

	// Create a decryption stream using the same IV and key
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt the data in place
	// In CFB mode, the XOR operation works the same for encryption and decryption
	stream.XORKeyStream(ciphertext, ciphertext)

	// Convert the decrypted bytes back to a string
	return string(ciphertext), nil
}
