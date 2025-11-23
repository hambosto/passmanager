package service

import (
	"fmt"

	"github.com/hambosto/passmanager/internal/domain/entity"
	"github.com/hambosto/passmanager/pkg/validator"
)

// SecurityService handles security-related operations
type SecurityService struct{}

// NewSecurityService creates a new security service
func NewSecurityService() *SecurityService {
	return &SecurityService{}
}

// CheckPasswordStrength checks the strength of a password
func (s *SecurityService) CheckPasswordStrength(password string) validator.PasswordStrength {
	entropy := validator.CalculateEntropy(password)
	return validator.GetStrengthFromEntropy(entropy)
}

// ValidatePassword validates a password
func (s *SecurityService) ValidatePassword(password string, minLength int) error {
	valid, _, message := validator.ValidatePassword(password, minLength)
	if !valid {
		return fmt.Errorf("%s", message)
	}
	return nil
}

// FindWeakPasswords finds entries with weak passwords
func (s *SecurityService) FindWeakPasswords(vault *entity.Vault) []*entity.Entry {
	var weak []*entity.Entry

	for _, entry := range vault.Entries {
		if entry.Password == "" {
			continue
		}

		strength := s.CheckPasswordStrength(entry.Password)
		if strength == validator.StrengthWeak {
			weak = append(weak, entry)
		}
	}

	return weak
}

// FindDuplicatePasswords finds entries with duplicate passwords
func (s *SecurityService) FindDuplicatePasswords(vault *entity.Vault) map[string][]*entity.Entry {
	passwordMap := make(map[string][]*entity.Entry)

	for _, entry := range vault.Entries {
		if entry.Password == "" {
			continue
		}

		passwordMap[entry.Password] = append(passwordMap[entry.Password], entry)
	}

	// Filter to only duplicates
	duplicates := make(map[string][]*entity.Entry)
	for password, entries := range passwordMap {
		if len(entries) > 1 {
			duplicates[password] = entries
		}
	}

	return duplicates
}

// CalculateSecurityScore calculates an overall security score for the vault
func (s *SecurityService) CalculateSecurityScore(vault *entity.Vault) float64 {
	if len(vault.Entries) == 0 {
		return 100.0
	}

	weakCount := len(s.FindWeakPasswords(vault))
	duplicateGroups := len(s.FindDuplicatePasswords(vault))

	totalIssues := weakCount + duplicateGroups
	score := 100.0 - (float64(totalIssues) / float64(len(vault.Entries)) * 100.0)

	if score < 0 {
		score = 0
	}

	return score
}
