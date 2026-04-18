// Package crypt provides the tools to manage codes securely.
package crypt

import (
	"crypto/rand"
	"errors"
	
	"golang.org/x/crypto/chacha20poly1305"
)

// Encrypt plaintext with a key. Returns the nonce and ciphertext combined.
func Encrypt(key, plaintext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
    return nil, err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
    return nil, err
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)

	return append(nonce, ciphertext...), nil
}

// Decrypt. Separates the nonce from the ciphertext and returns the plaintext.
func Decrypt(key, data []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
    return nil, err
	}

	nonceSize := aead.NonceSize()
	if len(data) < nonceSize {
	    return nil, errors.New("crypt: data too short")
	}
	
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	 if err != nil {
	     return nil, err
	 }
	 return plaintext, nil 
}
