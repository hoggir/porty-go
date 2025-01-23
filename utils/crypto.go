package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt encrypts the given text using AES encryption
func Encrypt(plainText, key string) (string, error) {
	if len(key) != 16 {
		return "", errors.New("key must be 16 bytes long")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	plainTextBytes := []byte(plainText)

	// Generate a new initialization vector (IV)
	cipherText := make([]byte, aes.BlockSize+len(plainTextBytes))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plainTextBytes)

	// Return the encrypted data as a base64 encoded string
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts the given text using AES encryption
func Decrypt(encryptedText, key string) (string, error) {
	if len(key) != 16 {
		return "", errors.New("key must be 16 bytes long")
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	// Decode the base64 encoded string
	cipherText, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
