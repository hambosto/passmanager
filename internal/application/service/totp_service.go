package service

import (
	"time"

	"github.com/hambosto/passmanager/pkg/totp"
)

// TOTPService handles TOTP operations
type TOTPService struct{}

// NewTOTPService creates a new TOTP service
func NewTOTPService() *TOTPService {
	return &TOTPService{}
}

// GenerateCode generates a TOTP code from a secret
func (s *TOTPService) GenerateCode(secret string) (string, time.Duration, error) {
	config := totp.DefaultConfig(secret)
	return config.GenerateCode()
}

// ValidateCode validates a TOTP code
func (s *TOTPService) ValidateCode(secret, code string) bool {
	config := totp.DefaultConfig(secret)
	return config.Validate(code)
}

// ParseURI parses an otpauth:// URI
func (s *TOTPService) ParseURI(uri string) (*totp.Config, error) {
	return totp.ParseURI(uri)
}

// GenerateURI generates an otpauth:// URI
func (s *TOTPService) GenerateURI(issuer, accountName, secret string) string {
	config := totp.DefaultConfig(secret)
	// ToURI returns the URI without parameters
	return config.ToURI()
}
