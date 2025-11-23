# Password Manager - Architecture

## Overview

This document describes the architecture of the Password Manager TUI application, built using Clean Architecture principles.

## Architecture Layers

### 1. Domain Layer (`internal/domain/`)

The innermost layer containing business logic and core entities.

**Entities** (`internal/domain/entity/`):
- `entry.go` - Password entries (Login, SecureNote, Card, Identity)
- `vault.go` - Vault container with settings
- `folder.go` - Folder organization
- `user.go` - User entity for future multi-user support  
- `util.go` - Utility functions (ID generation)

**Repositories** (`internal/domain/repository/`):
- `vault_repository.go` - Interface for vault persistence

**Services** (`internal/domain/service/`):
- `vault_service.go` - Interface for vault business logic

### 2. Infrastructure Layer (`internal/infrastructure/`)

Technical implementations of interfaces defined in the domain layer.

**Crypto** (`internal/infrastructure/crypto/`):
- `encryption.go` - AES-256-GCM encryption/decryption
- `key_derivation.go` - Argon2id key derivation
- Memory zeroing for sensitive data

**Storage** (`internal/infrastructure/storage/`):
- `file_repository.go` - File-based vault storage with custom binary format

**Clipboard** (`internal/infrastructure/clipboard/`):
- `clipboard.go` - Clipboard operations with auto-clear timeout

**Auto-lock**:
- `autolock.go` - Automatic vault locking after inactivity

### 3. Application Layer (`internal/application/`)

Use cases and application services orchestrating domain logic.

**Services** (`internal/application/service/`):
- `vault_service.go` - Vault operations implementation
- `totp_service.go` - TOTP code generation and validation
- `password_generator.go` - Password and passphrase generation
- `security_service.go` - Security auditing (weak passwords, duplicates)

**DTOs** (`internal/application/dto/`):
- `requests.go` - Request/response objects for decoupling

### 4. Presentation Layer (`internal/presentation/`)

User interfaces (TUI and CLI).

**TUI** (`internal/presentation/tui/`):
- `app.go` - Main Bubble Tea application
- **Screens**:
  - `login.go` - Login/vault creation
  - `vault_list.go` - Entry list with search
  - `entry_detail.go` - Entry details with TOTP
  - `entry_editor.go` - Create/edit entries
  - `settings.go` - Configuration
  - `help.go` - Keyboard shortcuts
- **Components**:
  - `password_generator_modal.go` - Password generation modal
- **Styles**:
  - `theme.go` - UI theme and styling
- **Utilities**:
  - `util.go` - Common TUI utilities

**CLI** (`internal/presentation/cli/`):
- `commands.go` - CLI commands (get, list, generate)

## Supporting Packages (`pkg/`)

Reusable packages independent of the application:

- `pkg/totp/` - RFC 6238 TOTP implementation
- `pkg/validator/` - Password validation and strength checking

## Configuration (`config/`)

- `config.go` - Application configuration with YAML support

## Security Architecture

### Encryption Stack

```
Master Password
      ↓
Argon2id (3 iter, 64MB, 4-way)
      ↓
256-bit Key
      ↓
AES-256-GCM
      ↓
Encrypted Vault
```

**Key Security Features**:
1. **Memory-hard KDF**: Argon2id resistant to GPU attacks
2. **Authenticated encryption**: AES-256-GCM prevents tampering
3. **Memory zeroing**: Sensitive data cleared from RAM
4. **No password storage**: Only derived key kept in memory
5. **Auto-lock**: Clears keys after timeout

### Vault File Format

```
[Header: PMVAULT1] [Version: 1] [Encrypted Data]
```

Encrypted data contains JSON-serialized vault with all entries.

## Data Flow

### Login/Unlock Flow
```
User Input → Login Screen → Vault Service → Key Derivation
                                ↓
              File Repository → Decryption → Vault Loaded
```

### Entry Management Flow
```
User Action → Editor Screen → Vault Update → Vault Service
                                   ↓
                          File Repository → Encryption → Saved
```

### TOTP Flow
```
Entry with Secret → TOTP Service → Generate Code
                                       ↓
                         Entry Detail Screen → Display with Countdown
```

## Design Patterns

1. **Repository Pattern**: Abstraction for data persistence
2. **Service Layer**: Business logic separation
3. **Dependency Injection**: Services injected via constructors
4. **Observer Pattern**: Bubble Tea's message-based architecture
5. **Strategy Pattern**: Different entry types, encryption algorithms

## Testing Strategy

### Unit Tests
- Domain entities
- Encryption/decryption
- TOTP generation (RFC 6238 vectors)
- Password validation

### Integration Tests
- Vault save/load cycles
- End-to-end encryption
- Service interactions

### Fixtures
- Test vaults
- Sample entries
- TOTP test vectors

## Future Enhancements

1. **Multi-user support**: Using `user.go` entity
2. **Cloud sync**: Additional repository implementations
3. **Browser extensions**: HTTP API layer
4. **Biometric unlock**: Platform-specific auth
5. **Backup/restore**: Automated backup service
6. **Audit logging**: Security event tracking

## Dependencies

**External**:
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/bubbles` - TUI components
- `github.com/charmbracelet/lipgloss` - TUI styling
- `github.com/atotto/clipboard` - Clipboard operations
- `golang.org/x/crypto` - Argon2id
- `gopkg.in/yaml.v3` - Configuration

**Standard Library**:
- `crypto/aes`, `crypto/cipher` - Encryption
- `crypto/rand` - Secure random
- `encoding/json` - Serialization
- `time` - TOTP timing

## Build & Deployment

Makefile targets:
- `make build` - Build single platform
- `make build-all` - Multi-platform builds
- `make test` - Run tests
- `make install` - Install to system

Binary output: `build/passmanager`
