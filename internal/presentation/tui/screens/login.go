package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hambosto/passmanager/internal/presentation/tui/styles"
)

// LoginScreen represents the login/unlock screen
type LoginScreen struct {
	passwordInput textinput.Model
	isNewVault    bool
	confirmInput  textinput.Model
	step          int // 0 = password, 1 = confirm (for new vault)
	error         string
	vaultExists   bool
	width         int
	height        int
}

// NewLoginScreen creates a new login screen
func NewLoginScreen(vaultExists bool) *LoginScreen {
	// Create password input
	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter master password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.Focus()
	passwordInput.Width = 40

	// Create confirm input
	confirmInput := textinput.New()
	confirmInput.Placeholder = "Confirm master password"
	confirmInput.EchoMode = textinput.EchoPassword
	confirmInput.EchoCharacter = '•'
	confirmInput.Width = 40

	return &LoginScreen{
		passwordInput: passwordInput,
		confirmInput:  confirmInput,
		isNewVault:    !vaultExists,
		vaultExists:   vaultExists,
		step:          0,
	}
}

// Init initializes the screen
func (s *LoginScreen) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (s *LoginScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return s, tea.Quit

		case "ctrl+n":
			// Toggle new vault mode
			s.isNewVault = !s.isNewVault
			s.step = 0
			s.error = ""
			s.passwordInput.SetValue("")
			s.confirmInput.SetValue("")
			return s, nil

		case "ctrl+h":
			// Toggle password visibility
			if s.passwordInput.EchoMode == textinput.EchoPassword {
				s.passwordInput.EchoMode = textinput.EchoNormal
				s.confirmInput.EchoMode = textinput.EchoNormal
			} else {
				s.passwordInput.EchoMode = textinput.EchoPassword
				s.confirmInput.EchoMode = textinput.EchoPassword
			}
			return s, nil

		case "enter":
			if s.isNewVault {
				if s.step == 0 {
					// Move to confirm step
					s.step = 1
					s.passwordInput.Blur()
					s.confirmInput.Focus()
					return s, textinput.Blink
				} else {
					// Validate passwords match
					if s.passwordInput.Value() != s.confirmInput.Value() {
						s.error = "Passwords do not match"
						s.step = 0
						s.passwordInput.SetValue("")
						s.confirmInput.SetValue("")
						s.confirmInput.Blur()
						s.passwordInput.Focus()
						return s, textinput.Blink
					}
					// Validate password strength (minimum 8 characters)
					if len(s.passwordInput.Value()) < 8 {
						s.error = "Password must be at least 8 characters"
						s.step = 0
						s.passwordInput.SetValue("")
						s.confirmInput.SetValue("")
						s.confirmInput.Blur()
						s.passwordInput.Focus()
						return s, textinput.Blink
					}
					return s, func() tea.Msg {
						return UnlockMsg{Password: s.passwordInput.Value(), IsNew: true}
					}
				}
			} else {
				// Unlock existing vault
				return s, func() tea.Msg {
					return UnlockMsg{Password: s.passwordInput.Value(), IsNew: false}
				}
			}

		case "tab":
			if s.isNewVault && s.step == 1 {
				s.step = 0
				s.confirmInput.Blur()
				s.passwordInput.Focus()
				return s, textinput.Blink
			}
		}
	}

	// Update the active input
	if s.step == 0 {
		s.passwordInput, cmd = s.passwordInput.Update(msg)
	} else {
		s.confirmInput, cmd = s.confirmInput.Update(msg)
	}

	return s, cmd
}

// View renders the screen
func (s *LoginScreen) View() string {
	var b strings.Builder

	// Title
	title := styles.TitleStyle.Render(styles.IconLock + " Password Manager")
	b.WriteString(styles.CenterHorizontal(s.width, title))
	b.WriteString("\n\n")

	// Main box content
	var boxContent strings.Builder

	if s.isNewVault {
		boxContent.WriteString(lipgloss.NewStyle().Foreground(styles.Info).Bold(true).Render("Create New Vault"))
		boxContent.WriteString("\n\n")

		if s.step == 0 {
			boxContent.WriteString("Master Password:\n")
			boxContent.WriteString(s.passwordInput.View())
			boxContent.WriteString("\n\n")
			boxContent.WriteString(styles.HelpStyle.Render("Press Enter to continue"))
		} else {
			boxContent.WriteString("Confirm Master Password:\n")
			boxContent.WriteString(s.confirmInput.View())
			boxContent.WriteString("\n\n")
			boxContent.WriteString(styles.HelpStyle.Render("Press Enter to create vault • Tab to go back"))
		}
	} else {
		boxContent.WriteString("Unlock Vault\n\n")
		boxContent.WriteString("Master Password:\n")
		boxContent.WriteString(s.passwordInput.View())
		boxContent.WriteString("\n\n")
		boxContent.WriteString(styles.HelpStyle.Render("Press Enter to unlock"))
	}

	// Show error if any
	if s.error != "" {
		boxContent.WriteString("\n\n")
		boxContent.WriteString(styles.ErrorStyle.Render(styles.IconError + " " + s.error))
	}

	// Render box
	box := styles.BoxStyle.Width(50).Render(boxContent.String())
	b.WriteString(styles.CenterHorizontal(s.width, box))
	b.WriteString("\n\n")

	// Help text
	var helpText string
	if s.isNewVault {
		if s.vaultExists {
			helpText = "[Ctrl+N] Back to unlock  •  [Ctrl+H] Show/Hide  •  [Esc] Quit"
		} else {
			helpText = "[Ctrl+H] Show/Hide  •  [Esc] Quit"
		}
	} else {
		helpText = "[Ctrl+N] Create new vault  •  [Ctrl+H] Show/Hide  •  [Esc] Quit"
	}
	b.WriteString(styles.CenterHorizontal(s.width, styles.HelpStyle.Render(helpText)))

	// Center vertically
	content := b.String()
	lines := strings.Count(content, "\n") + 1
	paddingTop := (s.height - lines) / 2
	if paddingTop > 0 {
		content = strings.Repeat("\n", paddingTop) + content
	}

	return content
}

// UnlockMsg is sent when the user attempts to unlock the vault
type UnlockMsg struct {
	Password string
	IsNew    bool
}
