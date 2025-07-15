# Requirements Document

## Introduction

The TUI (Terminal User Interface) Implementation feature will provide an interactive, visual interface for managing Claude Pilot sessions. This will complement the existing CLI interface by offering a more intuitive, dashboard-style experience for users who prefer graphical interaction over command-line operations. The TUI will be built using Bubble Tea framework and will provide real-time session monitoring, interactive session management, and a responsive design that adapts to different terminal sizes.

## Requirements

### Requirement 1

**User Story:** As a developer, I want to launch an interactive TUI interface so that I can visually manage my Claude sessions without memorizing CLI commands.

#### Acceptance Criteria

1. WHEN I run `claude-pilot tui` THEN the system SHALL launch a full-screen terminal user interface
2. WHEN I run `claude-pilot` with `ui.mode` set to "tui" in config THEN the system SHALL automatically launch the TUI interface
3. WHEN the TUI launches THEN the system SHALL display a dashboard with session overview and navigation options
4. WHEN I press 'q' or 'Ctrl+C' THEN the system SHALL gracefully exit the TUI and return to the terminal

### Requirement 2

**User Story:** As a developer, I want to see all my sessions in a visual table so that I can quickly understand the status and details of each session.

#### Acceptance Criteria

1. WHEN the TUI loads THEN the system SHALL display a table showing all sessions with columns for ID, Name, Status, Created, and Last Active
2. WHEN sessions are active THEN the system SHALL highlight them with green status indicators
3. WHEN sessions are inactive THEN the system SHALL show them with yellow status indicators
4. WHEN sessions have errors THEN the system SHALL display them with red status indicators
5. WHEN I use arrow keys THEN the system SHALL allow me to navigate between table rows
6. WHEN I press Enter on a selected session THEN the system SHALL show detailed session information

### Requirement 3

**User Story:** As a developer, I want to create new sessions through the TUI so that I can quickly start new Claude conversations without switching to CLI.

#### Acceptance Criteria

1. WHEN I press 'c' or 'n' THEN the system SHALL open a session creation modal
2. WHEN the creation modal opens THEN the system SHALL provide input fields for session name and optional description
3. WHEN I enter a session name and press Enter THEN the system SHALL create the session and update the session table
4. WHEN I press Escape in the modal THEN the system SHALL cancel creation and return to the main view
5. WHEN session creation fails THEN the system SHALL display an error message and allow retry

### Requirement 4

**User Story:** As a developer, I want to attach to sessions directly from the TUI so that I can seamlessly transition from management to active coding.

#### Acceptance Criteria

1. WHEN I select a session and press 'a' or Enter THEN the system SHALL attach to that session
2. WHEN attaching to a session THEN the system SHALL exit the TUI and connect to the tmux/zellij session
3. WHEN I select an inactive session and try to attach THEN the system SHALL display a warning and offer to restart the session
4. WHEN attachment fails THEN the system SHALL show an error message and remain in the TUI

### Requirement 5

**User Story:** As a developer, I want to terminate sessions from the TUI so that I can clean up unused sessions without using CLI commands.

#### Acceptance Criteria

1. WHEN I select a session and press 'd' or 'Delete' THEN the system SHALL prompt for confirmation before terminating
2. WHEN I confirm termination THEN the system SHALL kill the session and remove it from the table
3. WHEN I press 'D' (shift+d) THEN the system SHALL offer to terminate all sessions with confirmation
4. WHEN termination fails THEN the system SHALL display an error message and keep the session in the table

### Requirement 6

**User Story:** As a developer, I want the TUI to be responsive to different terminal sizes so that I can use it effectively on various screen configurations.

#### Acceptance Criteria

1. WHEN the terminal width is less than 80 characters THEN the system SHALL use a compact layout with abbreviated columns
2. WHEN the terminal width is between 80-120 characters THEN the system SHALL use a balanced layout with standard columns
3. WHEN the terminal width is greater than 120 characters THEN the system SHALL use a full layout with all details visible
4. WHEN I resize the terminal THEN the system SHALL automatically adjust the layout without losing functionality
5. WHEN the terminal is too small to display content THEN the system SHALL show a message indicating minimum size requirements

### Requirement 7

**User Story:** As a developer, I want real-time updates in the TUI so that I can see session status changes without manually refreshing.

#### Acceptance Criteria

1. WHEN sessions change status THEN the system SHALL automatically update the display within 5 seconds
2. WHEN new sessions are created externally THEN the system SHALL add them to the table automatically
3. WHEN sessions are terminated externally THEN the system SHALL remove them from the table automatically
4. WHEN I press 'r' THEN the system SHALL manually refresh all session data immediately

### Requirement 8

**User Story:** As a developer, I want keyboard shortcuts and help information so that I can efficiently navigate and use the TUI.

#### Acceptance Criteria

1. WHEN I press '?' or 'h' THEN the system SHALL display a help modal with all available keyboard shortcuts
2. WHEN I press 'j'/'k' or arrow keys THEN the system SHALL navigate up/down in the session table
3. WHEN I press '/' THEN the system SHALL open a search/filter input for finding specific sessions
4. WHEN I press Escape THEN the system SHALL close any open modals or return to the main view
5. WHEN I press Tab THEN the system SHALL cycle between different UI panels if multiple panels are visible

### Requirement 9

**User Story:** As a developer, I want to see detailed session information so that I can understand session context and history.

#### Acceptance Criteria

1. WHEN I select a session and press 'i' THEN the system SHALL show a detail panel with full session metadata
2. WHEN viewing session details THEN the system SHALL display creation time, last active time, project path, and description
3. WHEN a session has a long description THEN the system SHALL wrap text appropriately and provide scrolling if needed
4. WHEN I press Escape in the detail view THEN the system SHALL return to the main session table

### Requirement 10

**User Story:** As a developer, I want the TUI to maintain consistent theming with the CLI so that the experience feels cohesive across interfaces.

#### Acceptance Criteria

1. WHEN the TUI displays THEN the system SHALL use the same Claude orange theme colors as the CLI
2. WHEN showing status indicators THEN the system SHALL use consistent colors (green for active, yellow for inactive, red for errors)
3. WHEN displaying text THEN the system SHALL follow the same typography hierarchy as defined in the shared styles
4. WHEN rendering borders and panels THEN the system SHALL use the established design patterns from the theming standards