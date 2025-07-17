# Implementation Plan

- [x] 1. Fix missing dependencies and update shared package
  - Add missing Bubble Tea dependencies to shared package go.mod
  - Update shared/styles/theme.go to include missing Bubbles component styling functions
  - Fix import issues in claudepilot package related to shared/styles
  - _Requirements: 10.1, 10.2, 10.3, 10.4_

- [x] 2. Implement Summary Panel Model
  - [x] 2.1 Create SummaryPanelModel struct and basic methods
    - Write SummaryPanelModel with client, width, sessions, and stats fields
    - Implement Init(), Update(), and View() methods following Bubble Tea pattern
    - Create SessionStats struct with Total, Active, Inactive, Error, Backend fields
    - _Requirements: 1.1, 1.3, 7.1_

  - [x] 2.2 Implement session statistics calculation
    - Write calculateStats() method to analyze session data and compute statistics
    - Implement real-time statistics updates when session data changes
    - Add backend information display from client configuration
    - _Requirements: 7.1, 7.2_

  - [x] 2.3 Create responsive summary card layout
    - Implement renderSummaryCards() method using shared layout components
    - Create individual stat cards with color-coded status indicators
    - Add responsive behavior for different terminal sizes using styles.GetResponsiveWidth()
    - _Requirements: 6.1, 6.2, 6.3, 6.4_

- [x] 3. Implement Detail Panel Model
  - [x] 3.1 Create DetailPanelModel struct and basic functionality
    - Write DetailPanelModel with client, dimensions, session, scrollY, and content fields
    - Implement Init(), Update(), and View() methods for session detail display
    - Add SetSession() method to update displayed session information
    - _Requirements: 9.1, 9.2_

  - [x] 3.2 Implement scrollable content rendering
    - Write renderSessionDetails() method to format session metadata display
    - Implement vertical scrolling for long content using viewport management
    - Add scroll indicators and position display for user feedback
    - _Requirements: 9.3, 8.2_

  - [x] 3.3 Add session metadata formatting
    - Create formatSessionMetadata() method to display creation time, last active, project path
    - Implement message history display if session has conversation data
    - Add proper text wrapping and formatting for different field types
    - _Requirements: 9.2, 9.3_

- [x] 4. Enhance Session Table Model
  - [x] 4.1 Improve keyboard navigation and selection
    - Enhance existing navigation methods (moveUp, moveDown, pageUp, pageDown)
    - Add Home/End key support for jumping to first/last session
    - Implement proper viewport adjustment to keep selected item visible
    - _Requirements: 8.2, 8.3_

  - [x] 4.2 Add search and filtering functionality
    - Implement search input field that appears when '/' key is pressed
    - Write filterSessions() method to filter sessions by name, status, or project
    - Add search highlighting and result count display
    - _Requirements: 8.3_

  - [ ] 4.3 Enhance status display and real-time updates
    - Improve status formatting with consistent icons and colors using shared styles
    - Implement automatic status refresh every 5 seconds for real-time updates
    - Add visual indicators for recently updated sessions
    - _Requirements: 7.1, 7.2, 7.3, 10.2_

- [x] 5. Implement Create Modal Model
  - [x] 5.1 Create modal structure and form inputs
    - Write CreateModalModel with textinput fields for name and description
    - Implement modal lifecycle methods (Init, Update, View, Reset, IsCompleted)
    - Add form field navigation using Tab key and focus management
    - _Requirements: 3.1, 3.2, 8.4, 8.5_

  - [x] 5.2 Add input validation and error handling
    - Implement validateInput() method to check session name requirements
    - Add real-time validation feedback with error message display
    - Handle duplicate session name detection and user feedback
    - _Requirements: 3.3, 3.4_

  - [ ] 5.3 Implement session creation workflow
    - Write createSession() command to call core API for session creation
    - Handle creation success/failure with appropriate user feedback
    - Implement modal closing and data refresh after successful creation
    - _Requirements: 3.3, 3.5_

- [ ] 6. Enhance Dashboard Model coordination
  - [ ] 6.1 Implement comprehensive focus management
    - Enhance cycleFocus() method to handle all focusable components
    - Add visual focus indicators for each component type
    - Implement context-aware keyboard shortcuts based on focused component
    - _Requirements: 8.1, 8.2, 8.4, 8.5_

  - [ ] 6.2 Add modal overlay system
    - Enhance overlayModal() method to properly center modals on different screen sizes
    - Implement modal backdrop with proper z-index layering
    - Add modal animation or transition effects for better user experience
    - _Requirements: 3.1, 6.4_

  - [ ] 6.3 Implement responsive layout logic
    - Enhance renderMainContent() with improved responsive breakpoint handling
    - Add dynamic panel sizing based on available screen real estate
    - Implement panel hiding/showing logic for very small terminals
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 6.5_

- [ ] 7. Add session management actions
  - [ ] 7.1 Implement session attachment functionality
    - Write attachToSession() command that calls core API attach method
    - Handle attachment success by gracefully exiting TUI to CLI session
    - Add error handling for failed attachments with user feedback
    - _Requirements: 4.1, 4.2, 4.3, 4.4_

  - [ ] 7.2 Implement session termination with confirmation
    - Write killSession() command with confirmation dialog
    - Add bulk termination support for killing all sessions
    - Implement proper cleanup and session list refresh after termination
    - _Requirements: 5.1, 5.2, 5.3, 5.4_

  - [ ] 7.3 Add session restart functionality for inactive sessions
    - Detect inactive sessions and offer restart option in attachment flow
    - Implement restartSession() command to reactivate stopped sessions
    - Add user confirmation and error handling for restart operations
    - _Requirements: 4.3_

- [x] 8. Implement help system and keyboard shortcuts
  - [x] 8.1 Create help modal with keyboard shortcut reference
    - Write HelpModalModel to display comprehensive keyboard shortcut list
    - Implement context-sensitive help that shows relevant shortcuts
    - Add help modal toggle with '?' or 'h' key binding
    - _Requirements: 8.1, 8.4, 8.5_

  - [x] 8.2 Add status bar with current context and shortcuts
    - Enhance renderFooter() to show context-aware keyboard shortcuts
    - Display current selection information and available actions
    - Add dynamic shortcut hints based on focused component
    - _Requirements: 8.1, 8.5_

- [ ] 9. Implement real-time updates and refresh system
  - [ ] 9.1 Add automatic session data refresh
    - Implement periodic refresh timer that updates session data every 5 seconds
    - Write efficient refresh logic that only updates changed sessions
    - Add manual refresh capability with 'r' key binding
    - _Requirements: 7.1, 7.2, 7.4_

  - [ ] 9.2 Add external session change detection
    - Implement file system watching for session metadata changes
    - Handle external session creation/deletion with automatic UI updates
    - Add notification system for external changes
    - _Requirements: 7.2, 7.3_

- [ ] 10. Add comprehensive error handling and user feedback
  - [ ] 10.1 Implement error display system
    - Create ErrorDisplayModel for showing errors in different contexts
    - Add error severity levels (info, warning, critical) with appropriate styling
    - Implement error dismissal and automatic timeout for transient errors
    - _Requirements: 3.4, 4.4, 5.4_

  - [ ] 10.2 Add loading states and progress indicators
    - Implement loading spinners for long-running operations
    - Add progress feedback for session creation, attachment, and termination
    - Create loading overlays that don't block user interaction when possible
    - _Requirements: 1.1, 3.1, 4.1, 5.1_

- [ ] 11. Fix build issues and code cleanup
  - [ ] 11.1 Fix unused imports and build errors
    - Remove unused "strings" import from session_table.go
    - Fix any other build issues preventing compilation
    - Clean up unused variables and methods
    - _Requirements: All requirements_

  - [ ] 11.2 Add missing shared component integrations
    - Implement missing Bubbles component styling functions in shared/styles
    - Add ConfigureBubblesTable, ConfigureBubblesTextInput, ConfigureBubblesViewport, ConfigureBubblesHelp functions
    - Ensure all TUI models can properly use shared styling
    - _Requirements: 10.1, 10.2, 10.3, 10.4_

- [ ] 12. Implement comprehensive testing
  - [ ] 12.1 Write unit tests for all models
    - Create test files for each model with comprehensive test coverage
    - Mock core API client for isolated component testing
    - Test keyboard navigation, state transitions, and error handling
    - _Requirements: All requirements_

  - [ ] 12.2 Add integration tests for component interaction
    - Write tests for message passing between components
    - Test responsive layout behavior across different terminal sizes
    - Verify theming consistency and accessibility features
    - _Requirements: 6.1, 6.2, 6.3, 6.4, 10.1, 10.2, 10.3, 10.4_

- [ ] 13. Final integration and polish
  - [ ] 13.1 Integrate TUI with CLI command structure
    - Update root command to properly launch TUI when ui.mode is "tui"
    - Add `claude-pilot tui` subcommand for explicit TUI launch
    - Ensure proper configuration sharing between CLI and TUI modes
    - _Requirements: 1.1, 1.2_

  - [ ] 13.2 Add performance optimizations
    - Implement lazy loading for session details and large data sets
    - Add viewport rendering optimizations for large session lists
    - Optimize refresh cycles to minimize unnecessary re-renders
    - _Requirements: 6.4, 7.1_

  - [ ] 13.3 Final testing and documentation
    - Perform comprehensive manual testing across different terminal sizes
    - Test with large numbers of sessions (50+) for performance validation
    - Update README and documentation with TUI usage instructions
    - _Requirements: All requirements_