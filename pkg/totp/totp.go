package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Config represents TOTP configuration parameters
type Config struct {
	Secret    string        // Base32 encoded secret
	Period    time.Duration // Time period (default 30s)
	Digits    int           // Number of digits (default 6)
	Algorithm string        // Hash algorithm: "SHA1", "SHA256", "SHA512"
	Issuer    string        // Optional issuer name
	Account   string        // Optional account name
}

// DefaultConfig returns a standard TOTP configuration
func DefaultConfig(secret string) *Config {
	return &Config{
		Secret:    secret,
		Period:    30 * time.Second,
		Digits:    6,
		Algorithm: "SHA1",
	}
}

// GenerateCode generates the current TOTP code
// Returns: code, time until expiry, error
func (c *Config) GenerateCode() (string, time.Duration, error) {
	return c.GenerateCodeAt(time.Now())
}

// GenerateCodeAt generates a TOTP code for a specific time
func (c *Config) GenerateCodeAt(t time.Time) (string, time.Duration, error) {
	// Decode secret
	secret, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(strings.ToUpper(c.Secret))
	if err != nil {
		return "", 0, fmt.Errorf("invalid secret: %w", err)
	}

	// Calculate time counter
	counter := uint64(t.Unix()) / uint64(c.Period.Seconds())

	// Generate HOTP
	code, err := c.generateHOTP(secret, counter)
	if err != nil {
		return "", 0, err
	}

	// Calculate time until expiry
	nextCounter := (counter + 1) * uint64(c.Period.Seconds())
	expiresIn := time.Duration(nextCounter-uint64(t.Unix())) * time.Second

	return code, expiresIn, nil
}

// generateHOTP generates an HOTP code
func (c *Config) generateHOTP(secret []byte, counter uint64) (string, error) {
	// Create HMAC hash function
	var h func() hash.Hash
	switch c.Algorithm {
	case "SHA1":
		h = sha1.New
	case "SHA256":
		h = sha256.New
	case "SHA512":
		h = sha512.New
	default:
		return "", fmt.Errorf("unsupported algorithm: %s", c.Algorithm)
	}

	// Create HMAC
	mac := hmac.New(h, secret)

	// Write counter as big-endian uint64
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, counter)
	mac.Write(counterBytes)

	// Compute HMAC
	hmacResult := mac.Sum(nil)

	// Dynamic truncation
	offset := hmacResult[len(hmacResult)-1] & 0xf
	code := binary.BigEndian.Uint32(hmacResult[offset:offset+4]) & 0x7fffffff

	// Generate digits
	modulo := uint32(math.Pow10(c.Digits))
	code = code % modulo

	// Format with leading zeros
	format := fmt.Sprintf("%%0%dd", c.Digits)
	return fmt.Sprintf(format, code), nil
}

// Validate checks if a code is valid for the current time window
// Allows for ±1 time step tolerance
func (c *Config) Validate(code string) bool {
	return c.ValidateAt(code, time.Now())
}

// ValidateAt checks if a code is valid for a specific time
func (c *Config) ValidateAt(code string, t time.Time) bool {
	// Check current time
	currentCode, _, err := c.GenerateCodeAt(t)
	if err == nil && currentCode == code {
		return true
	}

	// Check ±1 time step
	prevTime := t.Add(-c.Period)
	prevCode, _, err := c.GenerateCodeAt(prevTime)
	if err == nil && prevCode == code {
		return true
	}

	nextTime := t.Add(c.Period)
	nextCode, _, err := c.GenerateCodeAt(nextTime)
	if err == nil && nextCode == code {
		return true
	}

	return false
}

// ParseURI parses an otpauth:// URI
// Format: otpauth://totp/{issuer}:{account}?secret={secret}&issuer={issuer}&period={period}&digits={digits}&algorithm={algorithm}
func ParseURI(uri string) (*Config, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("invalid URI: %w", err)
	}

	if u.Scheme != "otpauth" {
		return nil, fmt.Errorf("invalid scheme: expected otpauth, got %s", u.Scheme)
	}

	if u.Host != "totp" {
		return nil, fmt.Errorf("invalid type: expected totp, got %s", u.Host)
	}

	// Parse query parameters
	query := u.Query()
	secret := query.Get("secret")
	if secret == "" {
		return nil, fmt.Errorf("missing secret parameter")
	}

	config := DefaultConfig(secret)

	// Parse optional parameters
	if issuer := query.Get("issuer"); issuer != "" {
		config.Issuer = issuer
	}

	if period := query.Get("period"); period != "" {
		p, err := strconv.Atoi(period)
		if err != nil {
			return nil, fmt.Errorf("invalid period: %w", err)
		}
		config.Period = time.Duration(p) * time.Second
	}

	if digits := query.Get("digits"); digits != "" {
		d, err := strconv.Atoi(digits)
		if err != nil {
			return nil, fmt.Errorf("invalid digits: %w", err)
		}
		config.Digits = d
	}

	if algorithm := query.Get("algorithm"); algorithm != "" {
		config.Algorithm = strings.ToUpper(algorithm)
	}

	// Parse label (format: issuer:account or just account)
	label := strings.TrimPrefix(u.Path, "/")
	if label != "" {
		parts := strings.SplitN(label, ":", 2)
		if len(parts) == 2 {
			if config.Issuer == "" {
				config.Issuer = parts[0]
			}
			config.Account = parts[1]
		} else {
			config.Account = parts[0]
		}
	}

	return config, nil
}

// ToURI converts the config to an otpauth:// URI
func (c *Config) ToURI() string {
	values := url.Values{}
	values.Set("secret", c.Secret)

	if c.Issuer != "" {
		values.Set("issuer", c.Issuer)
	}

	if c.Period != 30*time.Second {
		values.Set("period", strconv.Itoa(int(c.Period.Seconds())))
	}

	if c.Digits != 6 {
		values.Set("digits", strconv.Itoa(c.Digits))
	}

	if c.Algorithm != "SHA1" {
		values.Set("algorithm", c.Algorithm)
	}

	// Build label
	label := c.Account
	if c.Issuer != "" && c.Account != "" {
		label = c.Issuer + ":" + c.Account
	} else if c.Issuer != "" {
		label = c.Issuer
	}

	return fmt.Sprintf("otpauth://totp/%s?%s", url.PathEscape(label), values.Encode())
}
