# Contributing to Recontronic CLI Client

Thank you for your interest in contributing to the Recontronic CLI Client! This document provides guidelines and instructions for contributing.

## Development Environment Setup

### Prerequisites

- Go 1.21 or higher
- Git
- Make (optional, but recommended)
- golangci-lint (for code quality)

### Initial Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/recontronic-cli-client.git
cd recontronic-cli-client

# Initialize Go modules
go mod download

# Install development tools
make install-tools

# Run tests to verify setup
make test
```

## Code Standards

### Go Style Guide

- Follow the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- Use `gofmt` for formatting (enforced by linter)
- Keep functions small and focused
- Write clear, descriptive variable names
- Add comments for exported functions and types

### Code Quality Requirements

All code must:
- Pass `golangci-lint` without errors
- Have at least 80% test coverage for new code
- Include unit tests for business logic
- Use proper error handling (no naked returns or ignored errors)
- Follow Go conventions (error as last return value, etc.)

### Linting

```bash
# Run linter
make lint

# Auto-fix issues where possible
make lint-fix
```

## Project Structure

```
recontronic-cli-client/
├── cmd/                    # Command implementations (Cobra commands)
├── pkg/                    # Reusable packages
│   ├── client/            # API clients (REST, gRPC)
│   ├── config/            # Configuration management
│   ├── ui/                # UI components (TUI, formatters)
│   └── models/            # Data models and types
├── proto/                 # Protocol buffer definitions
├── scripts/               # Build and utility scripts
├── docs/                  # Documentation
├── test/                  # Integration tests
└── main.go               # Application entry point
```

### Package Guidelines

- **cmd/**: Keep thin, delegate to pkg/ for logic
- **pkg/client/**: All external API communication
- **pkg/config/**: Configuration loading and validation
- **pkg/ui/**: All user interface code (formatters, TUI, etc.)
- **pkg/models/**: Data structures shared across packages

## Development Workflow

### 1. Create a Branch

```bash
# Create feature branch
git checkout -b feature/your-feature-name

# Create bugfix branch
git checkout -b fix/bug-description
```

### 2. Make Changes

- Write code following the style guide
- Add tests for new functionality
- Update documentation as needed
- Keep commits focused and atomic

### 3. Test Your Changes

```bash
# Run unit tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Run integration tests (if applicable)
make test-integration
```

### 4. Commit Changes

Follow conventional commit format:

```bash
# Format: <type>(<scope>): <subject>
git commit -m "feat(program): add delete program command"
git commit -m "fix(config): handle missing config file gracefully"
git commit -m "docs(readme): update installation instructions"
```

**Commit Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

### 5. Push and Create Pull Request

```bash
# Push branch
git push origin feature/your-feature-name

# Create PR via GitHub UI
```

## Pull Request Guidelines

### PR Title Format

Use conventional commit format:
- `feat: Add program delete functionality`
- `fix: Handle nil pointer in scan watch`
- `docs: Update contributing guidelines`

### PR Description Template

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix (non-breaking change fixing an issue)
- [ ] New feature (non-breaking change adding functionality)
- [ ] Breaking change (fix or feature causing existing functionality to break)
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass (if applicable)
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Comments added for complex code
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests added/updated
```

### Review Process

1. All PRs require at least one approval
2. All CI checks must pass
3. Code coverage must not decrease
4. No linting errors

## Testing Guidelines

### Unit Tests

```go
// Example unit test
func TestProgramAdd(t *testing.T) {
    tests := []struct {
        name    string
        input   ProgramInput
        want    *Program
        wantErr bool
    }{
        {
            name: "valid program",
            input: ProgramInput{
                Name:     "Test Corp",
                Platform: "hackerone",
                Scope:    []string{"*.example.com"},
            },
            want: &Program{
                Name:     "Test Corp",
                Platform: "hackerone",
                Scope:    []string{"*.example.com"},
            },
            wantErr: false,
        },
        // Add more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := AddProgram(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("AddProgram() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Add assertions
        })
    }
}
```

### Table-Driven Tests

Prefer table-driven tests for comprehensive coverage:

```go
func TestFormatOutput(t *testing.T) {
    tests := []struct {
        name   string
        format string
        data   interface{}
        want   string
    }{
        // Test cases
    }
    // Test execution
}
```

### Integration Tests

```go
// +build integration

func TestProgramAPIIntegration(t *testing.T) {
    // Requires RECON_SERVER environment variable
    if os.Getenv("RECON_SERVER") == "" {
        t.Skip("Integration tests require RECON_SERVER")
    }
    // Test implementation
}
```

## Documentation

### Code Documentation

```go
// AddProgram creates a new bug bounty program in the platform.
// It validates the input, sends a REST API request, and returns
// the created program details.
//
// Example:
//   prog, err := AddProgram(ProgramInput{
//       Name: "Example Corp",
//       Scope: []string{"*.example.com"},
//   })
func AddProgram(input ProgramInput) (*Program, error) {
    // Implementation
}
```

### Update Documentation

When adding features, update:
- README.md (if user-facing)
- Package documentation (godoc comments)
- Example usage in docs/
- Help text in Cobra commands

## Common Tasks

### Adding a New Command

1. Create command file in `cmd/`
2. Implement using Cobra framework
3. Add to root command
4. Add tests
5. Update README.md

### Adding REST API Endpoint

1. Add method to `pkg/client/rest.go`
2. Add corresponding model to `pkg/models/`
3. Add tests
4. Update documentation

### Adding gRPC Stream

1. Update proto definitions in `proto/`
2. Regenerate Go code: `make proto`
3. Implement client in `pkg/client/grpc.go`
4. Add tests
5. Update documentation

## Debugging

### Enable Debug Logging

```bash
export RECON_LOG_LEVEL=debug
recon-cli <command>
```

### Using Delve Debugger

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug a test
dlv test ./cmd -- -test.run TestProgramAdd

# Debug the CLI
dlv debug . -- program list
```

## Release Process

Releases are managed by maintainers:

1. Update version in `version.go`
2. Update CHANGELOG.md
3. Tag release: `git tag -a v1.0.0 -m "Release v1.0.0"`
4. Push tag: `git push origin v1.0.0`
5. GitHub Actions builds and publishes binaries

## Getting Help

- Review existing issues and PRs
- Check documentation in docs/
- Ask questions in GitHub Discussions
- Join Discord (link in README)

## Code of Conduct

### Our Standards

- Be respectful and inclusive
- Welcome newcomers
- Accept constructive criticism
- Focus on what's best for the community
- Show empathy towards others

### Unacceptable Behavior

- Harassment or discriminatory language
- Trolling or personal attacks
- Publishing others' private information
- Unprofessional conduct

## Attribution

By contributing, you agree that your contributions will be licensed under the MIT License.

Thank you for contributing! 🎉
