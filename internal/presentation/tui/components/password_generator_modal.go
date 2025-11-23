package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hambosto/passmanager/internal/application/service"
	"github.com/hambosto/passmanager/internal/presentation/tui/styles"
	"github.com/hambosto/passmanager/internal/presentation/tui/util"
	"github.com/hambosto/passmanager/pkg/validator"
)

// PasswordGeneratorModal is a modal for generating passwords
type PasswordGeneratorModal struct {
	visible bool
	width   int
	height  int

	// Generation mode
	usePassphrase bool

	// Password config
	passwordConfig service.PasswordConfig

	// Passphrase config
	passphraseConfig service.PassphraseConfig

	// Generated password
	password string
	entropy  float64
	strength validator.PasswordStrength

	// UI state
	focusedOption int
}

// NewPasswordGeneratorModal creates a new password generator modal
func NewPasswordGeneratorModal() *PasswordGeneratorModal {
	return &PasswordGeneratorModal{
		visible:          false,
		passwordConfig:   service.DefaultPasswordConfig(),
		passphraseConfig: service.DefaultPassphraseConfig(),
		usePassphrase:    false,
	}
}

// Show shows the modal
func (m *PasswordGeneratorModal) Show() {
	m.visible = true
	m.generatePassword()
}

// Hide hides the modal
func (m *PasswordGeneratorModal) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is visible
func (m *PasswordGeneratorModal) IsVisible() bool {
	return m.visible
}

// GetPassword returns the generated password
func (m *PasswordGeneratorModal) GetPassword() string {
	return m.password
}

// Update handles messages
func (m *PasswordGeneratorModal) Update(msg tea.Msg) tea.Cmd {
	if !m.visible {
		return nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.Hide()
			return nil

		case "enter":
			// Use the generated password
			m.Hide()
			return func() tea.Msg {
				return UsePasswordMsg{Password: m.password}
			}

		case "ctrl+r":
			// Regenerate
			m.generatePassword()
			return nil

		case "ctrl+c":
			// Copy to clipboard (handled by parent)
			return func() tea.Msg {
				return CopyPasswordMsg{Password: m.password}
			}

		case "tab":
			// Toggle between password/passphrase
			m.usePassphrase = !m.usePassphrase
			m.generatePassword()
			return nil

		case "up", "k":
			if m.focusedOption > 0 {
				m.focusedOption--
			}
			return nil

		case "down", "j":
			maxOptions := 4
			if m.usePassphrase {
				maxOptions = 2
			}
			if m.focusedOption < maxOptions {
				m.focusedOption++
			}
			return nil

		case "left", "h":
			m.adjustOption(-1)
			m.generatePassword()
			return nil

		case "right", "l":
			m.adjustOption(1)
			m.generatePassword()
			return nil
		}
	}

	return nil
}

// View renders the modal
func (m *PasswordGeneratorModal) View() string {
	if !m.visible {
		return ""
	}

	var content strings.Builder

	// Title
	content.WriteString(styles.TitleStyle.Render(styles.IconKey + " Password Generator"))
	content.WriteString("\n\n")

	// Mode selector
	modeText := "● Random Password    ○ Passphrase"
	if m.usePassphrase {
		modeText = "○ Random Password    ● Passphrase"
	}
	content.WriteString(lipgloss.NewStyle().Foreground(styles.Primary).Render(modeText))
	content.WriteString("\n\n")

	// Generated password box
	passwordBox := styles.BoxStyle.
		Width(util.MinInt(60, m.width-10)).
		Align(lipgloss.Center).
		Render(lipgloss.NewStyle().Bold(true).Foreground(styles.Success).Render(m.password))
	content.WriteString(passwordBox)
	content.WriteString("\n\n")

	// Strength meter
	strengthBar := m.renderStrengthMeter()
	content.WriteString(strengthBar)
	content.WriteString("\n")

	crackTime := validator.EstimateCrackTime(m.entropy)
	content.WriteString(lipgloss.NewStyle().Foreground(styles.Subtle).Render(
		fmt.Sprintf("Estimated crack time: %s", crackTime)))
	content.WriteString("\n\n")

	// Options
	optionsBox := m.renderOptions()
	content.WriteString(optionsBox)
	content.WriteString("\n\n")

	// Help
	helpText := "[Ctrl+R] Regenerate  •  [Ctrl+C] Copy  •  [Enter] Use  •  [Tab] Switch Mode  •  [Esc] Cancel"
	content.WriteString(styles.HelpStyle.Render(helpText))

	// Create modal box
	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Primary).
		Padding(2, 4).
		Width(util.MinInt(70, m.width-4)).
		Render(content.String())

	// Center the modal
	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modalBox,
	)

	return centered
}

// renderStrengthMeter renders the password strength meter
func (m *PasswordGeneratorModal) renderStrengthMeter() string {
	// Calculate percentage for progress bar
	percentage := m.entropy / 128.0
	if percentage > 1.0 {
		percentage = 1.0
	}

	var color lipgloss.Color
	var label string

	switch m.strength {
	case validator.StrengthWeak:
		color = styles.Danger
		label = "Weak"
	case validator.StrengthFair:
		color = styles.Warning
		label = "Fair"
	case validator.StrengthGood:
		color = styles.Info
		label = "Good"
	case validator.StrengthStrong:
		color = styles.Success
		label = "Strong"
	case validator.StrengthExcellent:
		color = styles.Success
		label = "Excellent"
	}

	// Create progress bar
	barWidth := 40
	filled := int(percentage * float64(barWidth))
	empty := barWidth - filled

	bar := lipgloss.NewStyle().Foreground(color).Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(styles.Gray700).Render(strings.Repeat("░", empty))

	text := fmt.Sprintf("Strength: %s (%.0f bits entropy)", label, m.entropy)

	return text + "\n" + bar
}

// renderOptions renders the configuration options
func (m *PasswordGeneratorModal) renderOptions() string {
	var content strings.Builder

	content.WriteString(lipgloss.NewStyle().Bold(true).Render("Options"))
	content.WriteString("\n\n")

	if m.usePassphrase {
		// Passphrase options
		content.WriteString(m.renderOption(0, fmt.Sprintf("Word Count: %d", m.passphraseConfig.WordCount)))
		content.WriteString("\n")
		content.WriteString(m.renderOption(1, fmt.Sprintf("Separator: %s", m.passphraseConfig.Separator)))
		content.WriteString("\n")

		capIcon := "☐"
		if m.passphraseConfig.Capitalize {
			capIcon = "☑"
		}
		content.WriteString(m.renderOption(2, capIcon+" Capitalize"))
	} else {
		// Password options
		content.WriteString(m.renderOption(0, fmt.Sprintf("Length: %d", m.passwordConfig.Length)))
		content.WriteString("\n")

		upperIcon := "☐"
		if m.passwordConfig.IncludeUpper {
			upperIcon = "☑"
		}
		content.WriteString(m.renderOption(1, upperIcon+" Uppercase"))
		content.WriteString("\n")

		numbersIcon := "☐"
		if m.passwordConfig.IncludeNumbers {
			numbersIcon = "☑"
		}
		content.WriteString(m.renderOption(2, numbersIcon+" Numbers"))
		content.WriteString("\n")

		symbolsIcon := "☐"
		if m.passwordConfig.IncludeSymbols {
			symbolsIcon = "☑"
		}
		content.WriteString(m.renderOption(3, symbolsIcon+" Symbols"))
		content.WriteString("\n")

		ambigIcon := "☐"
		if m.passwordConfig.ExcludeAmbiguous {
			ambigIcon = "☑"
		}
		content.WriteString(m.renderOption(4, ambigIcon+" Exclude Ambiguous"))
	}

	content.WriteString("\n")
	content.WriteString(styles.HelpStyle.Render("[↑↓] Navigate  [←→] Adjust"))

	return content.String()
}

// renderOption renders a single option
func (m *PasswordGeneratorModal) renderOption(index int, text string) string {
	style := lipgloss.NewStyle()
	if index == m.focusedOption {
		style = style.Foreground(styles.Primary).Bold(true)
		text = "> " + text
	} else {
		text = "  " + text
	}
	return style.Render(text)
}

// adjustOption adjusts the focused option's value
func (m *PasswordGeneratorModal) adjustOption(delta int) {
	if m.usePassphrase {
		switch m.focusedOption {
		case 0: // Word count
			m.passphraseConfig.WordCount += delta
			if m.passphraseConfig.WordCount < 3 {
				m.passphraseConfig.WordCount = 3
			}
			if m.passphraseConfig.WordCount > 10 {
				m.passphraseConfig.WordCount = 10
			}
		case 1: // Separator
			separators := []string{"-", "_", " ", ".", ""}
			currentIndex := 0
			for i, sep := range separators {
				if sep == m.passphraseConfig.Separator {
					currentIndex = i
					break
				}
			}
			currentIndex += delta
			if currentIndex < 0 {
				currentIndex = len(separators) - 1
			}
			if currentIndex >= len(separators) {
				currentIndex = 0
			}
			m.passphraseConfig.Separator = separators[currentIndex]
		case 2: // Capitalize
			m.passphraseConfig.Capitalize = !m.passphraseConfig.Capitalize
		}
	} else {
		switch m.focusedOption {
		case 0: // Length
			m.passwordConfig.Length += delta
			if m.passwordConfig.Length < 8 {
				m.passwordConfig.Length = 8
			}
			if m.passwordConfig.Length > 128 {
				m.passwordConfig.Length = 128
			}
		case 1: // Uppercase
			m.passwordConfig.IncludeUpper = !m.passwordConfig.IncludeUpper
		case 2: // Numbers
			m.passwordConfig.IncludeNumbers = !m.passwordConfig.IncludeNumbers
		case 3: // Symbols
			m.passwordConfig.IncludeSymbols = !m.passwordConfig.IncludeSymbols
		case 4: // Exclude ambiguous
			m.passwordConfig.ExcludeAmbiguous = !m.passwordConfig.ExcludeAmbiguous
		}
	}
}

// generatePassword generates a new password based on current config
func (m *PasswordGeneratorModal) generatePassword() {
	var password string
	var err error

	if m.usePassphrase {
		password, err = service.GeneratePassphrase(m.passphraseConfig)
	} else {
		password, err = service.GeneratePassword(m.passwordConfig)
	}

	if err != nil {
		password = "Error generating password"
	}

	m.password = password
	m.entropy = service.CalculatePasswordEntropy(password)
	m.strength = service.GetPasswordStrength(password)
}

// UsePasswordMsg signals to use the generated password
type UsePasswordMsg struct {
	Password string
}

// CopyPasswordMsg signals to copy the password
type CopyPasswordMsg struct {
	Password string
}
