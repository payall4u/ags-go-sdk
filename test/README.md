# Testing Guide

This directory contains integration tests for the Agent Sandbox Go SDK.

## Test Categories

### 1. Public Integration Tests
These tests run against the public API endpoints and can be executed in any environment.

**Files:**
- `code_integration_test.go`
- `command_integration_test.go`
- `filesystem_integration_test.go`
- `domain_integration_test.go`

**Run command:**
```bash
go test ./test/...
# or
make test-integration
```

### 2. Internal Network Tests
These tests require access to internal network endpoints and are only available in internal environments.

**Files:**
- `domain_internal_test.go` (requires `internal` build tag)

**Run command:**
```bash
go test -tags=internal ./test/...
# or
make test-internal
```

## Environment Setup

### Prerequisites
1. Set up cloud credentials (required for all integration tests)
2. For internal tests: Ensure access to internal network endpoints

### Environment Variables
The tests require proper cloud credentials to be configured. Check `helper_test.go` for the specific environment variables needed.

## Build Tags

We use Go build tags to control which tests are compiled and executed:

- **No tags**: Public integration tests only
- **`internal` tag**: Includes internal network tests

### Examples

```bash
# Run only public tests
go test ./test/...

# Run all tests including internal ones
go test -tags=internal ./test/...

# Run tests with verbose output
go test -v -tags=internal ./test/...

# Run specific test
go test -tags=internal -run TestCreate_WithInternalDomain ./test/...
```

## Makefile Commands

```bash
make test              # Run all tests (unit + integration)
make test-unit         # Run unit tests only (short tests)
make test-integration  # Run public integration tests
make test-internal     # Run internal network tests
```