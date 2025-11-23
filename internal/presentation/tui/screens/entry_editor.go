package screens

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hambosto/passmanager/internal/domain/entity"
	"github.com/hambosto/passmanager/internal/presentation/tui/styles"
)

// EntryEditorScreen allows creating/editing entries
type EntryEditorScreen struct {
	entry  *entity.Entry
	isNew  bool
	width  int
	height int

	// Form inputs
	nameInput     textinput.Model
	usernameInput textinput.Model
	passwordInput textinput.Model
	uriInput      textinput.Model
	totpInput     textinput.Model
	notesArea     textarea.Model

	// State
	focusIndex   int
	showPassword bool
	isFavorite   bool
	entryType    entity.EntryType
}

// NewEntryEditorScreen creates a new entry editor screen
func NewEntryEditorScreen(entry *entity.Entry, isNew bool) *EntryEditorScreen {
	// Create inputs
	nameInput := textinput.New()
	nameInput.Placeholder = "Entry name"
	nameInput.Width = 40
	nameInput.Focus()

	usernameInput := textinput.New()
	usernameInput.Placeholder = "Username or email"
	usernameInput.Width = 40

	passwordInput := textinput.New()
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.EchoCharacter = '•'
	passwordInput.Width = 40

	uriInput := textinput.New()
	uriInput.Placeholder = "https://example.com"
	uriInput.Width = 40

	totpInput := textinput.New()
	totpInput.Placeholder = "otpauth://totp/... or secret"
	totpInput.Width = 40

	notesArea := textarea.New()
	notesArea.Placeholder = "Additional notes..."
	notesArea.SetWidth(60)
	notesArea.SetHeight(4)

	// Populate if editing existing entry
	if !isNew && entry != nil {
		nameInput.SetValue(entry.Name)
		usernameInput.SetValue(entry.Username)
		passwordInput.SetValue(entry.Password)
		uriInput.SetValue(entry.URI)
		totpInput.SetValue(entry.TOTPSecret)
		notesArea.SetValue(entry.Notes)
	}

	entryType := entity.EntryTypeLogin
	isFavorite := false
	if entry != nil {
		entryType = entry.Type
		isFavorite = entry.IsFavorite
	}

	return &EntryEditorScreen{
		entry:         entry,
		isNew:         isNew,
		nameInput:     nameInput,
		usernameInput: usernameInput,
		passwordInput: passwordInput,
		uriInput:      uriInput,
		totpInput:     totpInput,
		notesArea:     notesArea,
		focusIndex:    0,
		showPassword:  false,
		isFavorite:    isFavorite,
		entryType:     entryType,
	}
}

// Init initializes the screen
func (s *EntryEditorScreen) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles messages
func (s *EntryEditorScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			// Cancel editing
			return s, func() tea.Msg { return CancelEditMsg{} }

		case "ctrl+c", "ctrl+q":
			return s, tea.Quit

		case "ctrl+s":
			// Save entry
			return s, s.saveEntry()

		case "ctrl+g":
			// Open password generator
			return s, func() tea.Msg { return OpenPasswordGeneratorMsg{} }

		case "ctrl+h":
			// Toggle password visibility
			s.showPassword = !s.showPassword
			if s.showPassword {
				s.passwordInput.EchoMode = textinput.EchoNormal
			} else {
				s.passwordInput.EchoMode = textinput.EchoPassword
			}
			return s, nil

		case "tab", "shift+tab":
			// Navigate between inputs
			if msg.String() == "tab" {
				s.focusIndex++
			} else {
				s.focusIndex--
			}

			if s.focusIndex > 6 {
				s.focusIndex = 0
			} else if s.focusIndex < 0 {
				s.focusIndex = 6
			}

			s.updateFocus()
			return s, textinput.Blink

		case "ctrl+f":
			// Toggle favorite
			s.isFavorite = !s.isFavorite
			return s, nil
		}
	}

	// Update the focused input
	switch s.focusIndex {
	case 0:
		s.nameInput, cmd = s.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	case 1:
		s.usernameInput, cmd = s.usernameInput.Update(msg)
		cmds = append(cmds, cmd)
	case 2:
		s.passwordInput, cmd = s.passwordInput.Update(msg)
		cmds = append(cmds, cmd)
	case 3:
		s.uriInput, cmd = s.uriInput.Update(msg)
		cmds = append(cmds, cmd)
	case 4:
		s.totpInput, cmd = s.totpInput.Update(msg)
		cmds = append(cmds, cmd)
	case 5:
		s.notesArea, cmd = s.notesArea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return s, tea.Batch(cmds...)
}

// View renders the screen
func (s *EntryEditorScreen) View() string {
	var b strings.Builder

	// Title
	title := "New Entry"
	if !s.isNew {
		title = "Edit Entry: " + s.entry.Name
	}
	b.WriteString(styles.TitleStyle.Render(styles.IconKey + " " + title))
	b.WriteString("\n\n")

	// Entry type
	typeStr := lipgloss.NewStyle().Foreground(styles.Subtle).Render("Type: Login")
	b.WriteString(typeStr)
	b.WriteString("  ")

	// Favorite checkbox
	favIcon := "☐"
	if s.isFavorite {
		favIcon = "☑"
	}
	favStr := styles.FavoriteStyle.Render(favIcon + " Favorite")
	b.WriteString(favStr)
	b.WriteString("\n\n")

	// Form
	var formContent strings.Builder

	// Name
	formContent.WriteString(s.renderField("Name:", s.nameInput.View(), s.focusIndex == 0))
	formContent.WriteString("\n\n")

	// Username
	formContent.WriteString(s.renderField("Username:", s.usernameInput.View(), s.focusIndex == 1))
	formContent.WriteString("\n\n")

	// Password
	passView := s.passwordInput.View()
	if s.focusIndex == 2 {
		passView += "\n" + styles.HelpStyle.Render("[Ctrl+G] Generate  [Ctrl+H] Show/Hide")
	}
	formContent.WriteString(s.renderField("Password:", passView, s.focusIndex == 2))
	formContent.WriteString("\n\n")

	// Website
	formContent.WriteString(s.renderField("Website:", s.uriInput.View(), s.focusIndex == 3))
	formContent.WriteString("\n\n")

	// TOTP
	totpView := s.totpInput.View()
	if s.focusIndex == 4 {
		totpView += "\n" + styles.HelpStyle.Render("Enter otpauth:// URI or Base32 secret")
	}
	formContent.WriteString(s.renderField("TOTP:", totpView, s.focusIndex == 4))
	formContent.WriteString("\n\n")

	// Notes
	formContent.WriteString(s.renderField("Notes:", s.notesArea.View(), s.focusIndex == 5))

	// Render form box
	box := styles.BoxStyle.
		Width(min(70, s.width-4)).
		Render(formContent.String())
	b.WriteString(box)
	b.WriteString("\n\n")

	// Help text
	helpText := "[Ctrl+S] Save  •  [Esc] Cancel  •  [Tab] Next Field  •  [Ctrl+F] Toggle Favorite"
	b.WriteString(styles.HelpStyle.Render(helpText))

	return b.String()
}

// renderField renders a form field
func (s *EntryEditorScreen) renderField(label, value string, focused bool) string {
	labelStyle := lipgloss.NewStyle().Bold(true)
	if focused {
		labelStyle = labelStyle.Foreground(styles.Primary)
	}

	return labelStyle.Render(label) + "\n" + value
}

// updateFocus updates which input is focused
func (s *EntryEditorScreen) updateFocus() {
	s.nameInput.Blur()
	s.usernameInput.Blur()
	s.passwordInput.Blur()
	s.uriInput.Blur()
	s.totpInput.Blur()
	s.notesArea.Blur()

	switch s.focusIndex {
	case 0:
		s.nameInput.Focus()
	case 1:
		s.usernameInput.Focus()
	case 2:
		s.passwordInput.Focus()
	case 3:
		s.uriInput.Focus()
	case 4:
		s.totpInput.Focus()
	case 5:
		s.notesArea.Focus()
	}
}

// saveEntry creates a command to save the entry
func (s *EntryEditorScreen) saveEntry() tea.Cmd {
	// Validate
	if s.nameInput.Value() == "" {
		// Show error - name is required
		return func() tea.Msg {
			return ShowErrorMsg{Error: "Name is required"}
		}
	}

	// Create or update entry
	if s.isNew {
		entry := entity.NewEntry(s.entryType, s.nameInput.Value())
		entry.Username = s.usernameInput.Value()
		entry.Password = s.passwordInput.Value()
		entry.URI = s.uriInput.Value()
		entry.TOTPSecret = s.totpInput.Value()
		entry.Notes = s.notesArea.Value()
		entry.IsFavorite = s.isFavorite

		return func() tea.Msg {
			return SaveEntryMsg{Entry: entry, IsNew: true}
		}
	} else {
		// Update existing entry
		s.entry.Name = s.nameInput.Value()
		s.entry.Username = s.usernameInput.Value()
		s.entry.Password = s.passwordInput.Value()
		s.entry.URI = s.uriInput.Value()
		s.entry.TOTPSecret = s.totpInput.Value()
		s.entry.Notes = s.notesArea.Value()
		s.entry.IsFavorite = s.isFavorite
		s.entry.Update()

		return func() tea.Msg {
			return SaveEntryMsg{Entry: s.entry, IsNew: false}
		}
	}
}

// SetPassword sets the password field value
func (s *EntryEditorScreen) SetPassword(password string) {
	s.passwordInput.SetValue(password)
}

// CancelEditMsg signals that editing was cancelled
type CancelEditMsg struct{}

// SaveEntryMsg signals that an entry should be saved
type SaveEntryMsg struct {
	Entry *entity.Entry
	IsNew bool
}

// OpenPasswordGeneratorMsg signals to open the password generator
type OpenPasswordGeneratorMsg struct{}

// ShowErrorMsg signals an error to display
type ShowErrorMsg struct {
	Error string
}
