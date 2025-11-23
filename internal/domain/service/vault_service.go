package service

import "github.com/hambosto/passmanager/internal/domain/entity"

// VaultService defines the interface for vault business logic
type VaultService interface {
	// CreateVault creates a new vault
	CreateVault(masterPassword string) (*entity.Vault, error)

	// UnlockVault unlocks an existing vault
	UnlockVault(masterPassword string) (*entity.Vault, error)

	// SaveVault saves the vault
	SaveVault(vault *entity.Vault) error

	// LockVault locks the vault
	LockVault() error

	// ChangePassword changes the master password
	ChangePassword(oldPassword, newPassword string) error
}
