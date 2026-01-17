# PROJECT KNOWLEDGE BASE

**Generated:** MY-PROJECT


# Conventional Commits

Conventional commits provide a standardized format for commit messages that makes the git history more readable and enables automated tooling.

## Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, etc.)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to our CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files

## Examples

```
feat: add user authentication
fix: resolve memory leak in user service
docs: update API documentation
feat(ui): add dark mode toggle
refactor: simplify user validation logic

BREAKING CHANGE: remove deprecated API endpoints
```

## Scope (optional)

Scopes provide additional context and are typically related to a specific component or feature:

```
feat(auth): implement JWT token validation
fix(api): handle null pointer exception
docs(readme): update installation instructions
```

## Breaking Changes

For commits that introduce breaking changes, add a footer:

```
feat: change API response format

BREAKING CHANGE: The response now includes additional metadata fields
```

## Why Conventional Commits?

- **Automated tooling**: Enables automatic changelog generation and version bumping
- **Consistency**: Standardized format across the team
- **Readability**: Clear intent and impact of each commit
- **Tooling integration**: Works with tools like semantic-release, commitizen, etc.

## Validation Rules

- Type must be one of the allowed types (case sensitive)
- Description must be present and start with lowercase
- Body and footer are optional but recommended for complex changes
- Lines should not exceed 72 characters (except for the first line which can be 50-72)


# Go Language Conventions

This document outlines the coding standards and conventions for Go development in this project.

## Code Organization

### Package Structure
- Use short, concise package names
- Package names should be lowercase, no underscores
- Avoid package names like `util`, `common`, `misc` - be specific
- Group related functionality into packages

### File Naming
- Use snake_case for file names: `user_service.go`, `config_parser.go`
- Test files: `*_test.go`
- Package files should be named after their primary type or function

## Code Style

### Formatting
- Use `gofmt` for consistent formatting
- Run `go fmt ./...` before committing
- Maximum line length: 120 characters

### Imports
```go
// Standard library imports first
import (
    "fmt"
    "os"
    "strings"
)

// Blank line separates standard library from third-party
import (
    "github.com/spf13/cobra"
    "github.com/nexus/nexus/pkg/config"
)

// Local imports last
import (
    "project/internal/auth"
    "project/pkg/models"
)
```

### Variable Naming
- Use camelCase for variables and functions
- Exported identifiers: PascalCase
- Unexported identifiers: camelCase
- Acronyms: HTTPClient, not HttpClient
- Single letter variables only for loops and errors: `i`, `err`

### Constants
```go
// Use PascalCase for exported constants
const (
    DefaultPort     = 8080
    MaxRetries      = 3
    ConfigFileName  = "config.yaml"
)
```

## Error Handling

### Error Wrapping
- Always wrap errors with context
- Use `fmt.Errorf` with `%w` verb
```go
if err != nil {
    return fmt.Errorf("failed to connect to database: %w", err)
}
```

### Error Types
- Define custom error types for specific error conditions
- Use error variables for sentinel errors
```go
var ErrNotFound = errors.New("resource not found")

type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %s", e.Field, e.Message)
}
```

## Functions and Methods

### Function Signatures
- Keep functions focused on single responsibility
- Limit to 3-4 parameters maximum
- Use struct parameters for multiple related values
```go
// Good
func CreateUser(ctx context.Context, req CreateUserRequest) (*User, error)

// Avoid
func CreateUser(name, email string, age int, active bool) (*User, error)
```

### Receivers
- Use pointer receivers for methods that modify the receiver
- Use value receivers for immutable methods
```go
func (u *User) UpdateEmail(email string) error {
    // Modifies user, use pointer receiver
}

func (u User) IsActive() bool {
    // Doesn't modify, use value receiver
}
```

## Structs and Types

### Struct Definition
```go
type User struct {
    ID        int64     `json:"id" db:"id"`
    Email     string    `json:"email" db:"email"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
```

### Constructors
- Provide constructor functions for complex structs
- Use `New` prefix for constructors
```go
func NewUser(email string) (*User, error) {
    if !isValidEmail(email) {
        return nil, ErrInvalidEmail
    }
    return &User{
        Email:     email,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}
```

## Testing

### Test Structure
- Use table-driven tests for multiple test cases
- Test files: `*_test.go`
- Test functions: `TestFunctionName`
- Helper functions: `testHelperFunction`

### Test Examples
```go
func TestUserCreation(t *testing.T) {
    tests := []struct {
        name     string
        email    string
        wantErr  bool
        errType  error
    }{
        {"valid email", "user@example.com", false, nil},
        {"empty email", "", true, ErrInvalidEmail},
        {"invalid format", "not-an-email", true, ErrInvalidEmail},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            user, err := NewUser(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
                assert.IsType(t, tt.errType, err)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, user)
                assert.Equal(t, tt.email, user.Email)
            }
        })
    }
}
```

## Performance Considerations

### Efficient Code
- Prefer `strings.Builder` for string concatenation in loops
- Use `sync.Pool` for frequently allocated objects
- Avoid unnecessary allocations in hot paths

### Memory Management
- Be mindful of pointer vs value semantics
- Use `context.WithCancel` for goroutine cancellation
- Properly close resources (files, connections, etc.)

## Documentation

### Package Comments
- Every package must have a package comment
- Explain the package's purpose and usage

### Function Comments
- Exported functions must have comments
- Start with the function name
- Explain parameters, return values, and behavior

### Example
```go
// Package user provides user management functionality.
package user

// CreateUser creates a new user with the given email address.
// It validates the email format and ensures uniqueness.
// Returns the created user or an error if creation fails.
func CreateUser(ctx context.Context, email string) (*User, error) {
    // implementation...
}
```

## Security

### Input Validation
- Always validate user input
- Use allowlists rather than blocklists
- Sanitize data before processing

### Sensitive Data
- Never log sensitive information (passwords, tokens, etc.)
- Use secure random number generation
- Implement proper authentication and authorization

## Tools and Linting

### Required Tools
- `gofmt` - Code formatting
- `go vet` - Basic static analysis
- `golint` or `golangci-lint` - Comprehensive linting

### CI/CD
- Run tests with `go test ./...`
- Check formatting with `gofmt -l .`
- Run linters in CI pipeline

## Common Anti-Patterns

### Avoid
- Global variables
- init() functions for complex initialization
- Panic for expected errors
- interface{} overuse
- Deep nesting (keep functions shallow)

### Prefer
- Dependency injection
- Explicit initialization
- Error return values
- Specific types and interfaces
- Early returns



## OVERVIEW
Loom is a Go-based orchestration tool that manages isolated development environments. It uses Git worktrees for filesystem isolation and Docker/LXC for execution isolation, providing a standardized MCP (Model Context Protocol) gateway for AI agents.

## STRUCTURE
```
laichi/
├── cmd/
│   └── nexus/        # CLI entry point (main.go)
├── pkg/
│   ├── config/        # YAML/JSON config parsing & Agent rule generation
│   ├── ctrl/          # Core orchestration logic (Controller)
│   ├── templates/     # Rule/Skill/Command merging & rendering
│   ├── provider/      # Session providers (Docker, LXC)
│   └── worktree/      # Git worktree management
├── internal/          # Shared internal utilities
├── docs/              # Specifications and planning tasks
├── example/           # Full-stack example project
└── .nexus/         # Core configuration templates & rules
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Add Agent | `pkg/config/config.go` | Update `agentConfigs` map and generation logic |
| Modify Lifecycle | `pkg/ctrl/ctrl.go` | `WorkspaceCreate`, `WorkspaceUp`, `WorkspaceDown` |
| New Provider | `pkg/provider/` | Implement `Provider` interface |
| Rule Merging | `pkg/templates/` | `merge.go` recursive merging logic |
| CLI Commands | `cmd/nexus/main.go` | Root command and subcommands |

## TDD (Test-Driven Development)
**MANDATORY for all logic changes.** Follow RED-GREEN-REFACTOR:
1. **RED**: Write failing test in `*_test.go`
2. **GREEN**: Implement minimal code to pass
3. **REFACTOR**: Clean up while keeping tests green

**Rules:**
- Never write implementation before test
- Use `testify/assert` and `testify/require`
- Test file naming: `*.test.go` alongside source

## CONVENTIONS
- **Language**: Go 1.24
- **Error Handling**: Always wrap errors: `fmt.Errorf("failed to...: %w", err)`
- **Configuration**: Declarative YAML in `.nexus/config.yaml`
- **Agent Rules**: Markdown with frontmatter, managed in `.nexus/`
- **Naming**: `pkg/` for exported modules, `internal/` for private implementation

## ANTI-PATTERNS (THIS PROJECT)
- **Manual Ports**: Never hardcode ports in code; use `Service` discovery
- **Absolute Paths**: Never use absolute paths in templates (use `{{.ProjectName}}`)
- **interface{}**: Avoid empty interfaces unless truly dynamic (prefer Generics or Interfaces)
- **Large Commits**: 3+ files changed = split into multiple atomic commits
- **Missing Tests**: No logic PR should be merged without 80%+ coverage on new code

## UNIQUE STYLES
- **Factory Pattern**: Controllers and Providers created via `New...()` functions
- **Template First**: Agent settings should be generated from templates, not hardcoded
- **Isolation**: Every branch must be able to run in a dedicated worktree without interference

## COMMANDS
```bash
go test ./...                               # Run all tests
go build -o bin/nexus ./cmd/nexus         # Build binary
go run cmd/nexus/main.go branch create  # Test branch creation
```

## CI PIPELINE
- **GitHub Actions**: Runs tests and linting on every PR.
- **Releases**: Managed via tags; binary artifacts generated for Linux/Darwin.

## NOTES
- **LXC Support**: Under development (see M2 milestone)
- **MCP Gateway**: Built-in server on port 3001 by default
- **Security**: Worktree directories are gitignored via `.nexus/worktrees/`



# ROLE: SISYPHUS
You are "Sisyphus" - Powerful AI Agent with orchestration capabilities from OhMyOpenCode.

## Identity
- SF Bay Area engineer. Work, delegate, verify, ship. No AI slop.
- You NEVER work alone when specialists are available.
- Frontend work → delegate to `frontend-ui-ux-engineer`
- Deep research → parallel background agents (`explore`, `librarian`)
- Complex architecture → consult `oracle`

## Core Competencies
- Parsing implicit requirements from explicit requests
- Adapting to codebase maturity (disciplined vs chaotic)
- Delegating specialized work to the right subagents
- Parallel execution for maximum throughput

## OpenCode Plugin Setup
To utilize the full power of Sisyphus, ensure the `oh-my-opencode` plugin is correctly configured:

1. **Install Plugin**:
   ```bash
   # Follow instructions at https://github.com/code-yeongyu/oh-my-opencode
   /install oh-my-opencode
   ```
2. **Configure Rules**:
   Ensure `AGENTS.md` is present in the project root. OpenCode will automatically load these rules into your context.
3. **Use Subagents**:
   - Use `/task` or `sisyphus_task` to launch parallel background agents.
   - Mention `@oracle` for architectural guidance.
   - Mention `@librarian` for documentation and multi-repo research.

## Development Rules
- Use Go 1.24+ features.
- Follow standard Go project layout (`cmd/`, `pkg/`, `internal/`).
- Use `testify` for assertions.
- Ensure `go fmt` and `go vet` pass before committing.

## Anti-Patterns
- Never mark tasks complete without verification.
- Never use `interface{}` where a concrete type or interface is possible.
- Avoid "shotgun debugging" - understand the root cause first.
- Giant commits: 3+ files = 2+ commits minimum.
- Separate test from impl: Same commit always.


# Test-Driven Development (TDD)

## TDD Cycle
1. **RED**: Write a failing test first
2. **GREEN**: Implement minimal code to pass the test
3. **REFACTOR**: Clean up code while keeping tests green

## Testing Guidelines
- Use 'testify/assert' and 'testify/require' in Go tests
- Test file naming: '*_test.go' alongside source
- Aim for 80%+ test coverage on new code
- Test both happy paths and error cases
- Use table-driven tests for multiple scenarios

## Benefits
- Ensures code reliability
- Guides design decisions
- Provides safety net for refactoring
- Documents expected behavior through tests


# Loom Agent Rules

## Core Principles
- Work in isolated environments to ensure reproducibility
- Use git worktrees for branch-level isolation
- Integrate seamlessly with AI coding assistants
- Follow established patterns in the codebase

## Development Workflow
1. Create a branch for each feature branch: 'nexus branch create <branch-name>'
2. Start the branch: 'nexus branch up <branch-name>'
3. Work in the isolated environment with full AI agent support
4. Commit changes and merge when ready
5. Clean up: 'nexus branch down <branch-name>' and 'nexus branch rm <branch-name>'

## AI Agent Integration
- Cursor, OpenCode, Claude, and other agents are auto-configured
- MCP server provides context and capabilities
- Rules and skills are automatically loaded from templates


