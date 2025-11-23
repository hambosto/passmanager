package entity

import "time"

// Vault represents the entire encrypted vault
type Vault struct {
	Version   string    `json:"version"`
	Entries   []*Entry  `json:"entries"`
	Folders   []*Folder `json:"folders"`
	Settings  Settings  `json:"settings"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Settings represents vault-specific settings
type Settings struct {
	AutoLockTimeout    int  `json:"auto_lock_timeout"` // minutes
	ClipboardTimeout   int  `json:"clipboard_timeout"` // seconds
	PasswordGenLength  int  `json:"password_gen_length"`
	PasswordGenUpper   bool `json:"password_gen_upper"`
	PasswordGenLower   bool `json:"password_gen_lower"`
	PasswordGenNumbers bool `json:"password_gen_numbers"`
	PasswordGenSymbols bool `json:"password_gen_symbols"`
}

// NewVault creates a new vault with default settings
func NewVault() *Vault {
	now := time.Now()
	return &Vault{
		Version:   "1.0",
		Entries:   make([]*Entry, 0),
		Folders:   make([]*Folder, 0),
		Settings:  DefaultSettings(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// DefaultSettings returns the default vault settings
func DefaultSettings() Settings {
	return Settings{
		AutoLockTimeout:    5,
		ClipboardTimeout:   30,
		PasswordGenLength:  16,
		PasswordGenUpper:   true,
		PasswordGenLower:   true,
		PasswordGenNumbers: true,
		PasswordGenSymbols: true,
	}
}

// AddEntry adds an entry to the vault
func (v *Vault) AddEntry(entry *Entry) {
	v.Entries = append(v.Entries, entry)
	v.UpdatedAt = time.Now()
}

// RemoveEntry removes an entry from the vault by ID
func (v *Vault) RemoveEntry(id string) bool {
	for i, entry := range v.Entries {
		if entry.ID == id {
			v.Entries = append(v.Entries[:i], v.Entries[i+1:]...)
			v.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// FindEntry finds an entry by ID
func (v *Vault) FindEntry(id string) *Entry {
	for _, entry := range v.Entries {
		if entry.ID == id {
			return entry
		}
	}
	return nil
}

// AddFolder adds a folder to the vault
func (v *Vault) AddFolder(folder *Folder) {
	v.Folders = append(v.Folders, folder)
	v.UpdatedAt = time.Now()
}

// RemoveFolder removes a folder from the vault by ID
func (v *Vault) RemoveFolder(id string) bool {
	for i, folder := range v.Folders {
		if folder.ID == id {
			v.Folders = append(v.Folders[:i], v.Folders[i+1:]...)
			v.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}

// FindFolder finds a folder by ID
func (v *Vault) FindFolder(id string) *Folder {
	for _, folder := range v.Folders {
		if folder.ID == id {
			return folder
		}
	}
	return nil
}

// Update updates the vault timestamp
func (v *Vault) Update() {
	v.UpdatedAt = time.Now()
}
