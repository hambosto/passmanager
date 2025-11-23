package clipboard

import (
	"sync"
	"time"

	goclipboard "github.com/tiagomelo/go-clipboard/clipboard"
)

// Manager handles clipboard operations with auto-clear functionality
type Manager struct {
	clipboard goclipboard.Clipboard
	timeout   time.Duration
	timer     *time.Timer
	mu        sync.Mutex
}

// NewManager creates a new clipboard manager with the specified timeout
func NewManager(timeout time.Duration) *Manager {
	return &Manager{
		clipboard: goclipboard.New(),
		timeout:   timeout,
	}
}

// CopyWithTimeout copies text to clipboard and clears it after timeout
func (m *Manager) CopyWithTimeout(text string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Copy to clipboard using go-clipboard (Wayland compatible)
	if err := m.clipboard.CopyText(text); err != nil {
		return err
	}

	// Stop existing timer if any
	if m.timer != nil {
		m.timer.Stop()
	}

	// Set new timer to clear clipboard
	if m.timeout > 0 {
		m.timer = time.AfterFunc(m.timeout, func() {
			m.Clear()
		})
	}

	return nil
}

// Clear clears the clipboard
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear clipboard by copying empty string
	_ = m.clipboard.CopyText("")

	// Stop timer if running
	if m.timer != nil {
		m.timer.Stop()
		m.timer = nil
	}
}

// SetTimeout updates the auto-clear timeout
func (m *Manager) SetTimeout(timeout time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.timeout = timeout
}
