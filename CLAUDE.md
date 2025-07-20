# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Building and Testing

```bash
# Build the Go application
go build

# Build the standalone aruba command
go build ./cmd/aruba

# Run unit tests
go test

# Run integration tests using the Go implementation
go run ./cmd/aruba --strict

# Install the aruba command globally
go install ./cmd/aruba
```

### Ruby-based Testing (Original Aruba)

```bash
# Install Ruby dependencies
bundle install

# Run integration tests using original Ruby Aruba
bundle exec cucumber --publish-quiet --strict
```

### Linting

```bash
# Run golangci-lint
golangci-lint run

# Run spell check
cspell "**/*.{md,go,sh}"
```

## Architecture

### Core Components

**Main Library (`scenario.go`)**

- Implements Godog step definitions for command-line testing
- Provides the `InitializeScenario` function that registers all Aruba steps
- Key functions:
  - `createFile()` - Creates test files with specified content
  - `runCommand()` - Executes shell commands in temporary directories
  - `exitStatus()` - Validates command exit codes
  - `stdout()/stderr()` - Validates command output
  - `fileContains()` - Validates file content

**Command-line Tool (`cmd/aruba/main.go`)**

- Standalone binary that runs Cucumber features using Godog
- Configures Godog with sensible defaults (pretty format, parallel execution)
- Accepts feature file paths as command line arguments

### Testing Framework Integration

**Context Management**

- Uses Go's `context.Context` to pass state between test steps
- Stores temporary directory, exit codes, stdout/stderr in context
- Each scenario gets a fresh temporary directory via `before()` hook

**String Processing**

- `quote()/unquote()` functions handle special character escaping
- Supports newlines (`\n`) and tabs (`\t`) in test strings
- Regex-based pattern matching for flexible string validation

**Step Definitions**

The library implements standard Aruba step patterns:
- `a file named "..." with:` - File creation with content
- `I (successfully) run \`...\`` - Command execution
- `the exit status should (not) be N` - Exit code validation
- `the stdout/stderr should (not) contain (exactly) "..."` - Output validation
- `a file named "..." should (not) contain (exactly) "..."` - File content validation

### Project Structure

**Go Module**

- Module: `github.com/raviqqe/aruba-go`
- Dependencies: Godog for BDD testing, pflag for command-line parsing
- Go version: 1.24.5

**Feature Files**

- Gherkin feature files in `features/` directory
- Test scenarios for command execution, exit codes, output validation
- Compatible with both Go and Ruby Aruba implementations

**Dependents Testing**

- `dependents/` directory contains projects that use aruba-go
- Includes complex projects like `stak` (Scheme interpreter) and `schemat` (JSON schema tool)
- CI tests both Go and Ruby Aruba implementations against dependents

### Development Notes

**Compatibility with Ruby Aruba**

- Maintains compatibility with original Ruby Aruba step definitions
- Test features can run with either Go or Ruby implementation
- Both implementations tested in CI to ensure behavioral consistency

**Command Execution Model**

- Commands run in temporary directories (created per scenario)
- Simple string splitting for command parsing (spaces separate arguments)
- Captures both stdout and stderr for validation
- Exit codes preserved for assertion testing

**Cross-platform Considerations**

- Uses standard Go libraries for file operations and command execution
- Should work on any platform supported by Go and the shell commands being tested