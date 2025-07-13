package multiplexer

import (
	"fmt"
	"slices"

	"claude-pilot/internal/interfaces"
)

// Factory implements the MultiplexerFactory interface
type Factory struct {
	sessionPrefix string
}

// NewFactory creates a new multiplexer factory
func NewFactory(sessionPrefix string) *Factory {
	if sessionPrefix == "" {
		sessionPrefix = "claude-pilot"
	}
	return &Factory{
		sessionPrefix: sessionPrefix,
	}
}

// CreateMultiplexer creates a multiplexer instance for the given backend
func (f *Factory) CreateMultiplexer(backend string) (interfaces.TerminalMultiplexer, error) {
	switch backend {
	case "tmux":
		return NewTmuxMultiplexer(f.sessionPrefix)
	case "zellij":
		return NewZellijMultiplexer(f.sessionPrefix)
	case "auto":
		return f.createAutoMultiplexer()
	default:
		return nil, fmt.Errorf("unsupported multiplexer backend: %s", backend)
	}
}

// GetAvailableBackends returns list of available multiplexer backends
func (f *Factory) GetAvailableBackends() []string {
	var available []string

	// Check tmux availability
	if tmux, err := NewTmuxMultiplexer(f.sessionPrefix); err == nil && tmux.IsAvailable() {
		available = append(available, "tmux")
	}

	// Check zellij availability
	if zellij, err := NewZellijMultiplexer(f.sessionPrefix); err == nil && zellij.IsAvailable() {
		available = append(available, "zellij")
	}

	return available
}

// GetDefaultBackend returns the preferred backend (first available)
func (f *Factory) GetDefaultBackend() string {
	available := f.GetAvailableBackends()
	if len(available) == 0 {
		return "tmux" // fallback default
	}

	// Prefer tmux if available, otherwise use first available
	for _, backend := range available {
		if backend == "tmux" {
			return backend
		}
	}

	return available[0]
}

// createAutoMultiplexer automatically selects the best available backend
func (f *Factory) createAutoMultiplexer() (interfaces.TerminalMultiplexer, error) {
	available := f.GetAvailableBackends()
	if len(available) == 0 {
		return nil, fmt.Errorf("no terminal multiplexer backends available (install tmux or zellij)")
	}

	// Prefer tmux if available
	if slices.Contains(available, "tmux") {
		return f.CreateMultiplexer("tmux")
	}

	// Otherwise use first available
	return f.CreateMultiplexer(available[0])
}

// ValidateBackend checks if a backend name is valid
func (f *Factory) ValidateBackend(backend string) error {
	switch backend {
	case "tmux", "zellij", "auto":
		return nil
	default:
		return fmt.Errorf("invalid backend '%s', must be one of: tmux, zellij, auto", backend)
	}
}

// GetBackendInfo returns information about a specific backend
func (f *Factory) GetBackendInfo(backend string) (map[string]any, error) {
	if err := f.ValidateBackend(backend); err != nil {
		return nil, err
	}

	info := make(map[string]any)
	info["name"] = backend

	switch backend {
	case "tmux":
		if tmux, err := NewTmuxMultiplexer(f.sessionPrefix); err == nil {
			info["available"] = tmux.IsAvailable()
			info["description"] = "Terminal multiplexer with extensive features and scripting capabilities"
			info["features"] = []string{"mature", "scriptable", "extensive ecosystem"}
		}
	case "zellij":
		if zellij, err := NewZellijMultiplexer(f.sessionPrefix); err == nil {
			info["available"] = zellij.IsAvailable()
			info["description"] = "Modern terminal workspace with built-in session management"
			info["features"] = []string{"modern", "user-friendly", "plugin system"}
		}
	case "auto":
		info["description"] = "Automatically select the best available backend"
		info["default"] = f.GetDefaultBackend()
		info["available_backends"] = f.GetAvailableBackends()
	}

	return info, nil
}
