package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name      string
		plaintext []byte
		key       []byte
		wantErr   bool
	}{
		{
			name:      "Valid encryption and decryption",
			plaintext: []byte("Hello, World!"),
			key:       make([]byte, 32), // 32-byte (256-bit) key
			wantErr:   false,
		},
		{
			name:      "Empty plaintext",
			plaintext: []byte(""),
			key:       make([]byte, 32),
			wantErr:   false,
		},
		{
			name:      "Large plaintext",
			plaintext: bytes.Repeat([]byte("A"), 10000),
			key:       make([]byte, 32),
			wantErr:   false,
		},
		{
			name:      "Invalid key size (too short)",
			plaintext: []byte("test"),
			key:       make([]byte, 16), // 16-byte key (invalid for AES-256)
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt
			encrypted, err := Encrypt(tt.plaintext, tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Encrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Verify encrypted data is different from plaintext
			if len(tt.plaintext) > 0 && bytes.Equal(encrypted, tt.plaintext) {
				t.Error("Encrypted data should not equal plaintext")
			}

			// Decrypt
			decrypted, err := Decrypt(encrypted, tt.key)
			if err != nil {
				t.Errorf("Decrypt() error = %v", err)
				return
			}

			// Verify decrypted data matches original plaintext
			if !bytes.Equal(decrypted, tt.plaintext) {
				t.Errorf("Decrypted data doesn't match plaintext.\nGot: %v\nWant: %v",
					decrypted, tt.plaintext)
			}
		})
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	plaintext := []byte("Secret message")
	key1 := make([]byte, 32)
	key2 := make([]byte, 32)
	key2[0] = 1 // Make key2 different from key1

	// Encrypt with key1
	encrypted, err := Encrypt(plaintext, key1)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Try to decrypt with key2 (should fail)
	_, err = Decrypt(encrypted, key2)
	if err == nil {
		t.Error("Decrypt() should fail with wrong key")
	}
}

func TestDecryptCorruptedData(t *testing.T) {
	plaintext := []byte("Secret message")
	key := make([]byte, 32)

	// Encrypt
	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}

	// Corrupt the encrypted data
	encrypted[len(encrypted)/2] ^= 0xFF

	// Try to decrypt corrupted data (should fail)
	_, err = Decrypt(encrypted, key)
	if err == nil {
		t.Error("Decrypt() should fail with corrupted data")
	}
}

func TestZeroBytes(t *testing.T) {
	data := []byte("sensitive data")
	ZeroBytes(data)

	for i, b := range data {
		if b != 0 {
			t.Errorf("Byte at index %d is not zero: %v", i, b)
		}
	}
}

func BenchmarkEncrypt(b *testing.B) {
	plaintext := make([]byte, 1024) // 1KB
	key := make([]byte, 32)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Encrypt(plaintext, key)
	}
}

func BenchmarkDecrypt(b *testing.B) {
	plaintext := make([]byte, 1024) // 1KB
	key := make([]byte, 32)

	encrypted, _ := Encrypt(plaintext, key)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = Decrypt(encrypted, key)
	}
}
