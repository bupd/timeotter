package utils

import (
	"strings"
	"testing"
)

func TestExecuteShellCommand_Success(t *testing.T) {
	err := ExecuteShellCommand("echo hello")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestExecuteShellCommand_Failure(t *testing.T) {
	err := ExecuteShellCommand("exit 1")
	if err == nil {
		t.Error("expected error for failing command, got nil")
	}
}

func TestExecuteShellCommand_InvalidCommand(t *testing.T) {
	err := ExecuteShellCommand("nonexistent_command_12345")
	if err == nil {
		t.Error("expected error for invalid command, got nil")
	}
}

func TestExecuteShellCommand_MultipleCommands(t *testing.T) {
	err := ExecuteShellCommand("echo one && echo two")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestExecuteShellCommand_ErrorContainsOutput(t *testing.T) {
	err := ExecuteShellCommand("echo 'error message' && exit 1")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "error message") {
		t.Errorf("expected error to contain output, got %v", err)
	}
}

func TestExecuteShellCommand_EmptyCommand(t *testing.T) {
	// Empty command should succeed (bash returns 0 for empty input)
	err := ExecuteShellCommand("")
	if err != nil {
		t.Errorf("expected no error for empty command, got %v", err)
	}
}

func TestExecuteShellCommand_WithVariables(t *testing.T) {
	err := ExecuteShellCommand("VAR=test && echo $VAR")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
