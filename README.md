# Password Manager

A secure, terminal-based password manager with TOTP support, built in Go.

## Features

- 🔐 **Strong Encryption**: AES-256-GCM encryption with Argon2id key derivation
- 🔑 **TOTP Support**: RFC 6238 compliant 2FA code generation
- 🎨 **Beautiful TUI**: Modern terminal user interface built with Bubble Tea
- 📁 **Folder Organization**: Organize entries into folders
- 🔍 **Fast Search**: Real-time filtering and search
- 📋 **Smart Clipboard**: Auto-clear clipboard after timeout
- 💪 **Password Generator**: Generate strong passwords and passphrases
- 🛡️ **Security Audit**: Password health checks and strength validation
- 📥 **Import/Export**: Compatible with Bitwarden, 1Password, LastPass formats

## Installation

### From Source

```bash
git clone https://github.com/hambosto/passmanager.git
cd passmanager
make build
make install
```

### Binary Downloads

Download pre-built binaries from the [Releases](https://github.com/hambosto/passmanager/releases) page.

## Quick Start

```bash
# Start the application
passmanager

# On first run, you'll be prompted to create a master password
# The vault will be stored at ~/.config/passmanager/vault.enc
```

## Usage

### Keyboard Shortcuts

**General:**
- `Ctrl+Q` - Quit application
- `Ctrl+L` - Lock vault
- `Esc` - Go back / Cancel
- `?` - Show help

**Navigation:**
- `↑/↓` or `k/j` - Move up/down
- `Enter` - View/Open entry
- `/` - Search

**Entry Management:**
- `Ctrl+N` - New entry
- `Ctrl+E` - Edit entry
- `Ctrl+D` - Delete entry
- `Space` - Toggle favorite

**Clipboard:**
- `Ctrl+C` - Copy password
- `Ctrl+U` - Copy username
- `Ctrl+T` - Copy TOTP code

**Other:**
- `Ctrl+G` - Generate password
- `Ctrl+H` - Show/Hide password
- `Ctrl+,` - Open settings

### TOTP Setup

To add TOTP to an entry:

1. Create or edit an entry
2. Enter the TOTP secret or otpauth:// URI
3. The TOTP code will be displayed with a countdown timer

Example otpauth URI:
```
otpauth://totp/GitHub:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=GitHub
```

## Security

- **Master password never stored** - Only used to derive encryption key
- **Argon2id key derivation** - Memory-hard, resistant to GPU attacks
- **AES-256-GCM encryption** - Authenticated encryption
- **Auto-lock** - Locks vault after configured inactivity
- **Clipboard auto-clear** - Clears sensitive data from clipboard
- **Memory security** - Sensitive data zeroed after use

See [docs/SECURITY.md](docs/SECURITY.md) for detailed security information.

## Configuration

Configuration file located at `~/.config/passmanager/config.yaml`

```yaml
security:
  auto_lock_timeout: 5        # minutes
  clipboard_timeout: 30       # seconds

password_generator:
  length: 16
  include_uppercase: true
  include_lowercase: true
  include_numbers: true
  include_symbols: true
```

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run benchmarks
make bench
```

### Project Structure

```
passmanager/
├── cmd/passmanager/        # Main application entry point
├── internal/
│   ├── domain/             # Domain entities and interfaces
│   ├── infrastructure/     # Crypto, storage, clipboard
│   ├── application/        # Business logic services
│   └── presentation/       # TUI implementation
├── pkg/
│   ├── totp/              # TOTP implementation
│   └── validator/         # Password validation
└── tests/                 # Test files
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) TUI framework
- TOTP implementation based on [RFC 6238](https://tools.ietf.org/html/rfc6238)
- Encryption using Go's standard crypto library

## Disclaimer

This is a personal password manager. While it implements strong encryption and security practices, use at your own risk. Always maintain encrypted backups of your vault.
