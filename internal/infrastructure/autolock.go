package infrastructure

import (
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// AutoLocker handles automatic vault locking after inactivity
type AutoLocker struct {
	timeout      time.Duration
	lastActivity time.Time
	ticker       *time.Ticker
	lockCallback func() tea.Msg
	mu           sync.Mutex
	running      bool
}

// NewAutoLocker creates a new auto-locker
func NewAutoLocker(timeout time.Duration, lockCallback func() tea.Msg) *AutoLocker {
	return &AutoLocker{
		timeout:      timeout,
		lastActivity: time.Now(),
		lockCallback: lockCallback,
		running:      false,
	}
}

// Start starts the auto-lock timer
func (a *AutoLocker) Start() tea.Cmd {
	if a.timeout == 0 {
		// Auto-lock disabled
		return nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if a.running {
		return nil
	}

	a.running = true
	a.lastActivity = time.Now()
	a.ticker = time.NewTicker(1 * time.Second)

	return func() tea.Msg {
		for range a.ticker.C {
			a.mu.Lock()
			if !a.running {
				a.mu.Unlock()
				return nil
			}

			elapsed := time.Since(a.lastActivity)
			if elapsed >= a.timeout {
				a.mu.Unlock()
				return a.lockCallback()
			}
			a.mu.Unlock()
		}
		return nil
	}
}

// Stop stops the auto-lock timer
func (a *AutoLocker) Stop() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if !a.running {
		return
	}

	a.running = false
	if a.ticker != nil {
		a.ticker.Stop()
	}
}

// Reset Resets the inactivity timer
func (a *AutoLocker) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lastActivity = time.Now()
}

// TimeUntilLock returns the time until auto-lock
func (a *AutoLocker) TimeUntilLock() time.Duration {
	if a.timeout == 0 {
		return 0
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	elapsed := time.Since(a.lastActivity)
	remaining := a.timeout - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}

// SetTimeout sets the auto-lock timeout
func (a *AutoLocker) SetTimeout(timeout time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.timeout = timeout
}

// IsEnabled returns whether auto-lock is enabled
func (a *AutoLocker) IsEnabled() bool {
	return a.timeout > 0
}

// AutoLockMsg signals that the vault should be locked
type AutoLockMsg struct{}
