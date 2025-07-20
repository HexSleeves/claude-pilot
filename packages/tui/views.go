package tui

import (
	"claude-pilot/shared/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

// renderTableView renders the main dashboard with session table, header, and footer.
// This is the primary view showing the list of sessions in a table format with
// navigation controls and status information.
func renderTableView(m Model) string {
	var b strings.Builder

	// Header
	header := renderHeader(m)
	b.WriteString(header)
	b.WriteString("\n\n")

	// Status message if present
	if m.statusMessage != "" {
		statusLine := styles.InfoStyle.Render(m.statusMessage)
		b.WriteString(statusLine)
		b.WriteString("\n\n")
	}

	// Selection information when multi-select is enabled
	if len(m.tableSelectedRows) > 0 {
		selectionInfo := fmt.Sprintf("%d row(s) selected", len(m.tableSelectedRows))
		b.WriteString(styles.InfoStyle.Render(selectionInfo))
		b.WriteString("\n")
	}

	// Table
	if len(m.sessions) == 0 {
		emptyMessage := styles.MutedTextStyle.Render("No sessions found. Press 'c' to create a new session.")
		b.WriteString(emptyMessage)
	} else {
		b.WriteString(m.table.View())
	}

	b.WriteString("\n\n")

	// Footer with key shortcuts
	footer := renderFooter(m)
	b.WriteString(footer)

	return b.String()
}

// renderCreateView renders the session creation form with input fields.
// This view provides a form interface for creating new sessions with
// name, description, and project path inputs.
func renderCreateView(m Model) string {
	var b strings.Builder

	// Header
	title := styles.TitleStyle.Render("Create New Session")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Form fields
	b.WriteString(styles.BoldStyle.Render("Session Name:"))
	b.WriteString("\n")
	if m.activeInput == 0 {
		b.WriteString(styles.InputFocusedStyle.Render(m.nameInput.View()))
	} else {
		b.WriteString(styles.InputStyle.Render(m.nameInput.View()))
	}
	b.WriteString("\n\n")

	b.WriteString(styles.BoldStyle.Render("Description (optional):"))
	b.WriteString("\n")
	if m.activeInput == 1 {
		b.WriteString(styles.InputFocusedStyle.Render(m.descriptionInput.View()))
	} else {
		b.WriteString(styles.InputStyle.Render(m.descriptionInput.View()))
	}
	b.WriteString("\n\n")

	b.WriteString(styles.BoldStyle.Render("Project Path (optional):"))
	b.WriteString("\n")
	if m.activeInput == 2 {
		b.WriteString(styles.InputFocusedStyle.Render(m.pathInput.View()))
	} else {
		b.WriteString(styles.InputStyle.Render(m.pathInput.View()))
	}
	b.WriteString("\n\n")

	// Instructions
	instructions := []string{
		"Tab/Shift+Tab: Navigate fields",
		"Enter: Create session",
		"Esc: Cancel",
	}

	for _, instruction := range instructions {
		b.WriteString(styles.MutedTextStyle.Render(instruction))
		b.WriteString("\n")
	}

	return b.String()
}

// renderLoadingView renders a loading spinner during API operations.
// This view is displayed while waiting for asynchronous operations
// such as loading sessions or creating new sessions.
func renderLoadingView(m Model) string {
	var b strings.Builder

	// Create a simple spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle

	b.WriteString("\n\n")
	b.WriteString(styles.SpinnerStyle.Render("●"))
	b.WriteString(" ")
	b.WriteString(styles.PrimaryTextStyle.Render("Loading..."))
	b.WriteString("\n\n")

	return b.String()
}

// renderErrorView renders error display with retry options.
// This view shows error messages with user-friendly retry and quit options
// when API operations or other operations fail.
func renderErrorView(m Model) string {
	var b strings.Builder

	// Error title
	title := styles.ErrorStyle.Render("Error")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Error message
	errorBox := styles.ErrorBoxStyle.Render(m.errorMessage)
	b.WriteString(errorBox)
	b.WriteString("\n\n")

	// Instructions
	instructions := []string{
		"r: Retry",
		"q: Quit",
	}

	for _, instruction := range instructions {
		b.WriteString(styles.MutedTextStyle.Render(instruction))
		b.WriteString("\n")
	}

	return b.String()
}

// renderHelpView renders help modal with keyboard shortcuts.
// This view displays a comprehensive list of available keyboard shortcuts
// and navigation options for the TUI interface.
func renderHelpView(m Model) string {
	var b strings.Builder

	// Help title
	title := styles.TitleStyle.Render("Help - Keyboard Shortcuts")
	b.WriteString(title)

	// Create help content using the keymap
	h := help.New()
	helpContent := h.Styles.ShortSeparator.Inline(true).Render(h.ShortSeparator) + h.View(m.keymap)

	// Wrap in a styled box
	helpBox := styles.InfoBoxStyle.Render(helpContent)
	b.WriteString(helpBox)
	b.WriteString("\n\n")

	// Instructions
	instruction := styles.MutedTextStyle.Render("Press '?' again or Esc to close help")
	b.WriteString(instruction)

	return b.String()
}

// renderHeader renders the application header with title and backend info.
// The header includes the application title, backend connection information,
// and last refresh timestamp for user context.
func renderHeader(m Model) string {
	var b strings.Builder

	// App title
	title := styles.TitleStyle.Render("Claude Pilot - Session Manager")
	b.WriteString(title)

	// Backend info
	backend := fmt.Sprintf("Backend: %s", m.client.GetBackend())
	backendInfo := styles.SecondaryTextStyle.Render(backend)

	// Last refresh time
	var refreshInfo string
	if !m.lastRefresh.IsZero() {
		refreshInfo = fmt.Sprintf("Last refresh: %s", m.lastRefresh.Format("15:04:05"))
	} else {
		refreshInfo = "Loading..."
	}
	refreshText := styles.MutedTextStyle.Render(refreshInfo)

	// Layout header info
	headerInfo := lipgloss.JoinHorizontal(
		lipgloss.Top,
		backendInfo,
		strings.Repeat(" ", 10),
		refreshText,
	)

	b.WriteString("\n")
	b.WriteString(headerInfo)

	return b.String()
}

// renderFooter renders the footer with key shortcuts.
// The footer displays context-sensitive keyboard shortcuts based on
// the current view state for quick user reference.
func renderFooter(m Model) string {
	var shortcuts []string

	switch m.currentView {
	case TableView:
		if m.showTableHelp {
			shortcuts = []string{
				"s+n: Sort by name",
				"s+s: Sort by status",
				"s+c: Sort by created",
				"s+a: Sort by last active",
				"/: Filter",
				"space: Select row",
				"ctrl+a: Select all",
				"ctrl+n/p: Next/Prev page",
				"?: Toggle help",
			}
		} else {
			shortcuts = []string{
				"↑/↓: Navigate",
				"Enter: Attach",
				"c: Create",
				"k: Kill",
				"/: Filter",
				"space: Select",
				"r: Refresh",
				"?: Help",
				"q: Quit",
			}
		}
	case CreatePrompt:
		shortcuts = []string{
			"Tab: Next field",
			"Enter: Create",
			"Esc: Cancel",
		}
	case KillConfirmation:
		shortcuts = []string{
			"y: Yes, kill session",
			"n: No, cancel",
			"Esc: Cancel",
		}
	case Error:
		shortcuts = []string{
			"r: Retry",
			"q: Quit",
		}
	default:
		shortcuts = []string{
			"?: Help",
			"q: Quit",
		}
	}

	// Join shortcuts with separators
	shortcutText := strings.Join(shortcuts, " • ")
	return styles.FooterStyle.Render(shortcutText)
}

// renderKillConfirmationView renders the kill confirmation dialog
func renderKillConfirmationView(m Model) string {
	var b strings.Builder

	// Header
	header := renderHeader(m)
	b.WriteString(header)
	b.WriteString("\n\n")

	if m.sessionToKill == nil {
		// Fallback if no session selected
		b.WriteString(styles.ErrorStyle.Render("No session selected"))
		b.WriteString("\n\n")
		b.WriteString(renderFooter(m))
		return b.String()
	}

	// Confirmation dialog
	sessionInfo := fmt.Sprintf("Session: %s (ID: %s)",
		styles.BoldStyle.Render(m.sessionToKill.Name),
		styles.MutedTextStyle.Render(m.sessionToKill.ID))

	warningText := styles.ErrorStyle.Render("⚠️  Are you sure you want to kill this session?")
	confirmationText := "This action cannot be undone."

	optionsText := fmt.Sprintf("%s or %s",
		styles.KeyStyle.Render("y")+" Yes",
		styles.KeyStyle.Render("n")+" No")

	// Center the dialog content
	dialogContent := lipgloss.JoinVertical(lipgloss.Center,
		sessionInfo,
		"",
		warningText,
		confirmationText,
		"",
		optionsText,
	)

	// Create a bordered box around the dialog
	dialog := styles.DialogBoxStyle.Render(dialogContent)

	// Center the dialog on screen
	centered := lipgloss.Place(
		m.totalWidth,
		m.totalHeight-4, // Account for header/footer
		lipgloss.Center,
		lipgloss.Center,
		dialog,
	)

	b.WriteString(centered)
	b.WriteString("\n")
	b.WriteString(renderFooter(m))

	return b.String()
}

// renderFilterView renders the filter input dialog
func renderFilterView(m Model) string {
	var b strings.Builder

	// Header
	header := renderHeader(m)
	b.WriteString(header)
	b.WriteString("\n\n")

	// Filter dialog content
	title := styles.TitleStyle.Render("Filter Sessions")

	// Show number of matches if filter is active
	var matchInfo string
	if m.filterActive && m.filteredSessions != nil {
		matchCount := len(m.filteredSessions)
		totalCount := len(m.sessions)
		matchInfo = fmt.Sprintf("Showing %d of %d sessions", matchCount, totalCount)
	} else {
		matchInfo = fmt.Sprintf("Search across %d sessions", len(m.sessions))
	}

	matchText := styles.MutedTextStyle.Render(matchInfo)

	// Filter input
	inputLabel := styles.BoldStyle.Render("Filter:")
	filterInput := styles.InputFocusedStyle.Render(m.filterInput.View())

	// Help text
	helpText := "Search in: name, description, status, project path, session ID"
	helpLine := styles.MutedTextStyle.Render(helpText)

	// Instructions
	instructionsText := fmt.Sprintf("%s to apply • %s to cancel • Start typing to filter",
		styles.KeyStyle.Render("Enter"),
		styles.KeyStyle.Render("Esc"))
	instructions := styles.InfoStyle.Render(instructionsText)

	// Center the dialog content
	dialogContent := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		matchText,
		"",
		inputLabel,
		filterInput,
		"",
		helpLine,
		"",
		instructions,
	)

	// Create a bordered box around the dialog
	dialog := styles.DialogBoxStyle.Render(dialogContent)

	// Center the dialog on screen
	b.WriteString(lipgloss.Place(m.totalWidth, m.totalHeight-4,
		lipgloss.Center, lipgloss.Center, dialog))

	return b.String()
}

// renderExportView renders the export configuration dialog
func renderExportView(m Model) string {
	var b strings.Builder

	// Header
	header := renderHeader(m)
	b.WriteString(header)
	b.WriteString("\n\n")

	// Export dialog content
	title := styles.TitleStyle.Render("Export Sessions")

	// Show what will be exported
	var exportInfo string
	if m.filterActive && m.filteredSessions != nil {
		exportInfo = fmt.Sprintf("Exporting %d filtered sessions (of %d total)",
			len(m.filteredSessions), len(m.sessions))
	} else {
		exportInfo = fmt.Sprintf("Exporting all %d sessions", len(m.sessions))
	}

	exportText := styles.MutedTextStyle.Render(exportInfo)

	// Format selection
	formatLabel := styles.BoldStyle.Render("Format:")
	var formatDisplay string
	if m.exportFormat == "csv" {
		formatDisplay = styles.SuccessStyle.Render("● CSV") + " " + styles.MutedTextStyle.Render("○ JSON")
	} else {
		formatDisplay = styles.MutedTextStyle.Render("○ CSV") + " " + styles.SuccessStyle.Render("● JSON")
	}

	// Filename input
	filenameLabel := styles.BoldStyle.Render("Filename:")
	filenameInput := styles.InputFocusedStyle.Render(m.exportFilename.View())
	extensionHint := styles.MutedTextStyle.Render("(." + m.exportFormat + " will be added automatically)")

	// Instructions
	instructionsText := fmt.Sprintf("%s to export • %s to cancel • %s/%s to toggle format",
		styles.KeyStyle.Render("Enter"),
		styles.KeyStyle.Render("Esc"),
		styles.KeyStyle.Render("Tab"),
		styles.KeyStyle.Render("Shift+Tab"))
	instructions := styles.InfoStyle.Render(instructionsText)

	// Center the dialog content
	dialogContent := lipgloss.JoinVertical(lipgloss.Left,
		title,
		"",
		exportText,
		"",
		formatLabel,
		formatDisplay,
		"",
		filenameLabel,
		filenameInput,
		extensionHint,
		"",
		instructions,
	)

	// Create a bordered box around the dialog
	dialog := styles.DialogBoxStyle.Render(dialogContent)

	// Center the dialog on screen
	b.WriteString(lipgloss.Place(m.totalWidth, m.totalHeight-4,
		lipgloss.Center, lipgloss.Center, dialog))

	return b.String()
}
