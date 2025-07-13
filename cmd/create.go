package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"claude-pilot/internal/config"
	"claude-pilot/internal/manager"
	"claude-pilot/internal/ui"

	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [session-name]",
	Short: "Create a new Claude session",
	Long: `Create a new Claude coding session with an optional name.
If no name is provided, a timestamp-based name will be generated.

Examples:
  claude-pilot create                    # Create session with auto-generated name
  claude-pilot create my-project         # Create session named "my-project"
  claude-pilot create --desc "React app" # Create session with description
  claude-pilot create --project ./src    # Create session with project path`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var sessionName string
		if len(args) > 0 {
			sessionName = args[0]
		}

		// Get flags
		description, _ := cmd.Flags().GetString("description")
		projectPath, _ := cmd.Flags().GetString("project")

		// If project path is not provided, use current directory
		if projectPath == "" {
			cwd, err := os.Getwd()
			if err == nil {
				projectPath = cwd
			}
		} else {
			// Convert to absolute path
			absPath, err := filepath.Abs(projectPath)
			if err == nil {
				projectPath = absPath
			}
		}

		// Load configuration
		configManager := config.NewConfigManager("")
		cfg, err := configManager.Load()
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to load configuration: %v", err)))
			os.Exit(1)
		}

		// Create session manager
		sm, err := manager.NewSessionManager(cfg)
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to initialize session manager: %v", err)))
			os.Exit(1)
		}

		// Create the session
		sess, err := sm.CreateSession(sessionName, description, projectPath)
		if err != nil {
			fmt.Println(ui.ErrorMsg(fmt.Sprintf("Failed to create session: %v", err)))
			os.Exit(1)
		}

		// Success message
		fmt.Println(ui.SuccessMsg(fmt.Sprintf("Created session '%s'", sess.Name)))
		fmt.Println()

		// Show session details
		fmt.Printf("%-15s %s\n", ui.Bold("ID:"), sess.ID)
		fmt.Printf("%-15s %s\n", ui.Bold("Name:"), ui.Title(sess.Name))
		fmt.Printf("%-15s %s\n", ui.Bold("Status:"), ui.FormatStatus(string(sess.Status)))
		fmt.Printf("%-15s %s\n", ui.Bold("Backend:"), sm.GetConfig().Backend)
		fmt.Printf("%-15s %s\n", ui.Bold("Created:"), sess.CreatedAt.Format("2006-01-02 15:04:05"))
		if sess.ProjectPath != "" {
			fmt.Printf("%-15s %s\n", ui.Bold("Project:"), sess.ProjectPath)
		}
		if sess.Description != "" {
			fmt.Printf("%-15s %s\n", ui.Bold("Description:"), sess.Description)
		}

		fmt.Println()

		// Show next steps
		fmt.Println(ui.InfoMsg("Next steps:"))
		fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight(fmt.Sprintf("claude-pilot attach %s", sess.Name)))
		fmt.Printf("  %s %s\n", ui.Arrow(), ui.Highlight("claude-pilot list"))
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Add flags
	createCmd.Flags().StringP("description", "d", "", "Description for the session")
	createCmd.Flags().StringP("project", "p", "", "Project path for the session (defaults to current directory)")
}
