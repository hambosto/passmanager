package screens

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hambosto/passmanager/internal/domain/entity"
	"github.com/hambosto/passmanager/internal/infrastructure/clipboard"
	"github.com/hambosto/passmanager/internal/presentation/tui/styles"
	"github.com/hambosto/passmanager/internal/presentation/tui/util"
	"github.com/hambosto/passmanager/pkg/totp"
)

// EntryDetailScreen shows the details of a single entry
type EntryDetailScreen struct {
	entry     *entity.Entry
	clipboard *clipboard.Manager
	width     int
	height    int

	// TOTP state
	totpCode      string
	totpExpiresIn time.Duration
	ticker        *time.Ticker

	// UI state
	showPassword bool
	copyMessage  string
	copyTimer    *time.Timer
}

// NewEntryDetailScreen creates a new entry detail screen
func NewEntryDetailScreen(entry *entity.Entry, clipboardMgr *clipboard.Manager) *EntryDetailScreen {
	return &EntryDetailScreen{
		entry:        entry,
		clipboard:    clipboardMgr,
		ticker:       time.NewTicker(1 * time.Second),
		showPassword: false,
	}
}

// Init initializes the screen
func (s *EntryDetailScreen) Init() tea.Cmd {
	// Update access time
	s.entry.UpdateAccessTime()

	// Generate initial TOTP code if available
	if s.entry.TOTPSecret != "" {
		s.updateTOTP()
	}

	return s.tickCmd()
}

// Update handles messages
func (s *EntryDetailScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back to vault list
			return s, func() tea.Msg { return BackMsg{} }

		case "ctrl+c", "ctrl+q":
			return s, tea.Quit

		case "ctrl+h":
			// Toggle password visibility
			s.showPassword = !s.showPassword
			return s, nil

		case "ctrl+u":
			// Copy username
			if s.entry.Username != "" {
				s.clipboard.CopyWithTimeout(s.entry.Username)
				s.showCopyMessage("Username copied!")
				return s, s.clearCopyMessageCmd()
			}

		case "ctrl+p", "ctrl+shift+c":
			// Copy password
			if s.entry.Password != "" {
				s.clipboard.CopyWithTimeout(s.entry.Password)
				s.showCopyMessage("Password copied!")
				return s, s.clearCopyMessageCmd()
			}

		case "ctrl+t":
			// Copy TOTP code
			if s.totpCode != "" {
				s.clipboard.CopyWithTimeout(s.totpCode)
				s.showCopyMessage("TOTP code copied!")
				return s, s.clearCopyMessageCmd()
			}

		case "ctrl+e":
			// Edit entry
			return s, func() tea.Msg { return EditEntryMsg{Entry: s.entry} }

		case "ctrl+o":
			// Open URL in browser
			if s.entry.URI != "" {
				// Note: Browser launch requires system command integration
				// Showing message for now as this would use exec.Command
				s.showCopyMessage("URL: " + s.entry.URI)
				return s, s.clearCopyMessageCmd()
			}
		}

	case tickMsg:
		// Update TOTP code
		if s.entry.TOTPSecret != "" {
			s.updateTOTP()
		}
		return s, s.tickCmd()

	case clearCopyMsgMsg:
		s.copyMessage = ""
		return s, nil
	}

	return s, nil
}

// View renders the screen
func (s *EntryDetailScreen) View() string {
	var b strings.Builder

	// Title
	icon := s.getEntryIcon()
	title := styles.TitleStyle.Render(icon + " " + s.entry.Name)
	if s.entry.IsFavorite {
		title = styles.FavoriteStyle.Render(styles.IconStar) + " " + title
	}
	b.WriteString(title)
	b.WriteString("\n\n")

	// Type and folder
	typeStr := lipgloss.NewStyle().Foreground(styles.Subtle).Render(
		fmt.Sprintf("Type: %s", s.entry.Type.String()))
	b.WriteString(typeStr)
	b.WriteString("\n\n")

	// Credentials box
	if s.entry.Type == entity.EntryTypeLogin {
		credBox := s.renderCredentials()
		b.WriteString(credBox)
		b.WriteString("\n\n")
	}

	// TOTP box (if available)
	if s.entry.TOTPSecret != "" {
		totpBox := s.renderTOTP()
		b.WriteString(totpBox)
		b.WriteString("\n\n")
	}

	// Notes box (if available)
	if s.entry.Notes != "" {
		notesBox := s.renderNotes()
		b.WriteString(notesBox)
		b.WriteString("\n\n")
	}

	// Custom fields (if available)
	if len(s.entry.CustomFields) > 0 {
		customBox := s.renderCustomFields()
		b.WriteString(customBox)
		b.WriteString("\n\n")
	}

	// Timestamps
	timestamps := s.renderTimestamps()
	b.WriteString(timestamps)
	b.WriteString("\n\n")

	// Copy message
	if s.copyMessage != "" {
		msg := styles.SuccessStyle.Render(styles.IconSuccess + " " + s.copyMessage)
		b.WriteString(msg)
		b.WriteString("\n\n")
	}

	// Help text
	helpText := "[Esc] Back  •  [Ctrl+U] Copy Username  •  [Ctrl+P] Copy Password"
	if s.entry.TOTPSecret != "" {
		helpText += "  •  [Ctrl+T] Copy TOTP"
	}
	helpText += "  •  [Ctrl+H] Show/Hide  •  [Ctrl+E] Edit"
	b.WriteString(styles.HelpStyle.Render(helpText))

	return b.String()
}

// renderCredentials renders the credentials box
func (s *EntryDetailScreen) renderCredentials() string {
	var content strings.Builder

	// Username
	if s.entry.Username != "" {
		content.WriteString(lipgloss.NewStyle().Bold(true).Render("Username:"))
		content.WriteString("\n")
		content.WriteString(s.entry.Username)
		content.WriteString("\n\n")
	}

	// Password
	if s.entry.Password != "" {
		content.WriteString(lipgloss.NewStyle().Bold(true).Render("Password:"))
		content.WriteString("\n")
		if s.showPassword {
			content.WriteString(s.entry.Password)
		} else {
			content.WriteString(strings.Repeat("•", 16))
		}
		content.WriteString("\n")
	}

	// Website
	if s.entry.URI != "" {
		content.WriteString("\n")
		content.WriteString(lipgloss.NewStyle().Bold(true).Render("Website:"))
		content.WriteString("\n")
		content.WriteString(lipgloss.NewStyle().Foreground(styles.Info).Render(s.entry.URI))
	}

	return styles.BoxStyle.
		Width(util.MinInt(60, s.width-4)).
		Render(content.String())
}

// renderTOTP renders the TOTP code box
func (s *EntryDetailScreen) renderTOTP() string {
	var content strings.Builder

	// Format code with space in middle (123 456)
	formattedCode := s.totpCode
	if len(s.totpCode) == 6 {
		formattedCode = s.totpCode[:3] + " " + s.totpCode[3:]
	}

	// TOTP code (large)
	codeStyle := lipgloss.NewStyle().
		Foreground(styles.Success).
		Bold(true).
		Align(lipgloss.Center)

	content.WriteString(codeStyle.Render(formattedCode))
	content.WriteString("\n\n")

	// Expiry info
	expiryText := fmt.Sprintf("%s Expires in: %ds", styles.IconClock, int(s.totpExpiresIn.Seconds()))
	content.WriteString(lipgloss.NewStyle().Foreground(styles.Subtle).Render(expiryText))
	content.WriteString("\n")

	// Progress bar
	percentage := float64(s.totpExpiresIn.Seconds()) / 30.0
	progressBar := styles.RenderProgressBar(percentage, 40)
	content.WriteString(progressBar)

	title := lipgloss.NewStyle().
		Foreground(styles.Success).
		Bold(true).
		Render("TOTP Code")

	boxContent := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.Success).
		Padding(1, 2).
		Width(util.MinInt(50, s.width-4)).
		Render(content.String())

	return title + "\n" + boxContent
}

// renderNotes renders the notes box
func (s *EntryDetailScreen) renderNotes() string {
	title := lipgloss.NewStyle().Bold(true).Render("Notes")

	notesContent := s.entry.Notes
	if len(notesContent) > 200 {
		notesContent = notesContent[:200] + "..."
	}

	box := styles.BoxStyle.
		Width(util.MinInt(60, s.width-4)).
		Render(notesContent)

	return title + "\n" + box
}

// renderCustomFields renders custom fields
func (s *EntryDetailScreen) renderCustomFields() string {
	var content strings.Builder

	for key, value := range s.entry.CustomFields {
		content.WriteString(lipgloss.NewStyle().Bold(true).Render(key + ":"))
		content.WriteString("\n")
		content.WriteString(value)
		content.WriteString("\n\n")
	}

	title := lipgloss.NewStyle().Bold(true).Render("Custom Fields")

	box := styles.BoxStyle.
		Width(util.MinInt(60, s.width-4)).
		Render(strings.TrimSpace(content.String()))

	return title + "\n" + box
}

// renderTimestamps renders the timestamps
func (s *EntryDetailScreen) renderTimestamps() string {
	created := s.entry.CreatedAt.Format("2006-01-02 15:04")
	updated := s.entry.UpdatedAt.Format("2006-01-02 15:04")

	text := fmt.Sprintf("Created: %s  •  Modified: %s", created, updated)

	if !s.entry.AccessedAt.IsZero() {
		accessed := s.entry.AccessedAt.Format("2006-01-02 15:04")
		text += fmt.Sprintf("  •  Last Accessed: %s", accessed)
	}

	return lipgloss.NewStyle().Foreground(styles.Subtle).Render(text)
}

// updateTOTP updates the TOTP code
func (s *EntryDetailScreen) updateTOTP() {
	config := totp.DefaultConfig(s.entry.TOTPSecret)
	code, expiresIn, err := config.GenerateCode()
	if err == nil {
		s.totpCode = code
		s.totpExpiresIn = expiresIn
	}
}

// getEntryIcon returns the icon for the entry type
func (s *EntryDetailScreen) getEntryIcon() string {
	switch s.entry.Type {
	case entity.EntryTypeLogin:
		return styles.IconGlobe
	case entity.EntryTypeSecureNote:
		return styles.IconNote
	case entity.EntryTypeCard:
		return styles.IconCard
	case entity.EntryTypeIdentity:
		return styles.IconIdentity
	default:
		return styles.IconKey
	}
}

// showCopyMessage shows a temporary copy message
func (s *EntryDetailScreen) showCopyMessage(msg string) {
	s.copyMessage = msg
	if s.copyTimer != nil {
		s.copyTimer.Stop()
	}
}

// clearCopyMessageCmd returns a command to clear the copy message
func (s *EntryDetailScreen) clearCopyMessageCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(t time.Time) tea.Msg {
		return clearCopyMsgMsg{}
	})
}

// tickCmd returns a command that waits for the next tick
func (s *EntryDetailScreen) tickCmd() tea.Cmd {
	return func() tea.Msg {
		<-s.ticker.C
		return tickMsg{}
	}
}

// BackMsg signals to go back to the vault list
type BackMsg struct{}

// EditEntryMsg signals to edit an entry
type EditEntryMsg struct {
	Entry *entity.Entry
}

// clearCopyMsgMsg signals to clear the copy message
type clearCopyMsgMsg struct{}
