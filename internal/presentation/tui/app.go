package tui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/hambosto/passmanager/config"
	"github.com/hambosto/passmanager/internal/domain/entity"
	"github.com/hambosto/passmanager/internal/infrastructure"
	"github.com/hambosto/passmanager/internal/infrastructure/clipboard"
	"github.com/hambosto/passmanager/internal/infrastructure/crypto"
	"github.com/hambosto/passmanager/internal/infrastructure/storage"
	"github.com/hambosto/passmanager/internal/presentation/tui/components"
	"github.com/hambosto/passmanager/internal/presentation/tui/screens"
)

// Screen represents different screens in the app
type Screen int

const (
	ScreenLogin Screen = iota
	ScreenVaultList
	ScreenEntryDetail
	ScreenEntryEditor
	ScreenSettings
)

// App is the main TUI application model
type App struct {
	currentScreen  Screen
	previousScreen Screen

	// Screens
	loginScreen    *screens.LoginScreen
	vaultList      *screens.VaultListScreen
	entryDetail    *screens.EntryDetailScreen
	entryEditor    *screens.EntryEditorScreen
	settingsScreen *screens.SettingsScreen
	helpScreen     *screens.HelpScreen

	// Components
	passwordGenerator *components.PasswordGeneratorModal

	// Infrastructure
	autoLocker *infrastructure.AutoLocker

	// Vault state
	vault      *entity.Vault
	vaultPath  string
	masterKey  []byte
	repository *storage.FileRepository
	clipboard  *clipboard.Manager
	config     *config.Config

	// Window size
	width  int
	height int

	// Error and messages
	err     error
	message string
}

// NewApp creates a new TUI application
func NewApp(vaultPath string, clipboardTimeout int) *App {
	repo := storage.NewFileRepository(vaultPath)
	vaultExists := repo.Exists()

	cfg := config.DefaultConfig()

	return &App{
		currentScreen:     ScreenLogin,
		loginScreen:       screens.NewLoginScreen(vaultExists),
		vaultPath:         vaultPath,
		repository:        repo,
		clipboard:         clipboard.NewManager(time.Duration(clipboardTimeout) * time.Second),
		passwordGenerator: components.NewPasswordGeneratorModal(),
		config:            cfg,
	}
}

// NewAppWithConfig creates a new TUI application with config
func NewAppWithConfig(vaultPath string, cfg *config.Config) *App {
	repo := storage.NewFileRepository(vaultPath)
	vaultExists := repo.Exists()

	return &App{
		currentScreen:     ScreenLogin,
		loginScreen:       screens.NewLoginScreen(vaultExists),
		vaultPath:         vaultPath,
		repository:        repo,
		clipboard:         clipboard.NewManager(time.Duration(cfg.Security.ClipboardTimeout) * time.Second),
		passwordGenerator: components.NewPasswordGeneratorModal(),
		config:            cfg,
	}
}

// Init initializes the application
func (a *App) Init() tea.Cmd {
	return a.loginScreen.Init()
}

// Update handles messages
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle password generator first if visible
	if a.passwordGenerator.IsVisible() {
		cmd := a.passwordGenerator.Update(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}

		// Don't pass through most messages if modal is open
		switch msg.(type) {
		case components.UsePasswordMsg, components.CopyPasswordMsg:
			// Pass through these messages
		case tea.WindowSizeMsg:
			// Pass through window size
		default:
			return a, tea.Batch(cmds...)
		}
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height

	case tea.KeyMsg:
		// Global shortcuts
		switch msg.String() {
		case "ctrl+c":
			a.clipboard.Clear()
			return a, tea.Quit

		case "ctrl+q":
			a.clipboard.Clear()
			return a, tea.Quit
		}

	case screens.UnlockMsg:
		// Handle vault unlock
		return a.handleUnlock(msg)

	case screens.BackMsg:
		// Go back to previous screen
		if a.currentScreen == ScreenEntryDetail {
			a.currentScreen = ScreenVaultList
			return a, nil
		}

	case screens.EditEntryMsg:
		// Switch to entry editor
		a.entryEditor = screens.NewEntryEditorScreen(msg.Entry, false)
		a.previousScreen = a.currentScreen
		a.currentScreen = ScreenEntryEditor
		return a, a.entryEditor.Init()

	case screens.CancelEditMsg:
		// Cancel editing, go back
		a.currentScreen = a.previousScreen
		return a, nil

	case screens.SaveEntryMsg:
		// Save entry
		return a.handleSaveEntry(msg)

	case screens.SaveSettingsMsg:
		// Save settings to config file
		if err := msg.Config.Save(config.GetConfigPath()); err != nil {
			a.err = fmt.Errorf("failed to save settings: %w", err)
		} else {
			a.message = "Settings saved!"
			// Update clipboard timeout
			a.clipboard = clipboard.NewManager(time.Duration(msg.Config.Security.ClipboardTimeout) * time.Second)
		}
		a.currentScreen = ScreenVaultList
		return a, nil

	case screens.NewEntryMsg:
		// Create new entry
		newEntry := entity.NewEntry(entity.EntryTypeLogin, "")
		a.entryEditor = screens.NewEntryEditorScreen(newEntry, true)
		a.previousScreen = a.currentScreen
		a.currentScreen = ScreenEntryEditor
		return a, a.entryEditor.Init()

	case screens.OpenPasswordGeneratorMsg:
		// Show password generator modal
		a.passwordGenerator.Show()
		return a, nil

	case components.UsePasswordMsg:
		// Use generated password in editor
		if a.entryEditor != nil && a.currentScreen == ScreenEntryEditor {
			// Set the password in the editor
			a.entryEditor.SetPassword(msg.Password)
			a.message = "Password set!"
		}
		return a, nil

	case components.CopyPasswordMsg:
		// Copy generated password
		a.clipboard.CopyWithTimeout(msg.Password)
		a.message = "Password copied!"
		return a, nil

	case errMsg:
		a.err = msg.err
		return a, nil
	}

	// Route to current screen
	var cmd tea.Cmd
	switch a.currentScreen {
	case ScreenLogin:
		_, cmd = a.loginScreen.Update(msg)
		cmds = append(cmds, cmd)

	case ScreenVaultList:
		if a.vaultList != nil {
			// Handle keyboard shortcuts
			if keyMsg, ok := msg.(tea.KeyMsg); ok {
				switch keyMsg.String() {
				case "enter":
					selectedEntry := a.vaultList.GetSelectedEntry()
					if selectedEntry != nil {
						a.entryDetail = screens.NewEntryDetailScreen(selectedEntry, a.clipboard)
						a.currentScreen = ScreenEntryDetail
						return a, a.entryDetail.Init()
					}
				case "ctrl+,":
					// Open settings
					a.settingsScreen = screens.NewSettingsScreen(a.config)
					a.previousScreen = a.currentScreen
					a.currentScreen = ScreenSettings
					return a, a.settingsScreen.Init()
				case "?":
					// Open help
					a.helpScreen = screens.NewHelpScreen()
					a.previousScreen = a.currentScreen
					a.currentScreen = ScreenSettings
					return a, a.helpScreen.Init()
				}
			}
			_, cmd = a.vaultList.Update(msg)
			cmds = append(cmds, cmd)
		}

	case ScreenEntryDetail:
		if a.entryDetail != nil {
			_, cmd = a.entryDetail.Update(msg)
			cmds = append(cmds, cmd)
		}

	case ScreenEntryEditor:
		if a.entryEditor != nil {
			_, cmd = a.entryEditor.Update(msg)
			cmds = append(cmds, cmd)
		}

	case ScreenSettings:
		if a.settingsScreen != nil {
			_, cmd = a.settingsScreen.Update(msg)
			cmds = append(cmds, cmd)
		} else if a.helpScreen != nil {
			_, cmd = a.helpScreen.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return a, tea.Batch(cmds...)
}

// View renders the application
func (a *App) View() string {
	if a.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress Ctrl+C to quit", a.err)
	}

	var view string

	// Render current screen
	switch a.currentScreen {
	case ScreenLogin:
		view = a.loginScreen.View()

	case ScreenVaultList:
		if a.vaultList != nil {
			view = a.vaultList.View()
		}

	case ScreenEntryDetail:
		if a.entryDetail != nil {
			view = a.entryDetail.View()
		}

	case ScreenEntryEditor:
		if a.entryEditor != nil {
			view = a.entryEditor.View()
		}

	case ScreenSettings:
		if a.settingsScreen != nil {
			view = a.settingsScreen.View()
		} else if a.helpScreen != nil {
			view = a.helpScreen.View()
		}

	default:
		view = "Loading..."
	}

	// Overlay password generator if visible
	if a.passwordGenerator.IsVisible() {
		view = a.passwordGenerator.View()
	}

	return view
}

// handleUnlock handles the unlock message
func (a *App) handleUnlock(msg screens.UnlockMsg) (tea.Model, tea.Cmd) {
	if msg.IsNew {
		return a.createVault(msg.Password)
	}
	return a.unlockVault(msg.Password)
}

// createVault creates a new vault
func (a *App) createVault(password string) (tea.Model, tea.Cmd) {
	// Derive encryption key
	params, err := crypto.DefaultKeyDerivationParams()
	if err != nil {
		return a, func() tea.Msg {
			return errMsg{err: fmt.Errorf("failed to create key params: %w", err)}
		}
	}

	key := crypto.DeriveKey(password, params)
	a.masterKey = key

	// Create new vault
	vault := entity.NewVault()
	a.vault = vault

	// Save vault with KDF params
	if err := a.repository.Save(vault, key, params); err != nil {
		return a, func() tea.Msg {
			return errMsg{err: fmt.Errorf("failed to save vault: %w", err)}
		}
	}

	// Switch to vault list screen
	a.currentScreen = ScreenVaultList
	a.vaultList = screens.NewVaultListScreen(vault, a.clipboard)

	return a, a.vaultList.Init()
}

// unlockVault unlocks an existing vault
func (a *App) unlockVault(password string) (tea.Model, tea.Cmd) {
	// Load KDF params from vault file
	params, err := a.repository.LoadParams()
	if err != nil {
		return a, func() tea.Msg {
			return errMsg{err: fmt.Errorf("failed to load vault params: %w", err)}
		}
	}

	// Derive key with the SAME parameters used during vault creation
	key := crypto.DeriveKey(password, params)

	// Try to load vault
	vault, err := a.repository.Load(key)
	if err != nil {
		return a, func() tea.Msg {
			return errMsg{err: fmt.Errorf("failed to unlock vault (wrong password?): %w", err)}
		}
	}

	a.masterKey = key
	a.vault = vault

	// Initialize auto-locker if configured
	if a.config != nil && a.config.Security.AutoLockTimeout > 0 {
		timeout := time.Duration(a.config.Security.AutoLockTimeout) * time.Minute
		a.autoLocker = infrastructure.NewAutoLocker(timeout, func() tea.Msg {
			return infrastructure.AutoLockMsg{}
		})
	}

	// Switch to vault list screen
	a.currentScreen = ScreenVaultList
	a.vaultList = screens.NewVaultListScreen(vault, a.clipboard)

	return a, a.vaultList.Init()
}

// handleSaveEntry handles saving an entry
func (a *App) handleSaveEntry(msg screens.SaveEntryMsg) (tea.Model, tea.Cmd) {
	if msg.IsNew {
		// Add new entry
		a.vault.AddEntry(msg.Entry)
	} else {
		// Entry is already updated in place
		a.vault.Update()
	}

	// Load existing KDF params from file
	params, err := a.repository.LoadParams()
	if err != nil {
		a.err = fmt.Errorf("failed to load vault params: %w", err)
		return a, nil
	}

	// Save vault to disk with the existing KDF params
	if err := a.repository.Save(a.vault, a.masterKey, params); err != nil {
		a.err = fmt.Errorf("failed to save vault: %w", err)
		return a, nil
	}

	// Go back to vault list
	a.currentScreen = ScreenVaultList
	a.vaultList = screens.NewVaultListScreen(a.vault, a.clipboard)
	a.message = "Entry saved!"

	return a, a.vaultList.Init()
}

// errMsg represents an error message
type errMsg struct {
	err error
}
