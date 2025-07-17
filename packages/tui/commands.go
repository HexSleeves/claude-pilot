package main

import (
	"claude-pilot/core/api"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// loadSessionsCmd loads all sessions from the API
func loadSessionsCmd(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return sessionsLoadedMsg{
				sessions: nil,
				err:      fmt.Errorf("API client is nil"),
			}
		}

		sessions, err := client.ListSessions()
		return sessionsLoadedMsg{
			sessions: sessions,
			err:      err,
		}
	}
}

// createSessionCmd creates a new session with the specified parameters
func createSessionCmd(client *api.Client, name, description, projectPath string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return sessionCreatedMsg{
				session: nil,
				err:     fmt.Errorf("API client is nil"),
			}
		}

		if strings.TrimSpace(name) == "" {
			return sessionCreatedMsg{
				session: nil,
				err:     fmt.Errorf("session name cannot be empty"),
			}
		}

		req := api.CreateSessionRequest{
			Name:        strings.TrimSpace(name),
			Description: strings.TrimSpace(description),
			ProjectPath: strings.TrimSpace(projectPath),
		}

		session, err := client.CreateSession(req)
		return sessionCreatedMsg{
			session: session,
			err:     err,
		}
	}
}

// killSessionCmd terminates a session by ID
func killSessionCmd(client *api.Client, sessionID string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return sessionKilledMsg{
				sessionID: sessionID,
				err:       fmt.Errorf("API client is nil"),
			}
		}

		if strings.TrimSpace(sessionID) == "" {
			return sessionKilledMsg{
				sessionID: sessionID,
				err:       fmt.Errorf("session ID cannot be empty"),
			}
		}

		err := client.KillSession(sessionID)
		return sessionKilledMsg{
			sessionID: sessionID,
			err:       err,
		}
	}
}

// attachSessionCmd attaches to a session and hands control to the multiplexer
func attachSessionCmd(client *api.Client, sessionID string) tea.Cmd {
	return func() tea.Msg {
		if client == nil {
			return errorMsg{error: fmt.Errorf("API client is nil")}
		}

		if strings.TrimSpace(sessionID) == "" {
			return errorMsg{error: fmt.Errorf("session ID cannot be empty")}
		}

		// Get the session to find its name
		session, err := client.GetSession(sessionID)
		if err != nil {
			return errorMsg{error: fmt.Errorf("failed to get session: %w", err)}
		}

		if session == nil {
			return errorMsg{error: fmt.Errorf("session not found")}
		}

		if strings.TrimSpace(session.Name) == "" {
			return errorMsg{error: fmt.Errorf("session name is empty")}
		}

		// Create a command to attach to the session
		// This will hand control over to the multiplexer session
		cmd := exec.Command("tmux", "attach-session", "-t", session.Name)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		// Use tea.ExecProcess to hand control to the external process
		return tea.ExecProcess(cmd, func(err error) tea.Msg {
			if err != nil {
				return errorMsg{error: fmt.Errorf("failed to attach to session: %w", err)}
			}
			// After the session ends, quit the TUI
			return tea.Quit()
		})
	}
}
