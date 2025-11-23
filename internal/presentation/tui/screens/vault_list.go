package screens

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hambosto/passmanager/internal/domain/entity"
	"github.com/hambosto/passmanager/internal/infrastructure/clipboard"
	"github.com/hambosto/passmanager/internal/presentation/tui/styles"
)

// VaultListScreen shows the list of vault entries
type VaultListScreen struct {
	vault     *entity.Vault
	list      list.Model
	clipboard *clipboard.Manager
	width     int
	height    int

	// TOTP ticker for real-time updates
	ticker *time.Ticker
}

// entryItem wraps an entry for the list
type entryItem struct {
	entry *entity.Entry
}

func (i entryItem) Title() string       { return i.entry.Name }
func (i entryItem) FilterValue() string { return i.entry.Name }

// itemDelegate is a custom delegate for rendering list items
type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 3 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(entryItem)
	if !ok {
		return
	}

	fn := styles.ItemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return styles.SelectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	// Icon based on type
	icon := ""
	switch i.entry.Type {
	case entity.EntryTypeLogin:
		icon = styles.IconGlobe
	case entity.EntryTypeSecureNote:
		icon = styles.IconNote
	case entity.EntryTypeCard:
		icon = styles.IconCard
	case entity.EntryTypeIdentity:
		icon = styles.IconIdentity
	}

	// Title line
	var titleBuilder strings.Builder
	titleBuilder.WriteString(icon)
	titleBuilder.WriteString(" ")
	titleBuilder.WriteString(i.entry.Name)
	if i.entry.IsFavorite {
		titleBuilder.WriteString(" ")
		titleBuilder.WriteString(styles.IconStar)
	}

	fmt.Fprint(w, fn(titleBuilder.String()))
	fmt.Fprint(w, "\n")

	// Description line (URI / Username)
	var descBuilder strings.Builder
	if i.entry.Username != "" {
		descBuilder.WriteString(i.entry.Username)
	}
	if i.entry.URI != "" {
		if descBuilder.Len() > 0 {
			descBuilder.WriteString(" • ")
		}
		descBuilder.WriteString(i.entry.URI)
	}

	descStyle := styles.ItemDescriptionStyle
	if index == m.Index() {
		descStyle = styles.SelectedItemDescriptionStyle
	}

	if descBuilder.Len() > 0 {
		fmt.Fprint(w, descStyle.Render("  "+descBuilder.String()))
		fmt.Fprint(w, "\n")
	} else {
		fmt.Fprint(w, "\n")
	}

	// Meta line (TOTP / Last Used)
	var metaBuilder strings.Builder
	if i.entry.TOTPSecret != "" {
		metaBuilder.WriteString("TOTP")
	}
	if !i.entry.AccessedAt.IsZero() {
		if metaBuilder.Len() > 0 {
			metaBuilder.WriteString(" • ")
		}
		metaBuilder.WriteString(formatDuration(time.Since(i.entry.AccessedAt)))
	}

	if metaBuilder.Len() > 0 {
		fmt.Fprint(w, descStyle.Render("  "+metaBuilder.String()))
	}
}

// NewVaultListScreen creates a new vault list screen
func NewVaultListScreen(vault *entity.Vault, clipboardMgr *clipboard.Manager) *VaultListScreen {
	// Create list items from vault entries
	items := make([]list.Item, len(vault.Entries))
	for i, entry := range vault.Entries {
		items[i] = entryItem{entry: entry}
	}

	// Create list with custom delegate
	l := list.New(items, itemDelegate{}, 0, 0)
	l.Title = "Vault Entries"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.TitleStyle
	l.Styles.PaginationStyle = styles.PaginationStyle
	l.Styles.HelpStyle = styles.HelpStyle

	return &VaultListScreen{
		vault:     vault,
		list:      l,
		clipboard: clipboardMgr,
		ticker:    time.NewTicker(1 * time.Second),
	}
}

// Init initializes the screen
func (s *VaultListScreen) Init() tea.Cmd {
	return tea.Batch(
		s.tickCmd(),
	)
}

// Update handles messages
func (s *VaultListScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.width = msg.Width
		s.height = msg.Height
		s.list.SetWidth(msg.Width)
		s.list.SetHeight(msg.Height - 4) // Reserve space for help
		return s, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "ctrl+q":
			return s, tea.Quit

		case "ctrl+n":
			// New entry
			return s, func() tea.Msg {
				return NewEntryMsg{}
			}
		}

	case tickMsg:
		// Update TOTP codes
		return s, s.tickCmd()
	}

	// Update list
	s.list, cmd = s.list.Update(msg)
	return s, cmd
}

// GetSelectedEntry returns the currently selected entry
func (s *VaultListScreen) GetSelectedEntry() *entity.Entry {
	selectedItem := s.list.SelectedItem()
	if selectedItem == nil {
		return nil
	}

	if item, ok := selectedItem.(entryItem); ok {
		return item.entry
	}

	return nil
}

// View renders the screen
func (s *VaultListScreen) View() string {
	return s.list.View()
}

// tickCmd returns a command that waits for the next tick
func (s *VaultListScreen) tickCmd() tea.Cmd {
	return func() tea.Msg {
		<-s.ticker.C
		return tickMsg{}
	}
}

// tickMsg is sent when it's time to update TOTP codes
type tickMsg struct{}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
	return fmt.Sprintf("%dd ago", int(d.Hours()/24))
}

// NewEntryMsg signals to create a new entry
type NewEntryMsg struct{}
