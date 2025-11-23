package repository

import (
	"github.com/hambosto/passmanager/internal/domain/entity"
	"github.com/hambosto/passmanager/internal/infrastructure/crypto"
)

// VaultRepository defines the interface for vault persistence
type VaultRepository interface {
	// Save saves the vault with the given key and KDF params
	Save(vault *entity.Vault, key []byte, kdfParams *crypto.KeyDerivationParams) error

	// Load loads the vault using the given key
	Load(key []byte) (*entity.Vault, error)

	// LoadParams loads the KDF parameters from the vault file
	LoadParams() (*crypto.KeyDerivationParams, error)

	// Exists checks if a vault file exists
	Exists() bool

	// Delete removes the vault from storage
	Delete() error

	// GetPath returns the storage path
	GetPath() string
}
