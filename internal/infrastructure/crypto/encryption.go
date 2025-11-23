package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// Encrypt encrypts plaintext using AES-256-GCM
// Returns: nonce + ciphertext (with auth tag appended by GCM)
func Encrypt(plaintext, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: expected 32 bytes, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a random nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and authenticate
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// Return nonce + ciphertext (ciphertext already includes auth tag)
	result := make([]byte, len(nonce)+len(ciphertext))
	copy(result, nonce)
	copy(result[len(nonce):], ciphertext)

	return result, nil
}

// Decrypt decrypts ciphertext using AES-256-GCM
// Input format: nonce + ciphertext (with auth tag)
func Decrypt(encrypted, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: expected 32 bytes, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce := encrypted[:nonceSize]
	ciphertext := encrypted[nonceSize:]

	// Decrypt and verify authentication
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}

// ZeroBytes securely zeros out a byte slice
func ZeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
