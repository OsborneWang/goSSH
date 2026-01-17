# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

GoSSH is a cross-platform SSH command-line tool written in Go that manages SSH server configurations, connections, remote command execution, and file transfers. It supports Windows, Linux, and macOS with intelligent terminal detection for opening SSH sessions in new tabs/windows.

## Build and Development Commands

### Building
```bash
# Build for current platform
make build          # Creates goss.exe (Windows) or goss (Unix)
go build -o goss .  # Alternative direct build

# Cross-platform builds
make build-linux    # Linux AMD64
make build-macos    # macOS AMD64
```

### Running During Development
```bash
# Quick testing without building
go run main.go [command] [args]
go run main.go interactive
go run main.go connect server1

# Using Makefile shortcuts
make run                # Show help
make run-interactive    # Interactive mode
make run-list          # List servers
```

### Hot Reload Development
```bash
# Install Air (if not installed)
make install-tools

# Start hot reload mode (auto-recompile on file changes)
make watch
# or
air
```

### Code Quality
```bash
make fmt     # Format code with go fmt and goimports
make vet     # Run go vet for static analysis
make lint    # Run golangci-lint (requires installation)
make test    # Run tests
```

## Architecture

### Project Structure
```
goSSH/
├── main.go                      # Entry point, calls cmd.Execute()
├── cmd/                         # Cobra command definitions
│   ├── root.go                  # Root command setup
│   ├── add.go                   # Add server configuration
│   ├── list.go                  # List servers
│   ├── remove.go                # Remove server
│   ├── connect.go               # SSH connection (supports --no-new-tab flag)
│   ├── exec.go                  # Remote command execution
│   ├── transfer.go              # File upload/download via SFTP
│   └── interactive.go           # Interactive menu mode
├── internal/
│   ├── config/                  # Configuration management layer
│   │   └── config.go            # Manager wraps Storage, provides server CRUD
│   ├── ssh/                     # SSH functionality
│   │   ├── client.go            # SSH client wrapper (Connect, Close, Reconnect)
│   │   ├── executor.go          # Windows-specific (build tag: windows)
│   │   ├── executor_unix.go     # Unix-specific (build tag: !windows)
│   │   ├── terminal.go          # Terminal detection and new tab/window logic
│   │   └── transfer.go          # SFTP file transfer implementation
│   └── storage/                 # Persistent storage
│       └── storage.go           # JSON file I/O for server configs
└── models/
    └── server.go                # Server and ServerConfig structs
```

### Key Design Patterns

**Configuration Storage:**
- Server configs stored as JSON at `%APPDATA%\gossh\servers.json` (Windows) or `~/.config/gossh/servers.json` (Unix)
- Storage layer (internal/storage) handles file I/O
- Config layer (internal/config) provides business logic (CRUD operations)
- Passwords stored in plaintext (security consideration documented in README)

**SSH Client Architecture:**
- `ssh.Client` wraps `golang.org/x/crypto/ssh` connection
- `ssh.Executor` handles command execution and shell sessions
- Platform-specific terminal handling via build tags (`executor.go` for Windows, `executor_unix.go` for Unix)

**Terminal Detection System (internal/ssh/terminal.go):**
- `DetectTerminal()` identifies current terminal (Windows Terminal, iTerm2, GNOME Terminal, Konsole, etc.)
- `OpenInNewTab()` attempts to open SSH session in new tab of current terminal
- `OpenInNewWindow()` fallback for terminals without tab support
- `HasDesktopEnvironment()` checks for GUI availability (Linux)
- Each terminal type has specific implementation (AppleScript for macOS, wt.exe for Windows Terminal, etc.)

**Shell Execution Flow (connect command):**
1. User runs `goss connect server1`
2. `cmd/connect.go` checks `--no-new-tab` flag
3. If flag not set, `executor.ExecuteShell(true)` tries:
   - New tab via `ExecuteShellInNewTab()` → spawns `goss connect server1 --no-new-tab` in new tab
   - Falls back to new window via `ExecuteShellInNewWindow()`
   - Falls back to current terminal via `executeShellInCurrentTerminal()`
4. If `--no-new-tab` set, directly executes in current terminal (prevents recursion)

**Windows-Specific Handling:**
- `crlfFilterReader` in `executor.go` filters `\r` characters to prevent double carriage returns
- Terminal size detection uses `windows.GetConsoleScreenBufferInfo`
- SSH terminal modes set `ECHO: 0` to disable remote echo (prevents command duplication)

**Command Framework:**
- Uses `github.com/spf13/cobra` for CLI structure
- Interactive prompts via `github.com/manifoldco/promptui`
- Color output via `github.com/fatih/color`

## Important Implementation Notes

### When Modifying SSH Connection Logic
- Always test both `--no-new-tab` and default behavior to avoid recursion
- Terminal detection may fail gracefully; ensure fallback paths work
- Windows requires special CRLF handling in `crlfFilterReader`

### When Adding New Commands
- Add command file in `cmd/` directory
- Register in `init()` function with `rootCmd.AddCommand(yourCmd)`
- Use `config.NewManager()` to access server configurations
- Follow existing patterns for interactive vs direct argument modes

### When Modifying Terminal Detection
- Test across multiple terminal types (Windows Terminal, iTerm2, GNOME Terminal, etc.)
- Ensure `HasDesktopEnvironment()` correctly detects headless environments
- New tab/window functions should use `cmd.Start()` not `cmd.Run()` to avoid blocking

### Cross-Platform Considerations
- Use build tags for platform-specific code: `// +build windows` or `// +build !windows`
- Terminal size detection differs: Windows uses `windows.ConsoleScreenBufferInfo`, Unix uses `unix.TIOCGWINSZ`
- Path handling: Use `filepath.Join()` for cross-platform compatibility

## Testing

### Manual Testing Setup
The DEVELOPMENT.md file documents several approaches for testing SSH functionality:
- Local SSH server (OpenSSH on Windows 10+, Linux, macOS)
- Docker container: `docker run -d -p 2222:22 -e ROOT_PASSWORD=testpass123 panubo/sshd`
- Cloud servers or VMs

### Debugging
- VS Code debugging configurations available (see DEVELOPMENT.md for setup)
- Use `GOSSH_DEBUG=1` environment variable for debug logging (if implemented)
- Delve debugger: `dlv debug . -- connect server1`

## Dependencies

Key dependencies (from go.mod):
- `golang.org/x/crypto/ssh` - SSH protocol implementation
- `github.com/pkg/sftp` - SFTP file transfer
- `github.com/spf13/cobra` - CLI framework
- `github.com/manifoldco/promptui` - Interactive prompts
- `github.com/fatih/color` - Terminal colors
- `golang.org/x/sys/windows` - Windows system calls (terminal size)

## Configuration File Format

```json
{
  "servers": [
    {
      "name": "server1",
      "host": "192.168.1.100",
      "port": 22,
      "username": "root",
      "password": "plaintext_password"
    }
  ]
}
```

Note: Passwords are stored in plaintext. This is documented as a security consideration in README.md.
