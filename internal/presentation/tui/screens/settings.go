package screens

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hambosto/passmanager/config"
	"github.com/hambosto/passmanager/internal/presentation/tui/styles"
	"github.com/hambosto/passmanager/internal/presentation/tui/util"
)

// SettingsScreen allows configuring application settings
type SettingsScreen struct {
	config *config.Config
	width  int
	height int

	// Form inputs
	autoLockInput    textinput.Model
	clipboardInput   textinput.Model
	passwordLenInput textinput.Model

	// State
	focusIndex int
	modified   bool

	// Checkboxes
	clearOnLock      bool
	clearOnExit      bool
	includeUpper     bool
	includeLower     bool
	includeNumbers   bool
	includeSymbols   bool
	excludeAmbiguous bool
}

// NewSettingsScreen creates a new settings screen
func NewSettingsScreen(cfg *config.Config) *SettingsScreen {
	autoLockInput := textinput.New()
	autoLockInput.Placeholder = "minutes"
	autoLockInput.Width = 10
	autoLockInput.SetValue(fmt.Sprintf("%d", cfg.Security.AutoLockTimeout))
	autoLockInput.Focus()

	clipboardInput := textinput.New()
	clipboardInput.Placeholder = "seconds"
	clipboardInput.Width = 10
	clipboardInput.SetValue(fmt.Sprintf("%d", cfg.Security.ClipboardTimeout))

	passwordLenInput := textinput.New()
	passwordLenInput.Placeholder = "characters"
	passwordLenInput.Width = 10
	passwordLenInput.SetValue(fmt.Sprintf("%d", cfg.PasswordGenerator.Length))

	return &SettingsScreen{
		config:           cfg,
		autoLockInput:    autoLockInput,
		clipboardInput:   clipboardInput,
		passwordLenInput: passwordLenInput,
		focusIndex:       0,
		modified:         false,

		// Initialize checkboxes from config
		clearOnLock:      cfg.Security.ClearClipboardOnLock,
		clearOnExit:      cfg.Security.ClearClipboardOnExit,
		includeUpper:     cfg.PasswordGenerator.IncludeUppercase,
		includeLower:     cfg.PasswordGenerator.IncludeLowercase,
		includeNumbers:   cfg.PasswordGenerator.IncludeNumbers,
		includeSymbols:   cfg.PasswordGenerator.IncludeSymbols,
		excludeAmbiguous: cfg.PasswordGenerator.ExcludeAmbiguous,
	}
}

// Init initializes the screen
func (s *SettingsScreen) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (s *SettingsScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back without saving
			return s, func() tea.Msg { return BackMsg{} }

		case "ctrl+c", "ctrl+q":
			return s, tea.Quit

		case "ctrl+s":
			// Save settings
			return s, s.saveSettings()

		case "tab", "shift+tab":
			// Navigate between inputs
			if msg.String() == "tab" {
				s.focusIndex++
			} else {
				s.focusIndex--
			}

			maxIndex := 9 // Total number of focusable items
			if s.focusIndex > maxIndex {
				s.focusIndex = 0
			} else if s.focusIndex < 0 {
				s.focusIndex = maxIndex
			}

			s.updateFocus()
			return s, textinput.Blink

		case " ", "enter":
			// Toggle checkboxes
			s.modified = true
			switch s.focusIndex {
			case 3:
				s.clearOnLock = !s.clearOnLock
			case 4:
				s.clearOnExit = !s.clearOnExit
			case 6:
				s.includeUpper = !s.includeUpper
			case 7:
				s.includeLower = !s.includeLower
			case 8:
				s.includeNumbers = !s.includeNumbers
			case 9:
				s.includeSymbols = !s.includeSymbols
			case 10:
				s.excludeAmbiguous = !s.excludeAmbiguous
			}
			return s, nil
		}
	}

	// Mark as modified when inputs change
	oldAutoLock := s.autoLockInput.Value()
	oldClipboard := s.clipboardInput.Value()
	oldPasswordLen := s.passwordLenInput.Value()

	// Update the focused input
	switch s.focusIndex {
	case 0:
		s.autoLockInput, cmd = s.autoLockInput.Update(msg)
		cmds = append(cmds, cmd)
	case 1:
		s.clipboardInput, cmd = s.clipboardInput.Update(msg)
		cmds = append(cmds, cmd)
	case 5:
		s.passwordLenInput, cmd = s.passwordLenInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	if oldAutoLock != s.autoLockInput.Value() ||
		oldClipboard != s.clipboardInput.Value() ||
		oldPasswordLen != s.passwordLenInput.Value() {
		s.modified = true
	}

	return s, tea.Batch(cmds...)
}

// View renders the screen
func (s *SettingsScreen) View() string {
	var b strings.Builder

	// Title
	title := styles.TitleStyle.Render(styles.IconKey + " Settings")
	b.WriteString(title)
	if s.modified {
		b.WriteString(lipgloss.NewStyle().Foreground(styles.Warning).Render(" (modified)"))
	}
	b.WriteString("\n\n")

	// Security section
	securityBox := s.renderSecuritySection()
	b.WriteString(securityBox)
	b.WriteString("\n\n")

	// Password generator section
	generatorBox := s.renderPasswordGeneratorSection()
	b.WriteString(generatorBox)
	b.WriteString("\n\n")

	// Help text
	helpText := "[Ctrl+S] Save  •  [Esc] Cancel  •  [Tab] Next  •  [Space] Toggle"
	b.WriteString(styles.HelpStyle.Render(helpText))

	return b.String()
}

// renderSecuritySection renders the security settings section
func (s *SettingsScreen) renderSecuritySection() string {
	var content strings.Builder

	content.WriteString(lipgloss.NewStyle().Bold(true).Render("Security"))
	content.WriteString("\n\n")

	// Auto-lock timeout
	content.WriteString(s.renderField(0, "Auto-lock timeout:", s.autoLockInput.View()+" minutes (0 = disabled)"))
	content.WriteString("\n\n")

	// Clipboard timeout
	content.WriteString(s.renderField(1, "Clipboard timeout:", s.clipboardInput.View()+" seconds"))
	content.WriteString("\n\n")

	// Clear clipboard options
	content.WriteString(lipgloss.NewStyle().Foreground(styles.Subtle).Render("Clear clipboard on:"))
	content.WriteString("\n")
	content.WriteString(s.renderCheckbox(3, "Lock", s.clearOnLock))
	content.WriteString("  ")
	content.WriteString(s.renderCheckbox(4, "Exit", s.clearOnExit))

	return styles.BoxStyle.Width(util.MinInt(70, s.width-4)).Render(content.String())
}

// renderPasswordGeneratorSection renders the password generator settings
func (s *SettingsScreen) renderPasswordGeneratorSection() string {
	var content strings.Builder

	content.WriteString(lipgloss.NewStyle().Bold(true).Render("Password Generator Defaults"))
	content.WriteString("\n\n")

	// Length
	content.WriteString(s.renderField(5, "Length:", s.passwordLenInput.View()+" characters"))
	content.WriteString("\n\n")

	// Include options
	content.WriteString(lipgloss.NewStyle().Foreground(styles.Subtle).Render("Include:"))
	content.WriteString("\n")
	content.WriteString(s.renderCheckbox(6, "Uppercase", s.includeUpper))
	content.WriteString("  ")
	content.WriteString(s.renderCheckbox(7, "Lowercase", s.includeLower))
	content.WriteString("\n")
	content.WriteString(s.renderCheckbox(8, "Numbers", s.includeNumbers))
	content.WriteString("  ")
	content.WriteString(s.renderCheckbox(9, "Symbols", s.includeSymbols))
	content.WriteString("\n\n")

	content.WriteString(s.renderCheckbox(10, "Exclude ambiguous (0,O,l,1,I)", s.excludeAmbiguous))

	return styles.BoxStyle.Width(util.MinInt(70, s.width-4)).Render(content.String())
}

// renderField renders a form field
func (s *SettingsScreen) renderField(index int, label, value string) string {
	labelStyle := lipgloss.NewStyle().Bold(true)
	if index == s.focusIndex {
		labelStyle = labelStyle.Foreground(styles.Primary)
	}

	return labelStyle.Render(label) + " " + value
}

// renderCheckbox renders a checkbox
func (s *SettingsScreen) renderCheckbox(index int, label string, checked bool) string {
	icon := "☐"
	if checked {
		icon = "☑"
	}

	style := lipgloss.NewStyle()
	if index == s.focusIndex {
		style = style.Foreground(styles.Primary).Bold(true)
	}

	return style.Render(icon + " " + label)
}

// updateFocus updates which input is focused
func (s *SettingsScreen) updateFocus() {
	s.autoLockInput.Blur()
	s.clipboardInput.Blur()
	s.passwordLenInput.Blur()

	switch s.focusIndex {
	case 0:
		s.autoLockInput.Focus()
	case 1:
		s.clipboardInput.Focus()
	case 5:
		s.passwordLenInput.Focus()
	}
}

// saveSettings creates a command to save the settings
func (s *SettingsScreen) saveSettings() tea.Cmd {
	// Update config from inputs
	if val, err := parseInt(s.autoLockInput.Value()); err == nil {
		s.config.Security.AutoLockTimeout = val
	}
	if val, err := parseInt(s.clipboardInput.Value()); err == nil {
		s.config.Security.ClipboardTimeout = val
	}
	if val, err := parseInt(s.passwordLenInput.Value()); err == nil {
		s.config.PasswordGenerator.Length = val
	}

	s.config.Security.ClearClipboardOnLock = s.clearOnLock
	s.config.Security.ClearClipboardOnExit = s.clearOnExit
	s.config.PasswordGenerator.IncludeUppercase = s.includeUpper
	s.config.PasswordGenerator.IncludeLowercase = s.includeLower
	s.config.PasswordGenerator.IncludeNumbers = s.includeNumbers
	s.config.PasswordGenerator.IncludeSymbols = s.includeSymbols
	s.config.PasswordGenerator.ExcludeAmbiguous = s.excludeAmbiguous

	return func() tea.Msg {
		return SaveSettingsMsg{Config: s.config}
	}
}

// SaveSettingsMsg signals that settings should be saved
type SaveSettingsMsg struct {
	Config *config.Config
}

// parseInt parses an integer from a string
func parseInt(s string) (int, error) {
	var val int
	_, err := fmt.Sscanf(s, "%d", &val)
	return val, err
}
