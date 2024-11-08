// utils/encryption.go
package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"theransticslabs/m/config"
)

// getCipher creates an AEAD cipher using the provided base64-encoded key.
func getCipher(key string) (cipher.AEAD, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != 32 {
		return nil, errors.New("encryption key must be 32 bytes")
	}
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	return aead, nil
}

// Encrypt performs double encryption on the plaintext using two keys.
func Encrypt(plaintext string) (string, error) {
	// First encryption with EncryptionKey1
	aead1, err := getCipher(config.AppConfig.EncryptionKey1)
	if err != nil {
		return "", err
	}

	nonce1 := make([]byte, aead1.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce1); err != nil {
		return "", err
	}

	ciphertext1 := aead1.Seal(nil, nonce1, []byte(plaintext), nil)
	encrypted1 := append(nonce1, ciphertext1...)
	encrypted1B64 := base64.StdEncoding.EncodeToString(encrypted1)

	// Second encryption with EncryptionKey2
	aead2, err := getCipher(config.AppConfig.EncryptionKey2)
	if err != nil {
		return "", err
	}

	nonce2 := make([]byte, aead2.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce2); err != nil {
		return "", err
	}

	ciphertext2 := aead2.Seal(nil, nonce2, []byte(encrypted1B64), nil)
	encrypted2 := append(nonce2, ciphertext2...)
	encrypted2B64 := base64.StdEncoding.EncodeToString(encrypted2)

	return encrypted2B64, nil
}

// Decrypt reverses the double encryption process to retrieve the plaintext.
func Decrypt(ciphertext string) (string, error) {
	// Base64 decode
	encrypted2, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// First decryption with EncryptionKey2
	aead2, err := getCipher(config.AppConfig.EncryptionKey2)
	if err != nil {
		return "", err
	}

	if len(encrypted2) < aead2.NonceSize() {
		return "", errors.New("ciphertext too short for key2")
	}

	nonce2, ciphertext2 := encrypted2[:aead2.NonceSize()], encrypted2[aead2.NonceSize():]
	encrypted1B64, err := aead2.Open(nil, nonce2, ciphertext2, nil)
	if err != nil {
		return "", err
	}

	// Base64 decode again
	encrypted1, err := base64.StdEncoding.DecodeString(string(encrypted1B64))
	if err != nil {
		return "", err
	}

	// Second decryption with EncryptionKey1
	aead1, err := getCipher(config.AppConfig.EncryptionKey1)
	if err != nil {
		return "", err
	}

	if len(encrypted1) < aead1.NonceSize() {
		return "", errors.New("ciphertext too short for key1")
	}

	nonce1, ciphertext1 := encrypted1[:aead1.NonceSize()], encrypted1[aead1.NonceSize():]
	plaintextBytes, err := aead1.Open(nil, nonce1, ciphertext1, nil)
	if err != nil {
		return "", err
	}

	return string(plaintextBytes), nil
}
