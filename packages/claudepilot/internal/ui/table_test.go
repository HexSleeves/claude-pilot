package ui

import (
	"claude-pilot/core/api"
	"claude-pilot/internal/styles"
	"strings"
	"testing"
	"time"
)

func TestSessionStatusToMultiplexerDisplay(t *testing.T) {
	tests := []struct {
		name           string
		status         api.SessionStatus
		expectedResult string
	}{
		{
			name:           "Connected status shows attached",
			status:         api.StatusConnected,
			expectedResult: FormatTmuxStatus("attached"),
		},
		{
			name:           "Active status shows running",
			status:         api.StatusActive,
			expectedResult: FormatTmuxStatus("running"),
		},
		{
			name:           "Inactive status shows stopped",
			status:         api.StatusInactive,
			expectedResult: FormatTmuxStatus("stopped"),
		},
		{
			name:           "Error status shows error",
			status:         api.StatusError,
			expectedResult: FormatTmuxStatus("error"),
		},
		{
			name:           "Unknown status shows unknown",
			status:         api.SessionStatus("unknown"),
			expectedResult: Dim("unknown"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sessionStatusToMultiplexerDisplay(tt.status)
			if result != tt.expectedResult {
				t.Errorf("sessionStatusToMultiplexerDisplay(%v) = %v, want %v", tt.status, result, tt.expectedResult)
			}
		})
	}
}

func TestSessionTable(t *testing.T) {
	// Create test sessions
	sessions := []*api.Session{
		{
			ID:          "test-id-1",
			Name:        "test-session-1",
			Status:      api.StatusActive,
			CreatedAt:   time.Now(),
			LastActive:  time.Now(),
			Messages:    []api.Message{},
			ProjectPath: "/test/path",
		},
		{
			ID:          "test-id-2",
			Name:        "test-session-2",
			Status:      api.StatusInactive,
			CreatedAt:   time.Now().Add(-time.Hour),
			LastActive:  time.Now().Add(-time.Hour),
			Messages:    []api.Message{},
			ProjectPath: "/another/path",
		},
	}

	// Test with backend string (should not panic)
	result := SessionTable(sessions, "tmux")
	if result == "" {
		t.Error("SessionTable returned empty string with valid sessions")
	}

	// Test with empty sessions
	emptyResult := SessionTable([]*api.Session{}, "tmux")
	if emptyResult != Dim("No active sessions found.") {
		t.Errorf("SessionTable with empty sessions should return 'No active sessions found.', got: %s", emptyResult)
	}
}

func TestSessionDetail(t *testing.T) {
	// Create test session
	session := &api.Session{
		ID:          "test-detail-id",
		Name:        "test-detail-session",
		Status:      api.StatusConnected,
		CreatedAt:   time.Now(),
		LastActive:  time.Now(),
		Description: "Test session description",
		Messages: []api.Message{
			{
				ID:        "msg-1",
				Role:      "user",
				Content:   "Test message content",
				Timestamp: time.Now(),
			},
		},
		ProjectPath: "/test/detail/path",
	}

	// Test with backend string (should not panic)
	result := SessionDetail(session, "tmux")
	if result == "" {
		t.Error("SessionDetail returned empty string with valid session")
	}

	// Check that the result contains expected session information
	if !strings.Contains(result, session.Name) {
		t.Error("SessionDetail result should contain session name")
	}
	if !strings.Contains(result, session.Description) {
		t.Error("SessionDetail result should contain session description")
	}
}

func TestFormatTime(t *testing.T) {
	testTime := time.Date(2025, 7, 13, 15, 30, 45, 0, time.UTC)
	result := formatTime(testTime)
	expected := Dim("2025-07-13 15:30")

	if result != expected {
		t.Errorf("formatTime(%v) = %v, want %v", testTime, result, expected)
	}
}

func TestFormatTimeAgo(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "Just now",
			time:     now.Add(-30 * time.Second),
			expected: Success.Sprint("just now"),
		},
		{
			name:     "1 minute ago",
			time:     now.Add(-1 * time.Minute),
			expected: Info.Sprint("1min ago"),
		},
		{
			name:     "Minutes ago",
			time:     now.Add(-5 * time.Minute),
			expected: Info.Sprint("5mins ago"),
		},
		{
			name:     "1 hour ago",
			time:     now.Add(-1 * time.Hour),
			expected: Warning.Sprint("1 hour ago"),
		},
		{
			name:     "Hours ago",
			time:     now.Add(-2 * time.Hour),
			expected: Warning.Sprint("2 hours ago"),
		},
		{
			name:     "Days ago",
			time:     now.Add(-25 * time.Hour),
			expected: Dim("1 day ago"),
		},
		{
			name:     "Multiple days ago",
			time:     now.Add(-72 * time.Hour),
			expected: Dim("3 days ago"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTimeAgo(tt.time)
			if result != tt.expected {
				t.Errorf("formatTimeAgo(%v) = %v, want %v", tt.time, result, tt.expected)
			}
		})
	}
}

func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		maxLen   int
		expected string
	}{
		{
			name:     "Short text unchanged",
			text:     "short",
			maxLen:   10,
			expected: "short",
		},
		{
			name:     "Long text truncated",
			text:     "this is a very long text that should be truncated",
			maxLen:   10,
			expected: "this is...",
		},
		{
			name:     "Exact length unchanged",
			text:     "exactly10c",
			maxLen:   10,
			expected: "exactly10c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := styles.TruncateText(tt.text, tt.maxLen)
			if result != tt.expected {
				t.Errorf("truncateText(%q, %d) = %q, want %q", tt.text, tt.maxLen, result, tt.expected)
			}
		})
	}
}
