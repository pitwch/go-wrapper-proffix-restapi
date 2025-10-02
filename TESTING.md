# Testing Guide

## Overview

This project has comprehensive test coverage for all major components of the PROFFIX REST API wrapper.

## Running Tests

### Run all tests
```bash
go test -v ./...
```

### Run tests with coverage
```bash
go test -v -race -coverprofile=coverage.txt -covermode=atomic ./...
```

### View coverage report
```bash
go tool cover -html=coverage.txt
```

### Run specific test
```bash
go test -v -run TestClient_Requests ./proffixrest
```

## Test Structure

### Test Files

- **`client_test.go`** - Core client functionality (POST, PUT, GET, DELETE, Login, Logout)
- **`advanced_test.go`** - Advanced features (PATCH, ServiceLogin, concurrent access, options)
- **`batch_test.go`** - Batch request handling
- **`sync_batch_test.go`** - Synchronous batch operations
- **`list_test.go`** - List generation and retrieval
- **`check_test.go`** - API health checks
- **`helper_test.go`** - Helper functions (GetFiltererCount, ReaderToString, etc.)
- **`error_test.go`** - Error handling and PxError types
- **`tools_test.go`** - Utility functions (time conversion, ID extraction)

## Test Coverage

### Core Functions
- ✅ HTTP Methods: GET, POST, PUT, PATCH, DELETE
- ✅ Authentication: Login, Logout, ServiceLogin
- ✅ Session Management: GetPxSessionId, updatePxSessionId
- ✅ Batch Operations: GetBatch, SyncBatch
- ✅ File Operations: File upload/download
- ✅ List Generation: GetList

### Helper Functions
- ✅ GetFiltererCount - Parse metadata from headers
- ✅ ReaderToString/ReaderToByte - Stream conversions
- ✅ GetMaps/GetMap - JSON parsing
- ✅ WriteFile - File writing
- ✅ GetUsedLicences - License calculation
- ✅ GetFileTokens - File token extraction

### Error Handling
- ✅ PxError creation and formatting
- ✅ isInvalidFields / isNotFound type checks
- ✅ Error messages with field details
- ✅ Default error messages for edge cases

### Edge Cases
- ✅ Nil readers
- ✅ Empty headers
- ✅ Invalid URLs
- ✅ Concurrent access
- ✅ Multiple sequential requests
- ✅ Session ID updates

## CI/CD Pipeline

The project uses GitHub Actions for continuous integration with the following jobs:

### Test Job
- **Matrix Testing**: Tests against Go versions 1.20, 1.21, 1.22, 1.23
- **Race Detection**: Runs with `-race` flag to detect data races
- **Coverage**: Generates coverage reports and uploads to Codecov
- **Dependency Verification**: Ensures `go.mod` and `go.sum` are in sync

### Lint Job
- **golangci-lint**: Runs comprehensive linting with 19+ linters
- **Timeout**: 5 minutes to ensure thorough analysis
- **Configuration**: Uses `.golangci.yml` for custom rules

### Build Job
- **Dependencies**: Requires test and lint jobs to pass
- **Builds**: Compiles main package and cmd/proffix-rest
- **Examples**: Builds all example applications

## Linting

### Run linter locally
```bash
golangci-lint run --timeout=5m
```

### Enabled Linters
- errcheck - Check for unchecked errors
- gosimple - Simplify code
- govet - Vet examines Go source code
- ineffassign - Detect ineffectual assignments
- staticcheck - Static analysis
- unused - Check for unused code
- gofmt - Format checking
- goimports - Import formatting
- misspell - Spell checking
- unconvert - Remove unnecessary conversions
- unparam - Detect unused parameters
- bodyclose - Check HTTP response body closes
- gosec - Security issues
- prealloc - Preallocate slices
- gocritic - Comprehensive checks
- revive - Fast linter
- nilness - Detect nil pointer dereferences
- copylocks - Detect mutex copies

## Test Configuration

### Environment Variables
Tests use the PROFFIX public demo server by default:
- URL: `https://portal.proffix.net:11011`
- Database: `DEMODB`
- User: `Gast`

For private testing, set:
- `PXDEMO_URL`
- `PXDEMO_USER`
- `PXDEMO_PASSWORD`
- `PXDEMO_DATABASE`
- `PXDEMO_KEY`

### Test Options
```go
&Options{
    Key:           "...",
    VerifySSL:     false,
    Autologout:    false,
    VolumeLicence: true,
}
```

## Best Practices

1. **Always use context**: Pass `context.Background()` or custom context
2. **Cleanup resources**: Use `defer pxrest.Logout(ctx)` after connection
3. **Check errors**: Always verify error returns
4. **Close readers**: Close `io.ReadCloser` when done
5. **Test isolation**: Each test should be independent

## Troubleshooting

### Tests fail with connection errors
- Check if demo server is accessible
- Verify network connectivity
- Ensure credentials are valid

### Coverage not generated
- Run with `-coverprofile=coverage.txt`
- Ensure all packages are tested with `./...`

### Linter errors
- Run `golangci-lint run` locally first
- Check `.golangci.yml` for configuration
- Some test files exclude certain linters

## Contributing

When adding new features:
1. Write tests first (TDD approach)
2. Ensure tests pass locally
3. Run linter before committing
4. Maintain or improve coverage
5. Add test documentation for complex scenarios
