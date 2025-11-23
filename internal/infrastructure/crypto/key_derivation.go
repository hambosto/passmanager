package crypto

import (
	"crypto/rand"
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/argon2"
)

// KeyDerivationParams contains parameters for Argon2id key derivation
type KeyDerivationParams struct {
	Algorithm   string `json:"algorithm"`   // "argon2id"
	Iterations  uint32 `json:"iterations"`  // Time parameter (3-4 recommended)
	Memory      uint32 `json:"memory"`      // Memory in KB (64MB = 65536)
	Parallelism uint8  `json:"parallelism"` // Threads (4 recommended)
	Salt        []byte `json:"salt"`        // 32 bytes random salt
	KeyLength   uint32 `json:"key_length"`  // Output key length (32 bytes for AES-256)
}

// DefaultKeyDerivationParams returns secure default parameters for Argon2id
func DefaultKeyDerivationParams() (*KeyDerivationParams, error) {
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	return &KeyDerivationParams{
		Algorithm:   "argon2id",
		Iterations:  3,
		Memory:      64 * 1024, // 64 MB
		Parallelism: 4,
		Salt:        salt,
		KeyLength:   32, // 256 bits for AES-256
	}, nil
}

// DeriveKey derives an encryption key from a master password using Argon2id
func DeriveKey(masterPassword string, params *KeyDerivationParams) []byte {
	return argon2.IDKey(
		[]byte(masterPassword),
		params.Salt,
		params.Iterations,
		params.Memory,
		params.Parallelism,
		params.KeyLength,
	)
}

// MarshalParams serializes key derivation parameters to JSON
func MarshalParams(params *KeyDerivationParams) ([]byte, error) {
	return json.Marshal(params)
}

// UnmarshalParams deserializes key derivation parameters from JSON
func UnmarshalParams(data []byte) (*KeyDerivationParams, error) {
	params := &KeyDerivationParams{}
	if err := json.Unmarshal(data, params); err != nil {
		return nil, fmt.Errorf("failed to unmarshal params: %w", err)
	}
	return params, nil
}
