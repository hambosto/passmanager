# Security Documentation

## Overview

This password manager implements multiple layers of security to protect your sensitive information. This document details the security architecture, threat model, and best practices.

## Threat Model

### What We Protect Against

✅ **Offline attacks** - Encrypted vault file cannot be decrypted without master password  
✅ **Brute force attacks** - Argon2id makes password cracking computationally expensive  
✅ **Data tampering** - AES-GCM provides authenticated encryption  
✅ **Memory dumps** - Sensitive data is zeroed after use  
✅ **Clipboard snooping** - Auto-clear prevents long-term clipboard exposure  

###What We Don't Protect Against

❌ **Keyloggers** - Master password entered via keyboard can be captured  
❌ **Screen recording** - TUI content visible on screen  
❌ **Physical access** - Unlocked vault accessible to anyone at the terminal  
❌ **Compromised system** - Malware with root access can extract keys from memory  
❌ **Social engineering** - Users must protect their master password  

## Encryption Details

### Vault Encryption

- **Algorithm**: AES-256-GCM (Galois/Counter Mode)
- **Key Size**: 256 bits (32 bytes)
- **Nonce**: 96 bits (12 bytes), randomly generated per encryption
- **Authentication**: Built-in AEAD with 128-bit auth tag

**Why AES-GCM?**
- Industry standard, widely audited
- Authenticated encryption (detects tampering)
- Fast hardware acceleration on modern CPUs
- No padding oracle attacks

### Key Derivation

- **Algorithm**: Argon2id
- **Version**: Argon2 v1.3
- **Parameters**:
  - Time cost (iterations): 3
  - Memory cost: 64 MB (65536 KB)
  - Parallelism: 4 threads
  - Salt: 32 bytes (256 bits), randomly generated
  - Output: 32 bytes (256 bits)

**Why Argon2id?**
- Winner of Password Hashing Competition (PHC)
- Memory-hard: resistant to GPU/ASIC attacks
- Side-channel resistant (id variant)
- Configurable parameters for future-proofing

**Cost Analysis:**
```
Single hash attempt: ~50-100ms on modern CPU
Brute force for 8-char password: ~10^15 attempts = millions of years
```

### File Format

```
[Header: 8 bytes]     "PMVAULT1"
[Version: 4 bytes]    Little-endian uint32
[Encrypted Data]      Nonce (12 bytes) + Ciphertext + Auth Tag (16 bytes)
```

The key derivation parameters are currently hard-coded in the application. Future versions may store them in the file format.

## Security Features

### Master Password

**Requirements:**
-Minimum 8 characters (16+ recommended)
- No strength validation on creation yet (planned feature)

**Best Practices:**
- Use a passphrase (4+ random words)
- Don't reuse passwords from other services
- Consider using a password manager for your master password (yes, seriously)
- Write it down and store securely if needed - better than forgetting

**What NOT to do:**
- ❌ Use common passwords ("password123", "qwerty")
- ❌ Use personal information (birthday, name, etc.)
- ❌ Share your master password
- ❌ Store it in plain text files

### Auto-Lock

**Default**: 5 minutes of inactivity  
**Configurable**: 0 (disabled) to any value in minutes

When locked:
- Master key cleared from memory
- Clipboard cleared
- Vault must be unlocked with master password

**Triggers:**
- Inactivity timeout
- Manual lock (Ctrl+L)
- Application exit

### Clipboard Security

**Default timeout**: 30 seconds  
**Configurable**: 1-300 seconds

**Behavior:**
- Copied passwords auto-cleared after timeout
- Cleared on vault lock
- Cleared on application exit
- Visual countdown shown

**Limitations:**
- Uses system clipboard (visible to other applications during timeout)
- Clipboard managers may keep history
- Consider terminal multiplexer's clipboard

### Memory Security

**Implemented:**
- Sensitive byte slices zeroed after use
- Master password not stored (only derived key, temporarily)
- Encryption key cleared on vault lock

**Limitations:**
- Go's garbage collector may leave copies in memory
- No memory locking (mlock) - OS may swap to disk
- Process memory dumps can expose unlocked vault

## Security Best Practices

### For Users

1. **Strong Master Password**
   ```
   Good: "correct horse battery staple 92"
   Bad:  "password123"
   ```

2. **Secure Your System**
   - Use full-disk encryption
   - Keep OS and applications updated
   - Use antivirus/antimalware
   - Don't run untrusted software

3. **Regular Backups**
   - Encrypted backups only
   - Store in multiple secure locations
   - Test restoration periodically

4. **Physical Security**
   - Lock your screen when away
   - Don't leave vault unlocked
   - Be aware of shoulder surfing

5. **TOTP Best Practices**
   - Store TOTP secrets separately from passwords when possible
   - Keep backup codes in a separate secure note
   - Don't share QR codes or secrets

### For Developers

1. **Never Log Sensitive Data**
   - No passwords in logs
   - No encryption keys in logs
   - No TOTP secrets in logs

2. **Zero Sensitive Memory**
   ```go
   defer crypto.ZeroBytes(masterPasswordBytes)
   defer crypto.ZeroBytes(encryptionKey)
   ```

3. **Constant-Time Comparisons**
   - Use `subtle.ConstantTimeCompare` for password verification
   - Prevent timing attacks

4. **Input Validation**
   - Validate all user inputs
   - Sanitize before encryption
   - Check bounds and types

5. **Secure Random Generation**
   - Use `crypto/rand` for all random data
   - Check for errors
   - Never use `math/rand` for security

## Vulnerability Reporting

If you discover a security vulnerability, please email:
**security@example.com**

Please do not create public GitHub issues for security vulnerabilities.

### Responsible Disclosure

1. Email details to security contact
2. Wait for acknowledgment (24-48 hours)
3. Allow 90 days for fix before public disclosure
4. Coordinate disclosure timing

We appreciate security researchers and will acknowledge contributors.

## Known Limitations

### Terminal Security

- **Visible on screen**: Anyone with physical or remote access to your terminal can see the vault when unlocked
- **Terminal history**: Commands may be logged in shell history
- **Screen sharing**: Be careful during video calls

### Clipboard

- **System-wide**: Other apps can read clipboard
- **Clipboard managers**: May keep history
- **Remote desktop**: Clipboard may sync to remote system

### No Network Sync

- **No built-in sync**: Must manually transfer vault file
- **File conflicts**: Can occur with multiple devices
- **Transfer security**: Use encrypted channels (SSH, HTTPS)

### Platform-Specific

- **File permissions**: Depend on OS umask settings
- **Memory protection**: No mlock, vulnerable to swap
- **Clipboard APIs**: Platform-dependent behavior

## Compliance

### Standards

- **NIST SP 800-132**: Password-Based Key Derivation (via Argon2id)
- **NIST SP 800-38D**: AES-GCM recommendation
- **RFC 6238**: TOTP specification
- **OWASP**: Password storage best practices

### Cryptographic Libraries

All cryptography uses Go's standard library (`golang.org/x/crypto`):
- Maintained by Go team
- Regular security audits
- Peer-reviewed implementations

## Security Audits

**Status**: No professional security audit performed yet

**Self-Assessment:**
- Code reviewed by author
- Unit tests for crypto operations
- Test vectors from RFC 6238 for TOTP

**Future Plans:**
- Professional security audit
- Penetration testing
- Bug bounty program

## Incident Response

If a security breach occurs:

1. **Immediate Actions**:
   - Change master password
   - Create new vault
   - Rotate all stored passwords
   - Review access logs

2. **Investigation**:
   - Determine scope of breach
   - Identify attack vector
   - Document timeline

3. **Remediation**:
   - Patch vulnerabilities
   - Update security practices
   - Notify affected users

4. **Post-Mortem**:
   - Root cause analysis
   - Improve security measures
   - Update documentation

## References

- [Argon2 Specification](https://github.com/P-H-C/phc-winner-argon2/blob/master/argon2-specs.pdf)
- [AES-GCM](https://nvlpubs.nist.gov/nistpubs/Legacy/SP/nistspecialpublication800-38d.pdf)
- [RFC 6238 - TOTP](https://tools.ietf.org/html/rfc6238)
- [OWASP Password Storage](https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html)
