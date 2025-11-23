package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the application configuration
type Config struct {
	Security            SecurityConfig            `yaml:"security"`
	PasswordGenerator   PasswordGeneratorConfig   `yaml:"password_generator"`
	PassphraseGenerator PassphraseGeneratorConfig `yaml:"passphrase_generator"`
	Storage             StorageConfig             `yaml:"storage"`
	UI                  UIConfig                  `yaml:"ui"`
}

// SecurityConfig contains security-related settings
type SecurityConfig struct {
	AutoLockTimeout      int  `yaml:"auto_lock_timeout"` // minutes (0 = disabled)
	ClipboardTimeout     int  `yaml:"clipboard_timeout"` // seconds
	ClearClipboardOnLock bool `yaml:"clear_clipboard_on_lock"`
	ClearClipboardOnExit bool `yaml:"clear_clipboard_on_exit"`
	MaxUnlockAttempts    int  `yaml:"max_unlock_attempts"`
	UnlockCooldown       int  `yaml:"unlock_cooldown"` // seconds
}

// PasswordGeneratorConfig contains default password generator settings
type PasswordGeneratorConfig struct {
	Length           int  `yaml:"length"`
	IncludeUppercase bool `yaml:"include_uppercase"`
	IncludeLowercase bool `yaml:"include_lowercase"`
	IncludeNumbers   bool `yaml:"include_numbers"`
	IncludeSymbols   bool `yaml:"include_symbols"`
	ExcludeAmbiguous bool `yaml:"exclude_ambiguous"`
}

// PassphraseGeneratorConfig contains default passphrase generator settings
type PassphraseGeneratorConfig struct {
	WordCount     int    `yaml:"word_count"`
	Separator     string `yaml:"separator"`
	Capitalize    bool   `yaml:"capitalize"`
	IncludeNumber bool   `yaml:"include_number"`
}

// StorageConfig contains storage-related settings
type StorageConfig struct {
	VaultPath      string `yaml:"vault_path"`
	BackupPath     string `yaml:"backup_path"`
	AutoBackup     bool   `yaml:"auto_backup"`
	BackupInterval int    `yaml:"backup_interval"` // hours
	MaxBackups     int    `yaml:"max_backups"`
}

// UIConfig contains UI-related settings
type UIConfig struct {
	Theme             string `yaml:"theme"` // "dark", "light"
	ShowTOTPCountdown bool   `yaml:"show_totp_countdown"`
	CompactMode       bool   `yaml:"compact_mode"`
	DateFormat        string `yaml:"date_format"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "passmanager")

	return &Config{
		Security: SecurityConfig{
			AutoLockTimeout:      5,
			ClipboardTimeout:     30,
			ClearClipboardOnLock: true,
			ClearClipboardOnExit: true,
			MaxUnlockAttempts:    5,
			UnlockCooldown:       300,
		},
		PasswordGenerator: PasswordGeneratorConfig{
			Length:           16,
			IncludeUppercase: true,
			IncludeLowercase: true,
			IncludeNumbers:   true,
			IncludeSymbols:   true,
			ExcludeAmbiguous: true,
		},
		PassphraseGenerator: PassphraseGeneratorConfig{
			WordCount:     4,
			Separator:     "-",
			Capitalize:    true,
			IncludeNumber: true,
		},
		Storage: StorageConfig{
			VaultPath:      filepath.Join(configDir, "vault.enc"),
			BackupPath:     filepath.Join(configDir, "backups"),
			AutoBackup:     false,
			BackupInterval: 24,
			MaxBackups:     10,
		},
		UI: UIConfig{
			Theme:             "dark",
			ShowTOTPCountdown: true,
			CompactMode:       false,
			DateFormat:        "2006-01-02 15:04",
		},
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default config if file doesn't exist
			return DefaultConfig(), nil
		}
		return nil, err
	}

	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// SaveConfig saves configuration to a file
func (c *Config) Save(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// GetConfigPath returns the default config file path
func GetConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "passmanager", "config.yaml")
}
