package screens

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hambosto/passmanager/internal/presentation/tui/styles"
)

// HelpScreen shows keyboard shortcuts and help
type HelpScreen struct {
	width  int
	height int
}

// NewHelpScreen creates a new help screen
func NewHelpScreen() *HelpScreen {
	return &HelpScreen{}
}

// Init initializes the screen
func (s *HelpScreen) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (s *HelpScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		return s, nil

	case tea.KeyMsg:
		// Any key goes back
		return s, func() tea.Msg { return BackMsg{} }
	}

	return s, nil
}

// View renders the screen
func (s *HelpScreen) View() string {
	var b strings.Builder

	// Title
	title := styles.TitleStyle.Render("⌨️  Keyboard Shortcuts")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Sections
	sections := []struct {
		title string
		items [][2]string
	}{
		{
			title: "General",
			items: [][2]string{
				{"Ctrl+Q", "Quit application"},
				{"Ctrl+L", "Lock vault"},
				{"Esc", "Go back / Cancel"},
				{"?", "Show this help"},
				{"Ctrl+,", "Open settings"},
			},
		},
		{
			title: "Navigation",
			items: [][2]string{
				{"↑/k", "Move up"},
				{"↓/j", "Move down"},
				{"←/h", "Move left / Collapse"},
				{"→/l", "Move right / Expand"},
				{"Enter", "View / Open entry"},
				{"/", "Search / Filter"},
			},
		},
		{
			title: "Entry Management",
			items: [][2]string{
				{"Ctrl+N", "New entry"},
				{"Ctrl+E", "Edit entry"},
				{"Ctrl+D", "Delete entry"},
				{"Space", "Toggle favorite"},
				{"Ctrl+S", "Save (in editor)"},
			},
		},
		{
			title: "Clipboard Operations",
			items: [][2]string{
				{"Ctrl+U", "Copy username"},
				{"Ctrl+P", "Copy password"},
				{"Ctrl+T", "Copy TOTP code"},
				{"Ctrl+C", "Copy (in password generator)"},
			},
		},
		{
			title: "Password Generator",
			items: [][2]string{
				{"Ctrl+G", "Generate password (in editor)"},
				{"Ctrl+R", "Regenerate"},
				{"Tab", "Switch mode (password/passphrase)"},
				{"↑↓", "Navigate options"},
				{"←→", "Adjust values"},
			},
		},
		{
			title: "Other",
			items: [][2]string{
				{"Ctrl+H", "Show/Hide password"},
				{"Ctrl+O", "Open URL in browser"},
				{"Ctrl+F", "Toggle favorite (in editor)"},
			},
		},
	}

	for i, section := range sections {
		if i > 0 {
			b.WriteString("\n")
		}

		sectionBox := s.renderSection(section.title, section.items)
		b.WriteString(sectionBox)
		b.WriteString("\n")
	}

	// Footer
	footer := styles.HelpStyle.Render("Press any key to go back")
	b.WriteString("\n")
	b.WriteString(footer)

	return b.String()
}

// renderSection renders a help section
func (s *HelpScreen) renderSection(title string, items [][2]string) string {
	var content strings.Builder

	// Section title
	content.WriteString(lipgloss.NewStyle().Bold(true).Foreground(styles.Primary).Render(title))
	content.WriteString("\n\n")

	// Items
	keyStyle := lipgloss.NewStyle().Foreground(styles.Success).Bold(true)
	descStyle := lipgloss.NewStyle().Foreground(styles.Subtle)

	for _, item := range items {
		key := item[0]
		desc := item[1]

		content.WriteString(keyStyle.Render(key))
		content.WriteString(strings.Repeat(" ", max(0, 15-len(key))))
		content.WriteString(descStyle.Render(desc))
		content.WriteString("\n")
	}

	return content.String()
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
