# WebP Converter - Linting Guide

## GitHub Actions

This repository includes a GitHub Actions workflow that automatically runs linting on:
- Push to `main`, `master`, or `develop` branches
- Pull requests to `main`, `master`, or `develop` branches

The workflow includes:
- `go vet` - Go's built-in static analysis tool
- `golangci-lint` - Comprehensive Go linter with multiple analyzers
- Build verification
- Test execution (if tests exist)

## Local Development

### Install golangci-lint

```bash
# On macOS
brew install golangci-lint

# On Linux (install latest version compatible with Go 1.24.1)
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest

# Or using Go install (latest version)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Run linting locally

```bash
# Run all linters
golangci-lint run

# Run specific linters
golangci-lint run --enable-only=errcheck,gofmt,govet

# Run with fix (auto-fix some issues)
golangci-lint run --fix

# Run go vet
go vet ./...

# Format code
go fmt ./...
```

### Configuration

The linting configuration is in `.golangci.yml` and includes:
- Essential linters: `errcheck`, `gofmt`, `goimports`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `typecheck`, `unused`
- Code quality linters: `misspell`, `unconvert`, `gocyclo`, `funlen`
- Deprecated linters removed: `deadcode`, `varcheck`, `structcheck` (replaced by `unused`)
- Relaxed rules for test files

### Pre-commit Hook (Optional)

You can add a pre-commit hook to run linting automatically:

```bash
# Create .git/hooks/pre-commit
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/sh
golangci-lint run
if [ $? -ne 0 ]; then
    echo "Linting failed. Please fix the issues before committing."
    exit 1
fi
EOF

chmod +x .git/hooks/pre-commit
```
