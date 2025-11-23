# Password Manager - User Guide

## Table of Contents
1. [Getting Started](#getting-started)
2. [Creating Your Vault](#creating-your-vault)
3. [Managing Entries](#managing-entries)
4. [TOTP (2FA) Setup](#totp-2fa-setup)
5. [Password Generation](#password-generation)
6. [Settings](#settings)
7. [Security Best Practices](#security-best-practices)
8. [Keyboard Shortcuts](#keyboard-shortcuts)
9. [Troubleshooting](#troubleshooting)

## Getting Started

### Installation

```bash
# Build from source
cd passmanager
make build

# Run the application
./build/passmanager
```

The vault file will be created at `~/.config/passmanager/vault.enc`

### System Requirements
- Go 1.21 or later (for building)
- Linux, macOS, or Windows
- Terminal with Unicode support

## Creating Your Vault

### First Launch

When you first run the application, you'll see the login screen:

1. Press `Ctrl+N` to create a new vault
2. Enter a strong master password (minimum 8 characters)
3. Confirm your master password
4. Your encrypted vault is created!

**⚠️ Important**: Write down your master password! There is no password recovery.

### Master Password Tips

- Use at least 12 characters
- Include uppercase, lowercase, numbers, and symbols
- Use a passphrase you can remember
- Don't reuse passwords from other services
- Consider using the built-in password generator

## Managing Entries

### Creating an Entry

1. From the vault list, press `Ctrl+N`
2. Fill in the entry details:
   - **Name**: Entry title (required)
   - **Username**: Email or username
   - **Password**: Use `Ctrl+G` to generate
   - **Website**: URL (optional)
   - **TOTP**: 2FA secret (optional)
   - **Notes**: Additional information
3. Press `Ctrl+S` to save

### Viewing Entry Details

1. Navigate to an entry with ↑/↓ or j/k
2. Press `Enter` to view details
3. Use shortcuts to copy:
   - `Ctrl+U` - Copy username
   - `Ctrl+P` - Copy password
   - `Ctrl+T` - Copy TOTP code

### Editing an Entry

1. View the entry
2. Press `Ctrl+E` to edit
3. Make your changes
4. Press `Ctrl+S` to save

### Deleting an Entry

Currently done by editing and removing all content. Future versions will have explicit delete.

### Organizing with Favorites

- Press `Space` in the vault list to toggle favorite
- Or press `Ctrl+F` in the entry editor
- Favorites appear with a ⭐ icon

### Searching Entries

1. Press `/` in the vault list
2. Type your search query
3. Results filter in real-time

## TOTP (2FA) Setup

### What is TOTP?

TOTP (Time-based One-Time Password) provides two-factor authentication codes that change every 30 seconds, like Google Authenticator.

### Adding TOTP to an Entry

**Method 1: Using otpauth:// URI**
1. When setting up 2FA on a website, choose "Manual setup" or "Can't scan QR code"
2. Copy the `otpauth://totp/...` URI
3. Paste into the TOTP field in your entry

**Method 2: Using Base32 Secret**
1. When setting up 2FA, copy the secret key (usually shown as backup)
2. Enter it in the TOTP field

### Viewing TOTP Codes

1. Open an entry with TOTP configured
2. See the 6-digit code with countdown timer
3. Press `Ctrl+T` to copy the code
4. Paste it into the 2FA prompt

### TOTP Best Practices

- ✅ Store TOTP secrets in the password manager
- ✅ Save backup codes separately
- ✅ Test the TOTP before finishing 2FA setup
-❌ Don't share TOTP secrets
- ❌ Don't screenshot QR codes

## Password Generation

### Generating Passwords

1. In the entry editor, press `Ctrl+G`
2. The password generator modal opens
3. Configure options:
   - **Length**: 8-128 characters
   - **Character sets**: Upper, lower, numbers, symbols
   - **Exclude ambiguous**: Remove 0, O, l, 1, I
4. Press `Enter` to use the password
5. Or `Ctrl+C` to copy without using

### Generating Passphrases

1. Open password generator (`Ctrl+G`)
2. Press `Tab` to switch to passphrase mode
3. Configure options:
   - **Word count**: 3-10 words
   - **Separator**: -, _, space, etc.
   - **Capitalize**: Uppercase first letters
   - **Include number**: Add random number
4. Example: `Correct-Horse-Battery-Staple-42`

### Password Strength Meter

The generator shows:
- **Strength**: Very Weak / Weak / Fair / Strong / Very Strong
- **Entropy**: Bits of randomness
- **Crack time**: Estimated time to crack

Aim for "Strong" or "Very Strong" passwords.

## Settings

Access settings with `Ctrl+,` from the vault list.

### Security Settings

**Auto-lock timeout**: Minutes of inactivity before auto-lock (0 = disabled)
- Recommended: 5-15 minutes
- Clears master key from memory

**Clipboard timeout**: Seconds before clipboard auto-clears
- Default: 30 seconds
- Protects against clipboard hijacking

**Clear clipboard on lock**: Automatically clear clipboard when vault locks
**Clear clipboard on exit**: Clear clipboard when quitting

### Password Generator Defaults

Configure default settings for password generation:
- Length
- Character sets
- Exclude ambiguous

## Security Best Practices

### Master Password
✅ Use a strong, unique master password
✅ Store it in a secure location
✅ Change it periodically (annually)
❌ Never share your master password
❌ Don't use the same password elsewhere

### Vault Security
✅ Enable auto-lock
✅ Lock before leaving computer
✅ Use full disk encryption on your device
✅ Regular backups (export and encrypt)
❌ Don't store vault in cloud without extra encryption

### Entry Management
✅ Use unique passwords for every site
✅ Enable 2FA where available
✅ Store TOTP secrets in vault
✅ Update passwords after breaches
❌ Don't reuse passwords

### Physical Security
✅ Lock your computer when away
✅ Be aware of shoulder surfing
✅ Don't leave clipboard with passwords
✅ Use auto-lock feature

## Keyboard Shortcuts

### Global
- `Ctrl+Q` - Quit application
- `Ctrl+L` - Lock vault
- `?` - Show help
- `Ctrl+,` - Settings
- `Esc` - Go back / Cancel

### Vault List
- `↑↓` or `k/j` - Navigate entries
- `Enter` - View entry details
- `/` - Search
- `Ctrl+N` - New entry
- `Space` - Toggle favorite

### Entry Detail
- `Ctrl+U` - Copy username
- `Ctrl+P` - Copy password
- `Ctrl+T` - Copy TOTP code
- `Ctrl+H` - Show/hide password
- `Ctrl+E` - Edit entry

### Entry Editor
- `Tab` - Next field
- `Ctrl+S` - Save
- `Ctrl+G` - Generate password
- `Ctrl+F` - Toggle favorite
- `Ctrl+H` - Show/hide password
- `Esc` - Cancel

### Password Generator
- `Tab` - Switch mode (password/passphrase)
- `↑↓` - Navigate options
- `←→` - Adjust values
- `Ctrl+R` - Regenerate
- `Enter` - Use password
- `Ctrl+C` - Copy password
- `Esc` - Close

## Troubleshooting

### "Failed to unlock vault (wrong password?)"
- Double-check your master password
- Ensure Caps Lock is off
- If forgotten, vault cannot be recovered

### Clipboard not auto-clearing
- Check clipboard timeout in settings
- Some clipboard managers may interfere
- Manually clear with system clipboard manager

### TOTP codes not working
- Verify system clock is accurate
- Check the TOTP secret is correct
- Ensure 30-second period (default)
- Try re-entering the secret

### Vault file not found
- Default location: `~/.config/passmanager/vault.enc`
- Check file permissions
- Ensure directory exists

### Build errors
- Verify Go version (1.21+)
- Run `go mod tidy`
- Check `make deps`

## CLI Usage

For automation and scripts:

```bash
# List all entries
passmanager list

# Get specific entry
passmanager get "GitHub"

# Generate password
passmanager generate --length 20

# Show version
passmanager version
```

## Backup and Export

**Manual Backup:**
1. Copy `~/.config/passmanager/vault.enc`
2. Store encrypted copy securely
3. Test restore periodically

**Future Features:**
- CSV export with warnings
- Encrypted backup files
- Cloud sync options

## Getting Help

- Check `README.md` for feature overview
- Review `SECURITY.md` for security details
- Read `ARCHITECTURE.md` for technical docs
- Press `?` in the app for shortcuts
- File issues on GitHub for bugs

## Version History

**v1.0.0** - Initial release
- Full vault management
- TOTP support
- Password generation
- TUI with all screens
- Auto-lock
- Settings configuration
