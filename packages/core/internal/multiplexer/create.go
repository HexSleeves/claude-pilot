package multiplexer

import (
	"fmt"
	"slices"
	"sync"

	"claude-pilot/shared/interfaces"
)

// MultiplexerCache caches multiplexer instances to avoid repeated creation
var (
	cache      = make(map[string]interfaces.TerminalMultiplexer)
	cacheMutex sync.RWMutex
)

// CreateMultiplexer creates a multiplexer instance for the given backend
func CreateMultiplexer(backend, sessionPrefix string) (interfaces.TerminalMultiplexer, error) {
	if sessionPrefix == "" {
		sessionPrefix = "claude-pilot"
	}

	// Create cache key
	cacheKey := fmt.Sprintf("%s:%s", backend, sessionPrefix)

	// Check cache first
	cacheMutex.RLock()
	if cached, exists := cache[cacheKey]; exists {
		// Verify the cached instance is still available
		if cached.IsAvailable() {
			cacheMutex.RUnlock()
			return cached, nil
		}
		// Remove stale entry
		cacheMutex.RUnlock()
		cacheMutex.Lock()
		delete(cache, cacheKey)
		cacheMutex.Unlock()
	} else {
		cacheMutex.RUnlock()
	}

	// Create new instance
	var mux interfaces.TerminalMultiplexer
	var err error

	switch backend {
	case "tmux":
		mux, err = NewTmuxMultiplexer(sessionPrefix)
	case "zellij":
		mux, err = NewZellijMultiplexer(sessionPrefix)
	case "auto":
		mux, err = createAutoMultiplexer(sessionPrefix)
	default:
		return nil, fmt.Errorf("unsupported multiplexer backend: %s", backend)
	}

	if err != nil {
		return nil, err
	}

	// Cache the instance
	cacheMutex.Lock()
	cache[cacheKey] = mux
	cacheMutex.Unlock()

	return mux, nil
}

// GetAvailableBackends returns list of available multiplexer backends
func GetAvailableBackends(sessionPrefix string) []string {
	if sessionPrefix == "" {
		sessionPrefix = "claude-pilot"
	}

	var available []string

	// Check tmux availability
	if tmux, err := NewTmuxMultiplexer(sessionPrefix); err == nil && tmux.IsAvailable() {
		available = append(available, "tmux")
	}

	// Check zellij availability
	if zellij, err := NewZellijMultiplexer(sessionPrefix); err == nil && zellij.IsAvailable() {
		available = append(available, "zellij")
	}

	return available
}

// GetDefaultBackend returns the preferred backend (first available)
func GetDefaultBackend(sessionPrefix string) string {
	available := GetAvailableBackends(sessionPrefix)
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
func createAutoMultiplexer(sessionPrefix string) (interfaces.TerminalMultiplexer, error) {
	available := GetAvailableBackends(sessionPrefix)
	if len(available) == 0 {
		return nil, fmt.Errorf("no terminal multiplexer backends available (install tmux or zellij)")
	}

	// Prefer tmux if available
	if slices.Contains(available, "tmux") {
		return CreateMultiplexer("tmux", sessionPrefix)
	}

	// Otherwise use first available
	return CreateMultiplexer(available[0], sessionPrefix)
}
