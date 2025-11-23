package validator

import (
	"fmt"
	"math"
	"strings"
	"unicode"
)

// PasswordStrength represents the strength level of a password
type PasswordStrength int

const (
	StrengthWeak PasswordStrength = iota
	StrengthFair
	StrengthGood
	StrengthStrong
	StrengthExcellent
)

// String returns the string representation of the strength
func (s PasswordStrength) String() string {
	switch s {
	case StrengthWeak:
		return "Weak"
	case StrengthFair:
		return "Fair"
	case StrengthGood:
		return "Good"
	case StrengthStrong:
		return "Strong"
	case StrengthExcellent:
		return "Excellent"
	default:
		return "Unknown"
	}
}

// CommonPasswords is a list of commonly used weak passwords
var CommonPasswords = []string{
	"password", "123456", "12345678", "qwerty", "abc123",
	"monkey", "1234567", "letmein", "trustno1", "dragon",
	"baseball", "111111", "iloveyou", "master", "sunshine",
	"ashley", "bailey", "passw0rd", "shadow", "123123",
	"654321", "superman", "qazwsx", "michael", "football",
}

// ValidatePassword validates a password and returns its strength
func ValidatePassword(password string, minLength int) (bool, PasswordStrength, string) {
	// Check minimum length
	if len(password) < minLength {
		return false, StrengthWeak, "Password must be at least " + string(rune(minLength)) + " characters"
	}

	// Check for common passwords
	lowerPassword := strings.ToLower(password)
	for _, common := range CommonPasswords {
		if lowerPassword == common || strings.Contains(lowerPassword, common) {
			return false, StrengthWeak, "Password is too common"
		}
	}

	// Calculate entropy and determine strength
	entropy := CalculateEntropy(password)
	strength := GetStrengthFromEntropy(entropy)

	if strength == StrengthWeak {
		return false, strength, "Password is too weak"
	}

	return true, strength, ""
}

// CalculateEntropy calculates the entropy of a password in bits
func CalculateEntropy(password string) float64 {
	if len(password) == 0 {
		return 0
	}

	// Determine character pool size
	poolSize := 0
	hasLower := false
	hasUpper := false
	hasDigit := false
	hasSymbol := false

	for _, char := range password {
		if unicode.IsLower(char) {
			hasLower = true
		} else if unicode.IsUpper(char) {
			hasUpper = true
		} else if unicode.IsDigit(char) {
			hasDigit = true
		} else {
			hasSymbol = true
		}
	}

	if hasLower {
		poolSize += 26
	}
	if hasUpper {
		poolSize += 26
	}
	if hasDigit {
		poolSize += 10
	}
	if hasSymbol {
		poolSize += 32 // Approximation for common symbols
	}

	// Entropy = log2(poolSize^length) = length * log2(poolSize)
	return float64(len(password)) * math.Log2(float64(poolSize))
}

// GetStrengthFromEntropy determines password strength based on entropy
func GetStrengthFromEntropy(entropy float64) PasswordStrength {
	switch {
	case entropy < 40:
		return StrengthWeak
	case entropy < 60:
		return StrengthFair
	case entropy < 80:
		return StrengthGood
	case entropy < 100:
		return StrengthStrong
	default:
		return StrengthExcellent
	}
}

// EstimateCrackTime estimates the time to crack a password based on entropy
// Assumes 1 billion guesses per second
func EstimateCrackTime(entropy float64) string {
	if entropy <= 0 {
		return "instantly"
	}

	// Calculate total possible combinations
	combinations := math.Pow(2, entropy)

	// Assume modern GPU can try ~1 billion hashes/second
	// For Argon2id, this is much slower, but we'll use a conservative estimate
	guessesPerSecond := 1_000_000_000.0
	seconds := combinations / guessesPerSecond / 2.0 // Divide by 2 for average case

	// Convert to human-readable format
	return formatDuration(seconds)
}

// formatDuration formats seconds into a human-readable duration
func formatDuration(seconds float64) string {
	if seconds < 1 {
		return "instantly"
	}
	if seconds < 60 {
		return fmt.Sprintf("%.0f seconds", seconds)
	}
	if seconds < 3600 {
		return fmt.Sprintf("%.0f minutes", seconds/60)
	}
	if seconds < 86400 {
		return fmt.Sprintf("%.0f hours", seconds/3600)
	}
	if seconds < 31536000 {
		return fmt.Sprintf("%.0f days", seconds/86400)
	}
	if seconds < 31536000000 {
		return fmt.Sprintf("%.0f years", seconds/31536000)
	}
	if seconds < 31536000000000 {
		return fmt.Sprintf("%.0f thousand years", seconds/31536000000)
	}
	if seconds < 31536000000000000 {
		return fmt.Sprintf("%.0f million years", seconds/31536000000000)
	}
	return fmt.Sprintf("%.0f billion years", seconds/31536000000000000)
}
