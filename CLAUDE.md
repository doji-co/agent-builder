# Agent Builder - Development Guide

## Project Overview

Agent Builder is a CLI tool for building ADK (Agent Development Kit) multi-agents. This project is built with:

- **Language**: Go 1.23+
- **CLI Framework**: Cobra (github.com/spf13/cobra)
- **Testing**: Go standard library `testing` package
- **Development Approach**: Test-Driven Development (TDD)

## Test-Driven Development (TDD)

This project follows strict TDD practices. **All features must be developed using the Red-Green-Refactor cycle.**

### The Red-Green-Refactor Cycle

1. **Red**: Write a failing test first
   - Define the expected behavior
   - Run the test and watch it fail
   - Ensure the failure message is clear

2. **Green**: Write minimal code to make the test pass
   - Implement just enough to satisfy the test
   - Don't worry about perfection yet
   - Run the test and watch it pass

3. **Refactor**: Improve the code while keeping tests green
   - Clean up duplication
   - Improve readability
   - Optimize if needed
   - Ensure all tests still pass

### Coverage Target

**Aim for 80%+ test coverage** on all packages. Check coverage with:

```bash
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Code Style

### No Comments Policy

**Do not add comments to code.** Code should be self-documenting through:
- Clear function and variable names
- Small, focused functions with single responsibilities
- Descriptive test names that serve as documentation

Comments are only acceptable for:
- Package-level documentation (package comments)
- Complex algorithms requiring mathematical or academic explanation
- Workarounds for known issues/bugs in dependencies

### Example

Bad:
```go
// Create a new agent with the given parameters
func NewAgent(name string, agentType AgentType) *Agent {
    // Initialize the agent struct
    return &Agent{
        Name: name,
        Type: agentType,
    }
}
```

Good:
```go
func NewAgent(name string, agentType AgentType) *Agent {
    return &Agent{
        Name: name,
        Type: agentType,
    }
}
```

## Testing Standards

### Test File Organization

```
project/
├── cmd/
│   └── root.go
│   └── root_test.go          # Test files alongside source
├── internal/
│   ├── agent/
│   │   ├── builder.go
│   │   └── builder_test.go
│   └── config/
│       ├── config.go
│       └── config_test.go
└── pkg/
    └── ...
```

### Test Naming Conventions

- Test files: `*_test.go`
- Test functions: `TestFunctionName(t *testing.T)`
- Subtests: Use `t.Run("description", func(t *testing.T) {...})`
- Table-driven tests: Use descriptive test case names

### Example Test Structure

```go
func TestBuilder_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   BuilderConfig
		want    *Agent
		wantErr bool
	}{
		{
			name: "valid configuration creates agent",
			input: BuilderConfig{
				Name: "test-agent",
				Type: "multi",
			},
			want: &Agent{
				Name: "test-agent",
				Type: "multi",
			},
			wantErr: false,
		},
		{
			name: "empty name returns error",
			input: BuilderConfig{
				Name: "",
				Type: "multi",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			builder := NewBuilder()
			got, err := builder.Create(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
```

## Development Workflow

### Before Writing Any Feature

1. **Understand the requirement** - Know what behavior you're implementing
2. **Write the test first** - Define expected inputs and outputs
3. **Run tests** - Verify the test fails for the right reason
4. **Implement** - Write minimal code to pass the test
5. **Refactor** - Clean up while keeping tests green
6. **Commit** - Commit both test and implementation together

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests in a specific package
go test ./internal/agent

# Run a specific test
go test -run TestBuilder_Create ./internal/agent

# Run tests with coverage
go test -cover ./...

# Watch mode (use with tools like gotest or gotestsum)
gotestsum --watch
```

## Code Organization

### Package Structure

```
agent-builder/
├── cmd/                    # CLI commands (Cobra commands)
│   ├── root.go            # Root command
│   ├── create.go          # Subcommands
│   └── *_test.go          # Command tests
├── internal/              # Private application code
│   ├── agent/             # Core agent building logic
│   ├── config/            # Configuration handling
│   ├── template/          # Template processing
│   └── validation/        # Input validation
├── pkg/                   # Public libraries (if any)
└── testdata/              # Test fixtures and data
```

### Testability Principles

1. **Dependency Injection**: Pass dependencies as parameters or interfaces
2. **Interfaces over Concrete Types**: Use interfaces for testability
3. **Small Functions**: Easier to test and reason about
4. **Pure Functions**: Prefer functions with no side effects
5. **Separation of Concerns**: Business logic separate from I/O

### Example: Testable Design

**Bad** (hard to test):
```go
func ProcessAgent() error {
	config := loadConfigFromFile("/etc/config.yaml") // File I/O
	agent := buildAgent(config)                       // Hard to mock
	return saveToDatabase(agent)                      // Database I/O
}
```

**Good** (easy to test):
```go
type ConfigLoader interface {
	Load() (Config, error)
}

type AgentRepository interface {
	Save(agent Agent) error
}

func ProcessAgent(loader ConfigLoader, repo AgentRepository) error {
	config, err := loader.Load()
	if err != nil {
		return err
	}

	agent := buildAgent(config)
	return repo.Save(agent)
}
```

## CLI-Specific Testing with Cobra

### Testing Commands

```go
func TestCreateCommand(t *testing.T) {
	// Setup
	cmd := NewCreateCommand()
	buf := new(bytes.Buffer)
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	// Set arguments
	cmd.SetArgs([]string{"--name", "test-agent"})

	// Execute
	err := cmd.Execute()

	// Assert
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Agent created successfully") {
		t.Errorf("expected success message, got: %s", output)
	}
}
```

### Testing Flags and Arguments

```go
func TestCreateCommand_Flags(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "missing required name flag",
			args:    []string{},
			wantErr: true,
			errMsg:  "required flag \"name\" not set",
		},
		{
			name:    "valid flags",
			args:    []string{"--name", "agent1", "--type", "multi"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewCreateCommand()
			cmd.SetArgs(tt.args)

			err := cmd.Execute()

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr && !strings.Contains(err.Error(), tt.errMsg) {
				t.Errorf("error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}
```

## Test Types

### Unit Tests
- Test individual functions and methods
- Mock external dependencies
- Fast execution (<1s for all unit tests)
- Located alongside source files

### Integration Tests
- Test multiple components working together
- May use real implementations (file system, etc.)
- Place in `*_integration_test.go` or separate test package
- Use build tags if needed: `// +build integration`

### End-to-End Tests
- Test complete user workflows
- Execute actual CLI commands
- Verify output and side effects
- Can be slower but provide high confidence

## Best Practices

### What to Test

✅ **Do test:**
- Public API functions and methods
- Business logic and algorithms
- Edge cases and error conditions
- Input validation
- Command flag parsing
- Output formatting

❌ **Skip testing:**
- Trivial getters/setters with no logic
- Third-party library code
- Generated code

### Mocking External Dependencies

```go
// Define interface for external dependency
type FileSystem interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
}

// Mock implementation for testing
type MockFileSystem struct {
	ReadFileFunc  func(path string) ([]byte, error)
	WriteFileFunc func(path string, data []byte) error
}

func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
	if m.ReadFileFunc != nil {
		return m.ReadFileFunc(path)
	}
	return nil, nil
}

func (m *MockFileSystem) WriteFile(path string, data []byte) error {
	if m.WriteFileFunc != nil {
		return m.WriteFileFunc(path, data)
	}
	return nil
}
```

### Test Isolation

- Each test should be independent
- Use `t.Cleanup()` for teardown
- Don't rely on test execution order
- Create fresh test fixtures per test

```go
func TestSomething(t *testing.T) {
	// Setup
	tmpDir := t.TempDir() // Automatically cleaned up

	// Cleanup other resources
	t.Cleanup(func() {
		// cleanup code
	})

	// Test code...
}
```

### Table-Driven Tests

Use table-driven tests for testing multiple scenarios:

```go
func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid input", "valid-name", false},
		{"empty string", "", true},
		{"too long", strings.Repeat("a", 256), true},
		{"invalid chars", "name@#$", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
```

## CI/CD Integration

Tests will automatically run on:
- Every commit (via pre-commit hooks)
- Pull requests (via GitHub Actions)
- Before releases (via GitHub Actions)

The build will fail if:
- Any test fails
- Coverage drops below 80%
- `go vet` reports issues
- `go fmt` shows unformatted code

## Quick Reference

```bash
# Run tests
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...

# Run specific test
go test -run TestName ./path/to/package

# Coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Check formatting
go fmt ./...

# Lint
go vet ./...

# Build
go build -o agent-builder .

# Install Cobra CLI (helpful for generating commands)
go install github.com/spf13/cobra-cli@latest
```

## Adding New Features

### Step-by-Step Example

Let's say you want to add a new command `agent-builder validate`:

1. **Write test first** (`cmd/validate_test.go`):
```go
func TestValidateCommand(t *testing.T) {
	cmd := NewValidateCommand()
	cmd.SetArgs([]string{"agent.yaml"})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
```

2. **Run test** - it should fail:
```bash
go test ./cmd
```

3. **Implement command** (`cmd/validate.go`):
```go
func NewValidateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "validate [file]",
		Short: "Validate agent configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Implementation
			return nil
		},
	}
}
```

4. **Run test again** - it should pass:
```bash
go test ./cmd
```

5. **Refactor** if needed, keeping tests green

6. **Commit** both test and implementation

## Resources

- [Go Testing Package](https://pkg.go.dev/testing)
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Table Driven Tests](https://dave.cheney.net/2019/05/07/prefer-table-driven-tests)
- [Test Fixtures in Go](https://dave.cheney.net/2016/05/10/test-fixtures-in-go)

## Summary

Remember: **Test First, Code Second**. Every feature starts with a failing test. This ensures:
- Clear requirements
- Testable design
- High confidence in changes
- Living documentation
- Easier refactoring

Happy coding! 🚀
