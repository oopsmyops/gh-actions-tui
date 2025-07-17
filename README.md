# GitHub Actions TUI

A terminal user interface (TUI) for viewing GitHub Actions workflows, runs, jobs, and logs. Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and powered by the [GitHub CLI](https://cli.github.com/).

![GitHub Actions TUI Demo](https://via.placeholder.com/800x400/1a1a1a/ffffff?text=GitHub+Actions+TUI+Demo)

## Features

- 🔍 **Browse Workflows**: View all workflows in a repository
- 🏃 **Workflow Runs**: See recent runs with status indicators
- 📋 **Job Details**: Inspect individual jobs and their status
- 📜 **Log Viewer**: Full-featured log viewer with `less`-like navigation
- 🔎 **Search Logs**: Search through logs with `/`, navigate with `n`/`N`
- ⌨️ **Vim-like Navigation**: `gg` to top, `G` to bottom, `j`/`k` for movement
- 🎨 **Clean Interface**: Intuitive TUI with helpful key bindings

## Installation

### Prerequisites

- [GitHub CLI](https://cli.github.com/) (v2.75+) - Required for fetching GitHub data
- Git (for cloning the repository)

### Install GitHub CLI

```bash
# macOS
brew install gh

# Ubuntu/Debian
sudo apt install gh

# Windows
winget install GitHub.cli

# Or download from: https://cli.github.com/
```

### Authenticate with GitHub

```bash
gh auth login
```

### Install gh-actions-tui

#### From Source

```bash
git clone https://github.com/yourusername/gh-actions-tui.git
cd gh-actions-tui
go build -o gh-actions-tui .
```

#### From Releases (Coming Soon)

Download the latest binary from the [releases page](https://github.com/yourusername/gh-actions-tui/releases).

## Usage

### Basic Usage

```bash
# View workflows for a specific repository
./gh-actions-tui -r owner/repo

# Example
./gh-actions-tui -r dhth/act3
```

### Navigation

#### Main Navigation
- `↑/↓` or `j/k` - Navigate lists
- `Enter` - Select item / Go deeper
- `Esc` - Go back / Exit current view
- `q` or `Ctrl+C` - Quit application

#### Log Viewer (less-like commands)
- `↑/↓` or `j/k` - Scroll line by line
- `Page Up/Down` or `b/f` - Scroll page by page
- `u/d` - Half page up/down
- `Space` - Page down
- `gg` - Go to top
- `G` - Go to bottom
- `Home/End` - Go to top/bottom

#### Search in Logs
- `/` - Start search
- `Enter` - Execute search
- `n` - Next match
- `N` - Previous match
- `Esc` - Exit search mode

## Interface Overview

The application has four main views:

1. **Workflows View** - List all workflows in the repository
2. **Runs View** - Show recent runs for selected workflow
3. **Jobs View** - Display jobs for selected run
4. **Logs View** - Full log viewer with search capabilities

### Status Indicators

- ✅ Success
- ❌ Failure  
- 🚫 Cancelled
- ⏭️ Skipped
- 🔄 In Progress
- ⏳ Queued

## Examples

### Viewing Workflow Logs

```bash
# Start the TUI
./gh-actions-tui -r microsoft/vscode

# Navigate through:
# 1. Select a workflow (e.g., "CI")
# 2. Select a recent run
# 3. Select a job (e.g., "build")
# 4. View logs with full navigation
```

### Searching Logs

1. Navigate to logs view
2. Press `/` to start search
3. Type your search term (e.g., "error")
4. Press `Enter` to search
5. Use `n`/`N` to navigate between matches

## Development

### Prerequisites

- Go 1.21+
- GitHub CLI configured

### Building

```bash
git clone https://github.com/yourusername/gh-actions-tui.git
cd gh-actions-tui
go mod tidy
go build -o gh-actions-tui .
```

### Running Tests

```bash
go test ./...
```

### Project Structure

```
gh-actions-tui/
├── main.go           # Application entry point
├── tui/              # TUI components
│   ├── model.go      # Bubble Tea model
│   ├── view.go       # UI rendering
│   └── keys.go       # Key bindings
├── github/           # GitHub API client
│   └── client.go     # GitHub CLI wrapper
└── README.md
```

## Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [GitHub CLI Go](https://github.com/cli/go-gh) - GitHub CLI integration

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Troubleshooting

### Common Issues

**"gh: command not found"**
- Install GitHub CLI: https://cli.github.com/

**"authentication required"**
- Run: `gh auth login`

**"failed to list workflows"**
- Ensure you have access to the repository
- Check if the repository has GitHub Actions enabled

**Logs not displaying properly**
- Ensure GitHub CLI version is 2.75+ (`gh --version`)
- Some private repositories may have restricted log access

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Charm](https://charm.sh/) for the amazing TUI libraries
- [GitHub CLI](https://cli.github.com/) team for the excellent CLI tool
- Inspired by tools like `lazygit` and `gh-dash`