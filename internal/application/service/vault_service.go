package service

import (
	"github.com/hambosto/passmanager/internal/domain/entity"
	"github.com/hambosto/passmanager/internal/domain/repository"
	"github.com/hambosto/passmanager/internal/infrastructure/crypto"
)

// VaultService implements vault business logic
type VaultServiceImpl struct {
	repository repository.VaultRepository
	masterKey  []byte
}

// NewVaultService creates a new vault service
func NewVaultService(repo repository.VaultRepository) *VaultServiceImpl {
	return &VaultServiceImpl{
		repository: repo,
	}
}

// CreateVault creates a new vault with the given master password
func (s *VaultServiceImpl) CreateVault(masterPassword string) (*entity.Vault, error) {
	// Derive encryption key
	params, err := crypto.DefaultKeyDerivationParams()
	if err != nil {
		return nil, err
	}

	key := crypto.DeriveKey(masterPassword, params)
	s.masterKey = key

	// Create new vault
	vault := entity.NewVault()

	// Save vault with KDF params
	if err := s.repository.Save(vault, key, params); err != nil {
		return nil, err
	}

	return vault, nil
}

// UnlockVault unlocks an existing vault with the given master password
func (s *VaultServiceImpl) UnlockVault(masterPassword string) (*entity.Vault, error) {
	params, err := crypto.DefaultKeyDerivationParams()
	if err != nil {
		return nil, err
	}

	key := crypto.DeriveKey(masterPassword, params)

	vault, err := s.repository.Load(key)
	if err != nil {
		return nil, err
	}

	s.masterKey = key
	return vault, nil
}

// SaveVault saves the current vault
func (s *VaultServiceImpl) SaveVault(vault *entity.Vault) error { // Load KDF params and save vault
	if s.masterKey == nil {
		return ErrVaultLocked
	}
	params, err := s.repository.LoadParams()
	if err != nil {
		return err
	}
	return s.repository.Save(vault, s.masterKey, params)
}

// LockVault clears the master key from memory
func (s *VaultServiceImpl) LockVault() error {
	if s.masterKey != nil {
		crypto.ZeroBytes(s.masterKey)
		s.masterKey = nil
	}
	return nil
}

// ChangePassword changes the vault's master password
func (s *VaultServiceImpl) ChangePassword(oldPassword, newPassword string, vault *entity.Vault) error {
	// Verify old password
	_, err := s.UnlockVault(oldPassword)
	if err != nil {
		return err
	}

	// Derive new key
	params, err := crypto.DefaultKeyDerivationParams()
	if err != nil {
		return err
	}

	newKey := crypto.DeriveKey(newPassword, params)

	// Save with new key and new params
	if err := s.repository.Save(vault, newKey, params); err != nil {
		return err
	}

	// Update master key
	crypto.ZeroBytes(s.masterKey)
	s.masterKey = newKey

	return nil
}

// Common errors
var (
	ErrVaultLocked = &ServiceError{Code: "VAULT_LOCKED", Message: "Vault is locked"}
)

// ServiceError represents a service-level error
type ServiceError struct {
	Code    string
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}
