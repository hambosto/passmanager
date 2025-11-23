package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/hambosto/passmanager/pkg/validator"
)

// PasswordConfig contains configuration for password generation
type PasswordConfig struct {
	Length           int
	IncludeUpper     bool
	IncludeLower     bool
	IncludeNumbers   bool
	IncludeSymbols   bool
	ExcludeAmbiguous bool
	MinUpper         int
	MinLower         int
	MinNumbers       int
	MinSymbols       int
}

// PassphraseConfig contains configuration for passphrase generation
type PassphraseConfig struct {
	WordCount     int
	Separator     string
	Capitalize    bool
	IncludeNumber bool
}

const (
	uppercaseLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowercaseLetters = "abcdefghijklmnopqrstuvwxyz"
	numbers          = "0123456789"
	symbols          = "!@#$%^&*()-_=+[]{}|;:,.<>?"
	ambiguous        = "0O1lI"
)

// DefaultPasswordConfig returns a sensible default password configuration
func DefaultPasswordConfig() PasswordConfig {
	return PasswordConfig{
		Length:           16,
		IncludeUpper:     true,
		IncludeLower:     true,
		IncludeNumbers:   true,
		IncludeSymbols:   true,
		ExcludeAmbiguous: true,
		MinUpper:         1,
		MinLower:         1,
		MinNumbers:       1,
		MinSymbols:       1,
	}
}

// DefaultPassphraseConfig returns a sensible default passphrase configuration
func DefaultPassphraseConfig() PassphraseConfig {
	return PassphraseConfig{
		WordCount:     4,
		Separator:     "-",
		Capitalize:    true,
		IncludeNumber: true,
	}
}

// GeneratePassword generates a random password based on the configuration
func GeneratePassword(config PasswordConfig) (string, error) {
	// Validate configuration
	if config.Length < 4 {
		return "", fmt.Errorf("password length must be at least 4")
	}

	minRequired := config.MinUpper + config.MinLower + config.MinNumbers + config.MinSymbols
	if minRequired > config.Length {
		return "", fmt.Errorf("minimum requirements exceed password length")
	}

	// Build character pool
	var pool string
	if config.IncludeUpper {
		pool += uppercaseLetters
	}
	if config.IncludeLower {
		pool += lowercaseLetters
	}
	if config.IncludeNumbers {
		pool += numbers
	}
	if config.IncludeSymbols {
		pool += symbols
	}

	if pool == "" {
		return "", fmt.Errorf("no character sets selected")
	}

	// Remove ambiguous characters if requested
	if config.ExcludeAmbiguous {
		pool = removeChars(pool, ambiguous)
	}

	// Generate password with minimum requirements
	password := make([]byte, config.Length)

	// First, satisfy minimum requirements
	pos := 0
	if config.MinUpper > 0 && config.IncludeUpper {
		chars := uppercaseLetters
		if config.ExcludeAmbiguous {
			chars = removeChars(chars, ambiguous)
		}
		for i := 0; i < config.MinUpper; i++ {
			password[pos] = randomChar(chars)
			pos++
		}
	}

	if config.MinLower > 0 && config.IncludeLower {
		chars := lowercaseLetters
		if config.ExcludeAmbiguous {
			chars = removeChars(chars, ambiguous)
		}
		for i := 0; i < config.MinLower; i++ {
			password[pos] = randomChar(chars)
			pos++
		}
	}

	if config.MinNumbers > 0 && config.IncludeNumbers {
		chars := numbers
		if config.ExcludeAmbiguous {
			chars = removeChars(chars, ambiguous)
		}
		for i := 0; i < config.MinNumbers; i++ {
			password[pos] = randomChar(chars)
			pos++
		}
	}

	if config.MinSymbols > 0 && config.IncludeSymbols {
		for i := 0; i < config.MinSymbols; i++ {
			password[pos] = randomChar(symbols)
			pos++
		}
	}

	// Fill remaining positions with random characters from pool
	for i := pos; i < config.Length; i++ {
		password[i] = randomChar(pool)
	}

	// Shuffle the password to distribute required characters
	shuffle(password)

	return string(password), nil
}

// GeneratePassphrase generates a random passphrase based on the configuration
func GeneratePassphrase(config PassphraseConfig) (string, error) {
	if config.WordCount < 1 {
		return "", fmt.Errorf("word count must be at least 1")
	}

	// Select random words from the word list
	words := make([]string, config.WordCount)
	for i := 0; i < config.WordCount; i++ {
		wordIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(effWordList))))
		if err != nil {
			return "", fmt.Errorf("failed to generate random word: %w", err)
		}
		words[i] = effWordList[wordIndex.Int64()]

		// Capitalize if requested
		if config.Capitalize {
			words[i] = strings.Title(words[i])
		}
	}

	// Join with separator
	passphrase := strings.Join(words, config.Separator)

	// Add number if requested
	if config.IncludeNumber {
		num, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %w", err)
		}
		passphrase += config.Separator + fmt.Sprintf("%02d", num.Int64())
	}

	return passphrase, nil
}

// randomChar returns a random character from the given string
func randomChar(chars string) byte {
	if len(chars) == 0 {
		return 0
	}
	index, err := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	if err != nil {
		panic(err)
	}
	return chars[index.Int64()]
}

// removeChars removes all occurrences of characters in 'remove' from 'str'
func removeChars(str, remove string) string {
	result := ""
	for _, char := range str {
		if !strings.ContainsRune(remove, char) {
			result += string(char)
		}
	}
	return result
}

// shuffle randomly shuffles a byte slice
func shuffle(data []byte) {
	for i := len(data) - 1; i > 0; i-- {
		j, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			panic(err)
		}
		data[i], data[j.Int64()] = data[j.Int64()], data[i]
	}
}

// CalculatePasswordEntropy calculates the entropy of a password
func CalculatePasswordEntropy(password string) float64 {
	return validator.CalculateEntropy(password)
}

// EstimatePasswordCrackTime estimates the time to crack a password
func EstimatePasswordCrackTime(password string) string {
	entropy := validator.CalculateEntropy(password)
	return validator.EstimateCrackTime(entropy)
}

// GetPasswordStrength returns the strength level of a password
func GetPasswordStrength(password string) validator.PasswordStrength {
	entropy := validator.CalculateEntropy(password)
	return validator.GetStrengthFromEntropy(entropy)
}

// effWordList is a subset of the EFF long wordlist for passphrase generation
// Full list would contain 7776 words, using a small subset here for demonstration
var effWordList = []string{
	"aardvark", "absurd", "accrue", "acme", "adrift", "adult", "afflict", "ahead",
	"aimless", "algae", "allow", "alone", "amuse", "angel", "armor", "arrow",
	"bamboo", "basket", "battery", "beach", "beaver", "beside", "between", "beyond",
	"cable", "camera", "campus", "canyon", "captain", "castle", "casual", "caught",
	"damage", "dance", "danger", "daring", "dash", "dawn", "decent", "decide",
	"eagle", "earth", "easy", "echo", "edge", "effort", "eight", "either",
	"fabric", "face", "fact", "fade", "faint", "false", "fancy", "fatal",
	"galaxy", "game", "gap", "garden", "gather", "gave", "gear", "general",
	"habit", "half", "hammer", "hand", "handle", "hang", "happen", "happy",
	"ice", "icon", "idea", "ideal", "identify", "idle", "ignore", "image",
	"jacket", "jazz", "join", "joint", "joke", "judge", "juice", "jump",
	"keep", "ketchup", "kettle", "keyboard", "kickoff", "kinson", "kitchen", "kite",
	"label", "ladder", "lady", "lagoon", "lamp", "language", "large", "laser",
	"machine", "macro", "madness", "magic", "magnet", "maiden", "mailbox", "major",
	"nanny", "napkin", "narrow", "nation", "native", "nature", "naval", "necklace",
	"oak", "oasis", "oath", "obese", "object", "observe", "obtain", "ocean",
	"pacific", "package", "pagan", "pager", "palace", "palm", "panel", "panic",
	"quantum", "quarter", "queen", "query", "question", "queue", "quick", "quiet",
	"race", "radar", "radio", "rage", "railway", "rainbow", "random", "range",
	"sack", "sacred", "saddle", "safari", "safe", "safety", "saga", "sage",
	"table", "tackle", "tactics", "tadpole", "talent", "talking", "tango", "tank",
	"ultimate", "umbrella", "umpire", "unable", "uncover", "undergo", "unfair", "unfold",
	"vacancy", "vaccine", "vacuum", "vague", "valid", "valley", "valve", "vampire",
	"wage", "wagon", "waist", "wallet", "walnut", "walrus", "warfare", "warm",
	"xerox", "xray",
	"yacht", "yahoo", "yarn", "year", "yellow", "yield", "yodel", "yoga",
	"zebra", "zenith", "zero", "zigzag", "zinc", "zipper", "zombie", "zone",
}
