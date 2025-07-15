package utils

import (
	"claude-pilot/shared/interfaces"
	"claude-pilot/shared/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ParseSessionStatus is the exported version of parseSessionStatus
func ParseSessionStatus(status string) string {
	statusLower := strings.ToLower(strings.TrimSpace(status))

	switch statusLower {
	case "error", "failed", "failure":
		return "error"
	case "warning", "warn", "caution":
		return "warning"
	case "active", "running", "online":
		return "active"
	case "connected", "attached", "linked":
		return "connected"
	case "inactive", "stopped", "offline":
		return "inactive"
	default:
		return "unknown"
	}
}

// GetSessionStatusColor returns the appropriate color for a given status string
func GetSessionStatusColor(status string) lipgloss.Color {
	statusLower := strings.ToLower(strings.TrimSpace(status))
	var borderColor lipgloss.Color = styles.InfoColor

	switch statusLower {
	case string(interfaces.StatusError): // error
		borderColor = styles.ErrorColor
	case string(interfaces.StatusWarning): // warning
		borderColor = styles.WarningColor
	case string(interfaces.StatusActive): // active
		borderColor = styles.SuccessColor
	case string(interfaces.StatusInactive): // inactive
		borderColor = styles.WarningColor
	case string(interfaces.StatusConnected): // connected
		borderColor = styles.InfoColor
	default: // unknown
		// Check if status contains keywords as fallback
		if strings.Contains(statusLower, string(interfaces.StatusError)) {
			borderColor = styles.ErrorColor
		} else if strings.Contains(statusLower, string(interfaces.StatusWarning)) {
			borderColor = styles.WarningColor
		} else if strings.Contains(statusLower, string(interfaces.StatusActive)) || strings.Contains(statusLower, string(interfaces.StatusConnected)) {
			borderColor = styles.SuccessColor
		}
	}

	return borderColor
}

// Utility functions for status formatting
func FormatSessionStatus(status string) lipgloss.Style {
	parsedStatus := ParseSessionStatus(status)
	switch parsedStatus {
	case "active":
		return styles.SessionStatusActiveStyle
	case "inactive":
		return styles.SessionStatusInactiveStyle
	case "connected":
		return styles.SessionStatusConnectedStyle
	case "error":
		return styles.SessionStatusErrorStyle
	default:
		return styles.MutedTextStyle
	}
}
