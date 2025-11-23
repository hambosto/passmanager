package storage

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hambosto/passmanager/internal/domain/entity"
	"github.com/hambosto/passmanager/internal/infrastructure/crypto"
)

const (
	// VaultHeader is the magic header for vault files
	VaultHeader = "PMVAULT1"
	// VaultVersion is the current vault file format version
	VaultVersion = uint32(1)
)

// FileRepository implements vault storage using encrypted files
type FileRepository struct {
	path string
}

// NewFileRepository creates a new file-based vault repository
func NewFileRepository(path string) *FileRepository {
	return &FileRepository{
		path: path,
	}
}

// Save encrypts and saves the vault to a file
// File format: [Header: 8 bytes][Version: 4 bytes][KDF Params Length: 4 bytes][KDF Params][Encrypted Data]
func (r *FileRepository) Save(vault *entity.Vault, key []byte, kdfParams *crypto.KeyDerivationParams) error {
	// Ensure directory exists
	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Serialize vault to JSON
	vaultJSON, err := json.Marshal(vault)
	if err != nil {
		return fmt.Errorf("failed to marshal vault: %w", err)
	}

	// Encrypt vault data
	encrypted, err := crypto.Encrypt(vaultJSON, key)
	if err != nil {
		return fmt.Errorf("failed to encrypt vault: %w", err)
	}
	defer crypto.ZeroBytes(vaultJSON)

	// Serialize KDF params
	paramsJSON, err := crypto.MarshalParams(kdfParams)
	if err != nil {
		return fmt.Errorf("failed to marshal KDF params: %w", err)
	}

	// Create file buffer
	buf := new(bytes.Buffer)

	// Write header
	if _, err := buf.WriteString(VaultHeader); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write version
	if err := binary.Write(buf, binary.LittleEndian, VaultVersion); err != nil {
		return fmt.Errorf("failed to write version: %w", err)
	}

	// Write KDF params length
	paramsLen := uint32(len(paramsJSON))
	if err := binary.Write(buf, binary.LittleEndian, paramsLen); err != nil {
		return fmt.Errorf("failed to write params length: %w", err)
	}

	// Write KDF params
	if _, err := buf.Write(paramsJSON); err != nil {
		return fmt.Errorf("failed to write KDF params: %w", err)
	}

	// Write encrypted data
	if _, err := buf.Write(encrypted); err != nil {
		return fmt.Errorf("failed to write encrypted data: %w", err)
	}

	// Write to file atomically (write to temp file, then rename)
	tempPath := r.path + ".tmp"
	if err := os.WriteFile(tempPath, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := os.Rename(tempPath, r.path); err != nil {
		os.Remove(tempPath) // Clean up temp file
		return fmt.Errorf("failed to rename file: %w", err)
	}

	return nil
}

// Load decrypts and loads the vault from a file
// Returns: vault, kdfParams, error
func (r *FileRepository) Load(key []byte) (*entity.Vault, error) {
	// Read file
	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	// Verify minimum size
	minSize := len(VaultHeader) + 4 + 4 // header + version + params length
	if len(data) < minSize {
		return nil, fmt.Errorf("vault file corrupted: too small")
	}

	// Verify header
	header := string(data[:len(VaultHeader)])
	if header != VaultHeader {
		return nil, fmt.Errorf("invalid vault file: wrong header")
	}
	data = data[len(VaultHeader):]

	// Read version
	var version uint32
	buf := bytes.NewReader(data[:4])
	if err := binary.Read(buf, binary.LittleEndian, &version); err != nil {
		return nil, fmt.Errorf("failed to read version: %w", err)
	}
	data = data[4:]

	if version != VaultVersion {
		return nil, fmt.Errorf("unsupported vault version: %d", version)
	}

	// Read KDF params length
	var paramsLen uint32
	buf = bytes.NewReader(data[:4])
	if err := binary.Read(buf, binary.LittleEndian, &paramsLen); err != nil {
		return nil, fmt.Errorf("failed to read params length: %w", err)
	}
	data = data[4:]

	// Read KDF params
	if len(data) < int(paramsLen) {
		return nil, fmt.Errorf("vault file corrupted: params too short")
	}
	paramsJSON := data[:paramsLen]
	data = data[paramsLen:]

	// Parse KDF params (we don't actually need them here since key is already derived)
	_, err = crypto.UnmarshalParams(paramsJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal KDF params: %w", err)
	}

	// Decrypt data
	decrypted, err := crypto.Decrypt(data, key)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt vault (wrong password?): %w", err)
	}
	defer crypto.ZeroBytes(decrypted)

	// Deserialize vault
	vault := &entity.Vault{}
	if err := json.Unmarshal(decrypted, vault); err != nil {
		return nil, fmt.Errorf("failed to unmarshal vault: %w", err)
	}

	return vault, nil
}

// LoadParams loads just the KDF parameters from the vault file (without decrypting)
func (r *FileRepository) LoadParams() (*crypto.KeyDerivationParams, error) {
	// Read file
	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read vault file: %w", err)
	}

	// Verify minimum size
	minSize := len(VaultHeader) + 4 + 4
	if len(data) < minSize {
		return nil, fmt.Errorf("vault file corrupted: too small")
	}

	// Skip header
	data = data[len(VaultHeader):]

	// Skip version
	data = data[4:]

	// Read KDF params length
	var paramsLen uint32
	buf := bytes.NewReader(data[:4])
	if err := binary.Read(buf, binary.LittleEndian, &paramsLen); err != nil {
		return nil, fmt.Errorf("failed to read params length: %w", err)
	}
	data = data[4:]

	// Read KDF params
	if len(data) < int(paramsLen) {
		return nil, fmt.Errorf("vault file corrupted: params too short")
	}
	paramsJSON := data[:paramsLen]

	// Parse KDF params
	return crypto.UnmarshalParams(paramsJSON)
}

// Exists checks if the vault file exists
func (r *FileRepository) Exists() bool {
	_, err := os.Stat(r.path)
	return err == nil
}

// Delete removes the vault file
func (r *FileRepository) Delete() error {
	if err := os.Remove(r.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete vault: %w", err)
	}
	return nil
}

// GetPath returns the vault file path
func (r *FileRepository) GetPath() string {
	return r.path
}
