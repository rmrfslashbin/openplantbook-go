# Contributing to OpenPlantbook Go SDK

First off, thank you for considering contributing to the OpenPlantbook Go SDK! It's people like you that make this library better for everyone.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include as many details as possible:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples** (code snippets, minimal reproducible examples)
- **Describe the behavior you observed** and what you expected
- **Include your environment details** (Go version, OS, SDK version)

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide a detailed description** of the suggested enhancement
- **Explain why this enhancement would be useful** to most users
- **List some examples** of how it would be used

### Pull Requests

1. **Fork the repo** and create your branch from `main`
2. **Add tests** for any new functionality
3. **Ensure tests pass**: `make test`
4. **Run quality checks**: `make quality`
5. **Update documentation** if you're changing functionality
6. **Follow the Go coding style** (use `gofmt`, `golint`)
7. **Write clear commit messages** (see below)

## Development Setup

### Prerequisites

- Go 1.23 or higher
- Git
- Make

### Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/openplantbook-go.git
cd openplantbook-go

# Install dependencies
go mod download

# Run tests
make test

# Run quality checks
make quality
```

### Testing

We maintain high test coverage (90%+). Please add tests for new functionality:

```bash
# Run tests with coverage
make test

# Run tests with race detector
make test

# View coverage report
make coverage

# Run specific tests
go test -run TestName -v
```

### Code Quality

Before submitting a PR, ensure your code passes all quality checks:

```bash
# Run all quality checks
make quality

# Individual checks
make vet          # Go vet
make fmt          # Go fmt
make staticcheck  # Staticcheck linter
make deadcode     # Dead code analysis
```

## Coding Standards

### Go Style Guide

- Follow the [Effective Go](https://golang.org/doc/effective_go.html) guidelines
- Use `gofmt` to format your code
- Write idiomatic Go code
- Keep functions small and focused
- Use meaningful variable and function names

### Documentation

- Add godoc comments for all exported functions, types, and constants
- Include usage examples in godoc comments
- Update README.md if adding new features
- Update CHANGELOG.md following [Keep a Changelog](https://keepachangelog.com/)

### Testing

- Write table-driven tests when appropriate
- Test both success and error cases
- Use descriptive test names: `TestFunctionName_Scenario`
- Avoid external dependencies in unit tests (use mocks)
- Keep test coverage above 80%

### Commit Messages

Follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>: <description>

[optional body]

[optional footer]
```

Types:
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `test`: Adding or updating tests
- `refactor`: Code refactoring
- `perf`: Performance improvements
- `chore`: Maintenance tasks

Examples:
```
feat: add support for plant sensor data endpoints

fix: correct authentication header format for OAuth2

docs: update README with new authentication examples

test: improve coverage for error handling
```

## Project Structure

```
.
├── cache.go              # Caching implementation
├── client.go             # Client initialization
├── errors.go             # Error types
├── models.go             # API response models
├── options.go            # Client options
├── plants.go             # Plant API endpoints
├── cmd/
│   └── openplantbook/    # CLI tool
├── examples/             # Usage examples
└── testdata/             # Test fixtures
```

## Release Process

Releases are managed by maintainers:

1. Update CHANGELOG.md
2. Update version in relevant files
3. Create and push git tag: `git tag -a v1.0.0 -m "Release v1.0.0"`
4. GitHub Actions will build and publish the release

## Questions?

- Check existing [issues](https://github.com/rmrfslashbin/openplantbook-go/issues)
- Open a new issue for questions
- Review the [README](README.md) and [documentation](https://pkg.go.dev/github.com/rmrfslashbin/openplantbook-go)

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing! 🌱
