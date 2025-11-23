package styles

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

// Color palette
var (
	// Primary colors
	Primary   = lipgloss.Color("#7C3AED") // Purple
	Secondary = lipgloss.Color("#06B6D4") // Cyan
	Success   = lipgloss.Color("#10B981") // Green
	Warning   = lipgloss.Color("#F59E0B") // Amber
	Danger    = lipgloss.Color("#EF4444") // Red
	Info      = lipgloss.Color("#3B82F6") // Blue

	// Grays
	Gray50  = lipgloss.Color("#F9FAFB")
	Gray100 = lipgloss.Color("#F3F4F6")
	Gray200 = lipgloss.Color("#E5E7EB")
	Gray300 = lipgloss.Color("#D1D5DB")
	Gray400 = lipgloss.Color("#9CA3AF")
	Gray500 = lipgloss.Color("#6B7280")
	Gray600 = lipgloss.Color("#4B5563")
	Gray700 = lipgloss.Color("#374151")
	Gray800 = lipgloss.Color("#1F2937")
	Gray900 = lipgloss.Color("#111827")

	// UI colors
	Background  = Gray900
	Foreground  = Gray100
	Border      = Gray700
	BorderFocus = Primary
	Subtle      = Gray500
)

// Common styles
var (
	// Title style
	TitleStyle = lipgloss.NewStyle().
			Foreground(Primary).
			Bold(true).
			Padding(0, 1)

	// Subtitle style
	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Subtle).
			Padding(0, 1)

	// Box style
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Border).
			Padding(1, 2)

	// Focused box style
	FocusedBoxStyle = BoxStyle.Copy().
			BorderForeground(BorderFocus)

	// List item style
	ListItemStyle = lipgloss.NewStyle().
			Padding(0, 2)

	// Selected list item style
	SelectedListItemStyle = ListItemStyle.Copy().
				Background(Primary).
				Foreground(Gray900).
				Bold(true)

	// Input style
	InputStyle = lipgloss.NewStyle().
			Foreground(Foreground).
			Background(Gray800).
			Padding(0, 1).
			Width(30)

	// Focused input style
	FocusedInputStyle = InputStyle.Copy().
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(Primary)

	// Button style
	ButtonStyle = lipgloss.NewStyle().
			Foreground(Gray900).
			Background(Primary).
			Padding(0, 2).
			Bold(true)

	// Secondary button style
	SecondaryButtonStyle = ButtonStyle.Copy().
				Background(Gray700).
				Foreground(Foreground)

	// Help text style
	HelpStyle = lipgloss.NewStyle().
			Foreground(Subtle).
			Padding(1, 0)

	// Error style
	ErrorStyle = lipgloss.NewStyle().
			Foreground(Danger).
			Bold(true).
			Padding(0, 1)

	// Success style
	SuccessStyle = lipgloss.NewStyle().
			Foreground(Success).
			Bold(true).
			Padding(0, 1)

	// Badge style
	BadgeStyle = lipgloss.NewStyle().
			Foreground(Gray900).
			Background(Secondary).
			Padding(0, 1).
			Bold(true)

	// TOTP badge style
	TOTPBadgeStyle = BadgeStyle.Copy().
			Background(Success)

	// Favorite star style
	FavoriteStyle = lipgloss.NewStyle().
			Foreground(Warning).
			Bold(true)

	// Folder icon style
	FolderStyle = lipgloss.NewStyle().
			Foreground(Secondary).
			Bold(true)

	// Login icon style
	LoginIconStyle = lipgloss.NewStyle().
			Foreground(Info)

	// Card icon style
	CardIconStyle = lipgloss.NewStyle().
			Foreground(Success)

	// Note icon style
	NoteIconStyle = lipgloss.NewStyle().
			Foreground(Warning)

	// Identity icon style
	IdentityIconStyle = lipgloss.NewStyle().
				Foreground(Primary)

	// Item style
	ItemStyle = lipgloss.NewStyle().
			PaddingLeft(2)

	// Selected Item style
	SelectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(0).
				Foreground(Primary)

	// Item Description style
	ItemDescriptionStyle = lipgloss.NewStyle().
				Foreground(Subtle).
				PaddingLeft(2)

	// Selected Item Description style
	SelectedItemDescriptionStyle = lipgloss.NewStyle().
					Foreground(Subtle).
					PaddingLeft(2)

	// Pagination style
	PaginationStyle = list.DefaultStyles().PaginationStyle.PaddingLeft(2)
)

// Icons
const (
	IconLock     = "🔐"
	IconUnlock   = "🔓"
	IconKey      = "🔑"
	IconGlobe    = "🌐"
	IconCard     = "💳"
	IconNote     = "📝"
	IconIdentity = "👤"
	IconFolder   = "📁"
	IconStar     = "⭐"
	IconWarning  = "⚠️"
	IconError    = "❌"
	IconSuccess  = "✓"
	IconInfo     = "ℹ️"
	IconSearch   = "🔍"
	IconClock    = "⏱"
	IconCopy     = "📋"
)

// RenderProgressBar renders a progress bar for TOTP countdown
func RenderProgressBar(percentage float64, width int) string {
	filled := int(float64(width) * percentage)
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	empty := max(width-filled, 0)

	// Choose color based on percentage
	var color lipgloss.Color
	switch {
	case percentage > 0.5:
		color = Success
	case percentage > 0.25:
		color = Warning
	default:
		color = Danger
	}

	filledStyle := lipgloss.NewStyle().Foreground(color)
	emptyStyle := lipgloss.NewStyle().Foreground(Gray700)

	bar := filledStyle.Render(strings.Repeat("█", filled)) +
		emptyStyle.Render(strings.Repeat("░", empty))

	return bar
}

// CenterHorizontal centers text horizontally
func CenterHorizontal(width int, text string) string {
	return lipgloss.Place(width, 1, lipgloss.Center, lipgloss.Center, text)
}
