package config

import (
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	config := DefaultConfig()

	// Test UI mode default
	if config.UI.Mode != "cli" {
		t.Errorf("Default UI mode should be 'cli', got '%s'", config.UI.Mode)
	}

	// Test backend default
	if config.Backend != "auto" {
		t.Errorf("Default backend should be 'auto', got '%s'", config.Backend)
	}

	// Test sessions directory default (should be set)
	if config.SessionsDir == "" {
		t.Error("Default sessions directory should not be empty")
	}

	// Test default shell
	if config.DefaultShell != "claude" {
		t.Errorf("Default shell should be 'claude', got '%s'", config.DefaultShell)
	}
}

func TestConfigManagerCreation(t *testing.T) {
	// Test with empty config path (should use default)
	cm := NewConfigManager("")
	if cm == nil {
		t.Error("NewConfigManager should not return nil")
	}

	// Test with custom config path
	customPath := "/tmp/test-config.yaml"
	cm = NewConfigManager(customPath)
	if cm == nil {
		t.Error("NewConfigManager should not return nil with custom path")
	}

	// Test that we can get config from the manager
	config := cm.GetConfig()
	if config == nil {
		t.Error("GetConfig should not return nil")
	}
}

func TestConfigManagerOperations(t *testing.T) {
	cm := NewConfigManager("")
	
	// Test getting initial config
	config := cm.GetConfig()
	if config == nil {
		t.Error("GetConfig should not return nil")
	}

	// Test updating config
	newConfig := DefaultConfig()
	newConfig.UI.Mode = "tui"
	cm.UpdateConfig(newConfig)

	// Verify update
	updatedConfig := cm.GetConfig()
	if updatedConfig.UI.Mode != "tui" {
		t.Errorf("Expected UI mode to be 'tui', got '%s'", updatedConfig.UI.Mode)
	}
}

func TestUIConfig(t *testing.T) {
	config := DefaultConfig()
	
	// Test that UI config exists
	if config.UI.Mode == "" {
		t.Error("UI mode should not be empty")
	}

	// Test valid modes
	validModes := []string{"cli", "tui"}
	found := false
	for _, mode := range validModes {
		if config.UI.Mode == mode {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Default UI mode '%s' should be one of: %v", config.UI.Mode, validModes)
	}
}

func TestBackendConfig(t *testing.T) {
	config := DefaultConfig()
	
	// Test that backend is set
	if config.Backend == "" {
		t.Error("Backend should not be empty")
	}

	// Test valid backends
	validBackends := []string{"auto", "tmux", "zellij"}
	found := false
	for _, backend := range validBackends {
		if config.Backend == backend {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Default backend '%s' should be one of: %v", config.Backend, validBackends)
	}
}

func TestHomeDirectoryExpansion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Tilde only",
			input:    "~",
			expected: "/expected/home",
		},
		{
			name:     "Tilde with path",
			input:    "~/test/path",
			expected: "/expected/home/test/path",
		},
		{
			name:     "No tilde",
			input:    "/absolute/path",
			expected: "/absolute/path",
		},
		{
			name:     "Relative path",
			input:    "relative/path",
			expected: "relative/path",
		},
		{
			name:     "Empty path",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandHomePath(tt.input, "/expected/home")
			if result != tt.expected {
				t.Errorf("expandHomePath(%q, %q) = %q, want %q", 
					tt.input, "/expected/home", result, tt.expected)
			}
		})
	}
}