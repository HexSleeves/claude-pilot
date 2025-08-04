package cmd

import (
	"fmt"
	"strings"

	"claude-pilot/core/api"
	"claude-pilot/internal/cli"
	"claude-pilot/shared/components"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all active Claude sessions",
	Long: `List all active Claude coding sessions with their details.
Shows session ID, name, status, creation time, last activity, message count, and pane count.

Examples:
  claude-pilot list           	# List all sessions
  claude-pilot list --sort=name # Sort by name instead of last activity
	claude-pilot list --active 		# Show only active sessions
	claude-pilot list --inactive 	# Show only inactive sessions
	claude-pilot list --output=json # Output in JSON format
	claude-pilot list --quiet       # Output only session IDs`,
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize common command context
		ctx, err := InitializeCommand()
		if err != nil {
			return fmt.Errorf("failed to initialize command: %w", err)
		}

		// Get flags
		sortBy, _ := cmd.Flags().GetString("sort")
		active, _ := cmd.Flags().GetBool("active")
		inactive, _ := cmd.Flags().GetBool("inactive")
		idFilter, _ := cmd.Flags().GetString("id")
		jsonFlag, _ := cmd.Flags().GetBool("json")

		// Handle deprecation warnings
		if jsonFlag {
			err := CreateDeprecationWarning(ctx, "--json", "--output=json")
			if err != nil {
				return err
			}
		}

		// Handle positional arg deprecation
		if len(args) > 0 && idFilter == "" {
			idFilter = args[0]
			err := CreateDeprecationWarning(ctx, "positional session ID", "--id flag")
			if err != nil {
				return err
			}
		}

		// Validate mutually exclusive flags
		if active && inactive {
			return cli.NewValidationError("--active and --inactive flags are mutually exclusive", "Use either --active or --inactive, not both")
		}

		var sessions []*api.Session

		// Apply filters
		if active || inactive {
			filter := "active"
			if inactive {
				filter = "inactive"
			}

			sessions, err = ctx.Client.ListFilteredSessions(filter)
			if err != nil {
				return fmt.Errorf("failed to list filtered sessions: %w", err)
			}
		} else {
			// Get all sessions
			sessions, err = ctx.Client.ListSessions()
			if err != nil {
				return fmt.Errorf("failed to list sessions: %w", err)
			}
		}

		// Apply ID filter if specified
		if idFilter != "" {
			filteredSessions := make([]*api.Session, 0)
			for _, sess := range sessions {
				if strings.Contains(sess.ID, idFilter) || strings.Contains(sess.Name, idFilter) {
					filteredSessions = append(filteredSessions, sess)
				}
			}
			sessions = filteredSessions
		}

		// Handle empty results
		if len(sessions) == 0 {
			// Empty result is not an error - return success
			if ctx.OutputWriter.GetFormat() != cli.OutputFormatQuiet {
				var suggestions []string
				if active || inactive {
					suggestions = []string{
						"claude-pilot list",
						"claude-pilot list --active",
						"claude-pilot list --inactive",
						"claude-pilot create [session-name]",
					}
				} else {
					suggestions = []string{
						"claude-pilot create [session-name]",
					}
				}

				err := WriteHelpfulMessage(ctx, "No sessions found.", suggestions)
				if err != nil {
					return err
				}
			}
			return nil
		}

		// Convert API sessions to CLI output format
		sessionData := convertToSessionData(sessions, ctx.Client)

		// Apply sorting if requested
		direction := "asc"
		if sortBy != "" {
			if sortBy == "activity" {
				sortBy = "last_active"
				direction = "desc" // Most recent first for activity
			}
			// Note: Sorting will be handled by OutputWriter for structured formats
			// For human/table formats, we'll use the existing table component
		}

		// Use structured output for non-human formats
		if ctx.OutputWriter.GetFormat() != cli.OutputFormatHuman && ctx.OutputWriter.GetFormat() != cli.OutputFormatTable {
			// Convert to CLI output format
			cliSessions := make([]cli.SessionData, len(sessionData))
			for i, sess := range sessionData {
				cliSessions[i] = cli.SessionData{
					ID:        sess.ID,
					Name:      sess.Name,
					Project:   sess.ProjectPath,
					Status:    sess.Status,
					CreatedAt: sess.Created,
					UpdatedAt: sess.LastActive,
					PaneCount: sess.Panes,
				}
			}

			// Prepare metadata
			metadata := map[string]string{
				"backend": ctx.Client.GetBackend(),
				"total":   fmt.Sprintf("%d", len(sessions)),
			}
			if sortBy != "" {
				metadata["sortBy"] = sortBy
				metadata["sortDirection"] = direction
			}
			if active {
				metadata["filter"] = "active"
			} else if inactive {
				metadata["filter"] = "inactive"
			}

			return ctx.OutputWriter.WriteSessionList(cliSessions, metadata)
		}

		// Use existing table component for human/table output
		table := components.NewSessionTable(components.TableConfig{
			ShowHeaders: true,
			Interactive: false,
			MaxRows:     0, // Show all rows
			SortEnabled: true,
		})

		// Set the session data
		table.SetSessionData(sessionData)

		// Apply CLI sort option using table's built-in sorting
		if sortBy != "" {
			direction := "asc"
			if sortBy == "activity" {
				sortBy = "last_active"
				direction = "desc" // Most recent first for activity
			}
			err := table.SetSort(sortBy, direction)
			if err != nil {
				return fmt.Errorf("failed to set sort: %w", err)
			}
		}

		// Display sessions table using shared component
		if err := ctx.OutputWriter.WriteString("Claude Pilot Sessions\n"); err != nil {
			return err
		}
		if err := ctx.OutputWriter.WriteString(fmt.Sprintf("Backend: %s\n\n", ctx.Client.GetBackend())); err != nil {
			return err
		}
		if err := ctx.OutputWriter.WriteString(table.RenderCLI() + "\n"); err != nil {
			return err
		}

		// Show summary for human output
		activeCount := 0
		inactiveCount := 0
		for _, sess := range sessions {
			if sess.Status == api.StatusActive || sess.Status == api.StatusConnected {
				activeCount++
			} else {
				inactiveCount++
			}
		}

		summary := fmt.Sprintf("\nTotal: %d sessions (%d active, %d inactive)\n", len(sessions), activeCount, inactiveCount)
		if err := ctx.OutputWriter.WriteString(summary); err != nil {
			return err
		}

		// Show helpful commands
		suggestions := []string{
			"claude-pilot attach <session-name>",
			"claude-pilot kill <session-name>",
			"claude-pilot create [session-name]",
		}
		return WriteHelpfulMessage(ctx, "Available commands:", suggestions)
	},
}

// convertToSessionData converts API sessions to the shared table SessionData format
func convertToSessionData(sessions []*api.Session, client *api.Client) []components.SessionData {
	sessionData := make([]components.SessionData, len(sessions))

	for i, sess := range sessions {
		// Get pane count for active sessions
		paneCount := 0
		if sess.Status == api.StatusActive || sess.Status == api.StatusConnected {
			if count, err := client.GetSessionPaneCount(sess.Name); err == nil {
				paneCount = count
			}
			// If error getting pane count, just use 0 (don't fail the entire list)
		}

		sessionData[i] = components.SessionData{
			ID:          sess.ID,
			Name:        sess.Name,
			Status:      string(sess.Status),
			Backend:     sess.Backend,
			Created:     sess.CreatedAt,
			LastActive:  sess.LastActive,
			ProjectPath: sess.ProjectPath,
			Panes:       paneCount,
		}
	}

	return sessionData
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Add flags
	listCmd.Flags().BoolP("active", "a", false, "Show only active sessions")
	listCmd.Flags().Bool("inactive", false, "Show only inactive sessions")
	listCmd.Flags().StringP("sort", "s", "activity", "Sort by: name, created, status, activity, panes")
	listCmd.Flags().String("id", "", "Filter sessions by ID or name")

	// Add deprecated --json flag with deprecation notice
	listCmd.Flags().Bool("json", false, "Output in JSON format (deprecated: use --output=json)")
	// listCmd.Flags().MarkDeprecated("json", "use --output=json instead")
}
