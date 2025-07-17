# Contributing to GitHub Actions TUI

Thank you for your interest in contributing to GitHub Actions TUI! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.21 or later
- [GitHub CLI](https://cli.github.com/) v2.75+
- Git

### Setting up the development environment

1. Fork the repository
2. Clone your fork:
   ```bash
   git clone https://github.com/yourusername/gh-actions-tui.git
   cd gh-actions-tui
   ```
3. Install dependencies:
   ```bash
   go mod download
   ```
4. Build the project:
   ```bash
   go build -o gh-actions-tui .
   ```

## Development Workflow

### Running the application

```bash
go run main.go -r owner/repo
```

### Running tests

```bash
go test ./...
```

### Running with race detection

```bash
go test -race ./...
```

### Code formatting

```bash
go fmt ./...
```

### Linting

```bash
go vet ./...
```

## Project Structure

```
gh-actions-tui/
├── main.go              # Application entry point
├── tui/                 # TUI components
│   ├── model.go         # Bubble Tea model and state management
│   ├── view.go          # UI rendering and layout
│   └── keys.go          # Key bindings and help text
├── github/              # GitHub API integration
│   └── client.go        # GitHub CLI wrapper
├── .github/workflows/   # CI/CD workflows
├── README.md
└── CONTRIBUTING.md
```

## Code Style

- Follow standard Go conventions
- Use `gofmt` for formatting
- Write clear, descriptive commit messages
- Add comments for complex logic
- Keep functions focused and small

## Making Changes

### Before you start

1. Check existing issues and pull requests
2. Create an issue to discuss major changes
3. Fork the repository and create a feature branch

### Development process

1. Create a feature branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes following the code style guidelines

3. Add or update tests as needed

4. Ensure all tests pass:
   ```bash
   go test ./...
   ```

5. Commit your changes:
   ```bash
   git commit -m "feat: add your feature description"
   ```

6. Push to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

7. Create a pull request

### Commit Message Format

We follow conventional commits format:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

Examples:
- `feat: add search functionality to log viewer`
- `fix: resolve viewport scrolling issue`
- `docs: update installation instructions`

## Testing

### Writing tests

- Write unit tests for new functionality
- Ensure tests are deterministic and fast
- Use table-driven tests where appropriate
- Mock external dependencies

### Test coverage

- Aim for good test coverage on new code
- Run tests with coverage: `go test -cover ./...`

## UI/UX Guidelines

### TUI Design Principles

- Keep the interface intuitive and consistent
- Provide clear visual feedback for user actions
- Show helpful key bindings and status information
- Handle errors gracefully with clear messages
- Maintain responsive performance

### Key Bindings

- Follow common conventions (vim-like navigation)
- Provide consistent bindings across views
- Show available keys in help text
- Handle edge cases (empty lists, loading states)

## Submitting Pull Requests

### Pull Request Checklist

- [ ] Code follows project style guidelines
- [ ] Tests pass locally
- [ ] New functionality includes tests
- [ ] Documentation is updated if needed
- [ ] Commit messages follow conventional format
- [ ] PR description explains the changes

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass
- [ ] Manual testing completed

## Screenshots (if applicable)
Add screenshots for UI changes
```

## Reporting Issues

### Bug Reports

Include:
- Go version (`go version`)
- GitHub CLI version (`gh --version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Error messages or logs

### Feature Requests

Include:
- Clear description of the feature
- Use case and motivation
- Possible implementation approach
- Examples from similar tools

## Getting Help

- Check existing issues and discussions
- Ask questions in issue comments
- Reach out to maintainers for guidance

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help create a welcoming environment
- Follow GitHub's community guidelines

Thank you for contributing to GitHub Actions TUI! 🎉