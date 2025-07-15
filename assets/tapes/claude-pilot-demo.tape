# Claude Pilot CLI Demo
# A demonstration of the claude-pilot session management tool

# Output settings
Output claude-pilot-demo.gif
Output claude-pilot-demo.mp4

# Terminal appearance
Set FontSize 14
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set Padding 20
Set Margin 40
Set MarginFill "#1e1e2e"
Set BorderRadius 8
Set WindowBar Colorful

# Typing speed
Set TypingSpeed 50ms

# Required programs
#Require claude-pilot
Require tmux

# Start with a clean terminal
Hide
Type "clear"
Enter
Show

# Introduction
Type "# Claude Pilot - Multi-Session Claude Code Manager"
Enter
Sleep 1s
Type "# Let's explore the key features!"
Enter
Sleep 2s

# Clear and show help
Type "clear && ./claude-pilot --help"
Enter
Sleep 3s

# List initial sessions (likely empty)
Type "./claude-pilot list"
Enter
Sleep 2s

# Create first session with a name
Type "./claude-pilot create react-frontend"
Enter
Sleep 3s

# Create another session
Type "./claude-pilot create api-backend"
Enter
Sleep 3s

# Create a session without a name (random name)
Type "./claude-pilot create"
Enter
Sleep 3s

# List all sessions to show the table
Type "./claude-pilot list"
Enter
Sleep 4s

# Demonstrate attaching to a session
Type "# Let's attach to our react-frontend session"
Enter
Type "./claude-pilot attach react-frontend"
Enter
Sleep 2s

# Show we're inside tmux with Claude
Type "# We're now inside a Claude Code session!"
Enter
Type "# This is where you'd interact with Claude for your React project"
Enter
Sleep 2s

# Detach from tmux session (Ctrl+B, D)
Ctrl+B
Type "d"
Sleep 1s

# Back to main terminal
Type "# Detached! Session is still running in background"
Enter
Sleep 2s

# Show sessions again
Type "./claude-pilot list"
Enter
Sleep 3s

# Kill a specific session
Type "./claude-pilot kill api-backend"
Enter
Sleep 2s

# Show updated list
Type "./claude-pilot list"
Enter
Sleep 3s

# Demonstrate kill-all (but first create more sessions)
Type "./claude-pilot create test-session-1"
Enter
Sleep 1s
Type "./claude-pilot create test-session-2"
Enter
Sleep 1s

# Show all sessions
Type "./claude-pilot list"
Enter
Sleep 3s

# Kill all sessions
Type "./claude-pilot kill-all"
Enter
Sleep 2s

# Final list showing empty state
Type "./claude-pilot list"
Enter
Sleep 3s

# Closing message
Type "# Claude Pilot: Manage multiple Claude Code sessions with ease!"
Enter
Type "# Perfect for multi-project workflows 🚀"
Enter
Sleep 3s
