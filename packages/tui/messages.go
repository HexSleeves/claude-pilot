package main

import (
	"claude-pilot/shared/interfaces"
)

// Message types for the bubbletea message-passing system

// sessionsLoadedMsg contains the result of loading sessions from the API.
// This message is sent when the loadSessionsCmd completes, either successfully
// with a list of sessions or with an error.
type sessionsLoadedMsg struct {
	sessions []*interfaces.Session
	err      error
}

// sessionCreatedMsg contains the result of creating a new session.
// This message is sent when the createSessionCmd completes, containing
// either the newly created session or an error if creation failed.
type sessionCreatedMsg struct {
	session *interfaces.Session
	err     error
}

// sessionKilledMsg contains the result of terminating a session.
// This message is sent when the killSessionCmd completes, indicating
// success or failure of the session termination operation.
type sessionKilledMsg struct {
	sessionID string
	err       error
}

// errorMsg contains error information for display in the error view.
// This message is used to transition the TUI to an error state with
// user-friendly error display and retry options.
type errorMsg struct {
	error error
}

// statusMsg contains status text for temporary display in the main view.
// This message is used to show brief status updates to the user, such as
// "Session created successfully" or "Refreshing sessions...".
type statusMsg struct {
	message string
}

// Internal state messages

// viewStateMsg triggers a change between different TUI views.
// This message is used internally to transition between TableView,
// CreatePrompt, Loading, Error, and Help states.
type viewStateMsg struct {
	state ViewState
}
