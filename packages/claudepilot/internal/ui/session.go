package ui

import (
	"claude-pilot/core/api"
	"fmt"
)

// DisplaySessionDetails shows formatted session information consistently
// This eliminates the duplicated session detail display logic in create.go and kill.go
func DisplaySessionDetails(sess *api.Session, backend string) {
	fmt.Printf("%-15s %s\n", Bold("ID:"), sess.ID)
	fmt.Printf("%-15s %s\n", Bold("Name:"), Title(sess.Name))
	fmt.Printf("%-15s %s\n", Bold("Status:"), FormatStatus(string(sess.Status)))
	fmt.Printf("%-15s %s\n", Bold("Backend:"), backend)
	fmt.Printf("%-15s %s\n", Bold("Created:"), sess.CreatedAt.Format("2006-01-02 15:04:05"))
	if sess.ProjectPath != "" {
		fmt.Printf("%-15s %s\n", Bold("Project:"), sess.ProjectPath)
	}
	if sess.Description != "" {
		fmt.Printf("%-15s %s\n", Bold("Description:"), sess.Description)
	}
}

// DisplaySessionDetailsWithMessages shows session details including message count
// This is used in kill.go where message count is displayed
func DisplaySessionDetailsWithMessages(sess *api.Session, backend string) {
	DisplaySessionDetails(sess, backend)
	fmt.Printf("%-15s %d\n", Bold("Messages:"), len(sess.Messages))
}

// DisplayAvailableSessions shows available sessions in error scenarios
// This eliminates the duplicated session listing logic in attach.go
func DisplayAvailableSessions(sessions []*api.Session) {
	if len(sessions) == 0 {
		fmt.Println(Dim("  No sessions available"))
		fmt.Printf("  %s %s\n", Arrow(), Highlight("claude-pilot create [session-name]"))
	} else {
		for _, s := range sessions {
			idDisplay := s.ID
			if len(s.ID) > 8 {
				idDisplay = s.ID[:8]
			}

			fmt.Printf("  %s %s (%s)\n", Arrow(), Highlight(s.Name), Dim(idDisplay))
		}
	}
}

// DisplayNextSteps shows helpful next commands consistently
// This eliminates the duplicated "next steps" logic across multiple commands
func DisplayNextSteps(commands ...string) {
	fmt.Println()
	fmt.Println(NextSteps(commands...))
}

// DisplayAvailableCommands shows helpful commands consistently
// This eliminates the duplicated "available commands" logic in list.go
func DisplayAvailableCommands(commands ...string) {
	fmt.Println()
	fmt.Println(AvailableCommands(commands...))
}

// DisplaySessionSummary shows session count summary
// This eliminates the duplicated summary logic in list.go and kill.go
func DisplaySessionSummary(totalSessions, activeSessions, inactiveSessions int, showAll bool) {
	if showAll {
		fmt.Printf("%s Total: %d sessions (%d active, %d inactive)\n",
			InfoMsg("Summary:"), totalSessions, activeSessions, inactiveSessions)
	} else {
		fmt.Printf("%s Active sessions: %d\n",
			InfoMsg("Summary:"), activeSessions)
		if inactiveSessions > 0 {
			fmt.Printf("  %s Use --all to show %d inactive sessions\n",
				Dim("Note:"), inactiveSessions)
		}
	}
}

// DisplayRemainingSessionsInfo shows information about remaining sessions after deletion
// This eliminates the duplicated remaining sessions logic in kill.go
func DisplayRemainingSessionsInfo(remainingSessions []*api.Session) {
	if len(remainingSessions) > 0 {
		fmt.Printf("%s %d sessions remaining\n", InfoMsg("Status:"), len(remainingSessions))
		fmt.Printf("  %s %s\n", Arrow(), Highlight("claude-pilot list"))
	} else {
		fmt.Println(InfoMsg("No sessions remaining"))
		fmt.Printf("  %s %s\n", Arrow(), Highlight("claude-pilot create [session-name]"))
	}
}
