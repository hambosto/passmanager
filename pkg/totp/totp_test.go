package totp

import (
	"testing"
	"time"
)

func TestGenerateCode(t *testing.T) {
	// Test with known secret and time (from RFC 6238 test vectors)
	config := &Config{
		Secret:    "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ", // Base32 encoded "12345678901234567890"
		Period:    30 * time.Second,
		Digits:    6,
		Algorithm: "SHA1",
	}

	// Test at known time: 59 seconds (should generate "287082")
	testTime := time.Unix(59, 0)
	code, expiresIn, err := config.GenerateCodeAt(testTime)
	if err != nil {
		t.Fatalf("GenerateCodeAt() error = %v", err)
	}

	if code != "287082" {
		t.Errorf("GenerateCodeAt() code = %v, want 287082", code)
	}

	if expiresIn != 1*time.Second {
		t.Errorf("GenerateCodeAt() expiresIn = %v, want 1s", expiresIn)
	}
}

func TestGenerateCodeSHA256(t *testing.T) {
	config := &Config{
		Secret:    "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ",
		Period:    30 * time.Second,
		Digits:    6,
		Algorithm: "SHA256",
	}

	testTime := time.Unix(59, 0)
	code, _, err := config.GenerateCodeAt(testTime)
	if err != nil {
		t.Fatalf("GenerateCodeAt() error = %v", err)
	}

	// SHA256 should produce different code than SHA1
	if code == "287082" {
		t.Error("SHA256 should produce different code than SHA1")
	}

	// Code should be 6 digits
	if len(code) != 6 {
		t.Errorf("Code length = %d, want 6", len(code))
	}
}

func TestGenerateCode8Digits(t *testing.T) {
	config := &Config{
		Secret:    "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ",
		Period:    30 * time.Second,
		Digits:    8,
		Algorithm: "SHA1",
	}

	testTime := time.Unix(59, 0)
	code, _, err := config.GenerateCodeAt(testTime)
	if err != nil {
		t.Fatalf("GenerateCodeAt() error = %v", err)
	}

	if len(code) != 8 {
		t.Errorf("Code length = %d, want 8", len(code))
	}
}

func TestValidateCode(t *testing.T) {
	config := &Config{
		Secret:    "GEZDGNBVGY3TQOJQGEZDGNBVGY3TQOJQ",
		Period:    30 * time.Second,
		Digits:    6,
		Algorithm: "SHA1",
	}

	testTime := time.Unix(59, 0)

	// Valid code
	if !config.ValidateAt("287082", testTime) {
		t.Error("ValidateAt() should accept correct code")
	}

	// Invalid code
	if config.ValidateAt("000000", testTime) {
		t.Error("ValidateAt() should reject incorrect code")
	}

	// Code from previous time window (should still be valid due to ±1 tolerance)
	prevTime := testTime.Add(-30 * time.Second)
	prevCode, _, _ := config.GenerateCodeAt(prevTime)
	if !config.ValidateAt(prevCode, testTime) {
		t.Error("ValidateAt() should accept code from previous time window")
	}
}

func TestParseURI(t *testing.T) {
	tests := []struct {
		name    string
		uri     string
		want    *Config
		wantErr bool
	}{
		{
			name: "Basic URI",
			uri:  "otpauth://totp/Test?secret=JBSWY3DPEHPK3PXP",
			want: &Config{
				Secret:    "JBSWY3DPEHPK3PXP",
				Period:    30 * time.Second,
				Digits:    6,
				Algorithm: "SHA1",
				Account:   "Test",
			},
			wantErr: false,
		},
		{
			name: "URI with issuer",
			uri:  "otpauth://totp/GitHub:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=GitHub",
			want: &Config{
				Secret:    "JBSWY3DPEHPK3PXP",
				Period:    30 * time.Second,
				Digits:    6,
				Algorithm: "SHA1",
				Issuer:    "GitHub",
				Account:   "user@example.com",
			},
			wantErr: false,
		},
		{
			name: "URI with custom parameters",
			uri:  "otpauth://totp/Test?secret=JBSWY3DPEHPK3PXP&period=60&digits=8&algorithm=SHA256",
			want: &Config{
				Secret:    "JBSWY3DPEHPK3PXP",
				Period:    60 * time.Second,
				Digits:    8,
				Algorithm: "SHA256",
				Account:   "Test",
			},
			wantErr: false,
		},
		{
			name:    "Missing secret",
			uri:     "otpauth://totp/Test",
			wantErr: true,
		},
		{
			name:    "Invalid scheme",
			uri:     "http://totp/Test?secret=JBSWY3DPEHPK3PXP",
			wantErr: true,
		},
		{
			name:    "Invalid type (not totp)",
			uri:     "otpauth://hotp/Test?secret=JBSWY3DPEHPK3PXP",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseURI(tt.uri)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseURI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if got.Secret != tt.want.Secret {
				t.Errorf("Secret = %v, want %v", got.Secret, tt.want.Secret)
			}
			if got.Period != tt.want.Period {
				t.Errorf("Period = %v, want %v", got.Period, tt.want.Period)
			}
			if got.Digits != tt.want.Digits {
				t.Errorf("Digits = %v, want %v", got.Digits, tt.want.Digits)
			}
			if got.Algorithm != tt.want.Algorithm {
				t.Errorf("Algorithm = %v, want %v", got.Algorithm, tt.want.Algorithm)
			}
			if got.Issuer != tt.want.Issuer {
				t.Errorf("Issuer = %v, want %v", got.Issuer, tt.want.Issuer)
			}
			if got.Account != tt.want.Account {
				t.Errorf("Account = %v, want %v", got.Account, tt.want.Account)
			}
		})
	}
}

func TestToURI(t *testing.T) {
	config := &Config{
		Secret:    "JBSWY3DPEHPK3PXP",
		Period:    30 * time.Second,
		Digits:    6,
		Algorithm: "SHA1",
		Issuer:    "GitHub",
		Account:   "user@example.com",
	}

	uri := config.ToURI()

	// Parse it back
	parsed, err := ParseURI(uri)
	if err != nil {
		t.Fatalf("ParseURI() error = %v", err)
	}

	if parsed.Secret != config.Secret {
		t.Errorf("Secret = %v, want %v", parsed.Secret, config.Secret)
	}
	if parsed.Issuer != config.Issuer {
		t.Errorf("Issuer = %v, want %v", parsed.Issuer, config.Issuer)
	}
	if parsed.Account != config.Account {
		t.Errorf("Account = %v, want %v", parsed.Account, config.Account)
	}
}

func BenchmarkGenerateCode(b *testing.B) {
	config := DefaultConfig("JBSWY3DPEHPK3PXP")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, _ = config.GenerateCode()
	}
}
