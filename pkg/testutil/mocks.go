package testutil

import (
	"errors"
	"sync"
)

// MockShellExecutor is a mock for shell command execution.
type MockShellExecutor struct {
	mu             sync.Mutex
	ExecutedCmds   []string
	ShouldFail     bool
	FailureMessage string
}

// NewMockShellExecutor creates a new MockShellExecutor.
func NewMockShellExecutor() *MockShellExecutor {
	return &MockShellExecutor{
		ExecutedCmds: make([]string, 0),
	}
}

// ExecuteShellCommand records the command and optionally returns an error.
func (m *MockShellExecutor) ExecuteShellCommand(command string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ExecutedCmds = append(m.ExecutedCmds, command)

	if m.ShouldFail {
		msg := m.FailureMessage
		if msg == "" {
			msg = "mock shell execution failed"
		}
		return errors.New(msg)
	}
	return nil
}

// GetExecutedCommands returns a copy of executed commands.
func (m *MockShellExecutor) GetExecutedCommands() []string {
	m.mu.Lock()
	defer m.mu.Unlock()

	cmds := make([]string, len(m.ExecutedCmds))
	copy(cmds, m.ExecutedCmds)
	return cmds
}

// Reset clears the executed commands and resets state.
func (m *MockShellExecutor) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ExecutedCmds = make([]string, 0)
	m.ShouldFail = false
	m.FailureMessage = ""
}

// MockCronManager is a mock for cron operations.
type MockCronManager struct {
	mu           sync.Mutex
	AddedCrons   []CronEntry
	ClearedCount int
	ShouldFail   bool
}

// CronEntry represents a cron job entry.
type CronEntry struct {
	CronJob   string
	CmdToExec string
}

// NewMockCronManager creates a new MockCronManager.
func NewMockCronManager() *MockCronManager {
	return &MockCronManager{
		AddedCrons: make([]CronEntry, 0),
	}
}

// AddCrons records a cron addition.
func (m *MockCronManager) AddCrons(cronJob, cmdToExec string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ShouldFail {
		return errors.New("mock cron add failed")
	}

	m.AddedCrons = append(m.AddedCrons, CronEntry{
		CronJob:   cronJob,
		CmdToExec: cmdToExec,
	})
	return nil
}

// ClearCronJobs records a cron clear operation.
func (m *MockCronManager) ClearCronJobs(backupFile, cronMarker string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.ShouldFail {
		return errors.New("mock cron clear failed")
	}

	m.ClearedCount++
	return nil
}

// GetAddedCrons returns a copy of added cron entries.
func (m *MockCronManager) GetAddedCrons() []CronEntry {
	m.mu.Lock()
	defer m.mu.Unlock()

	entries := make([]CronEntry, len(m.AddedCrons))
	copy(entries, m.AddedCrons)
	return entries
}

// Reset clears the state.
func (m *MockCronManager) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.AddedCrons = make([]CronEntry, 0)
	m.ClearedCount = 0
	m.ShouldFail = false
}
