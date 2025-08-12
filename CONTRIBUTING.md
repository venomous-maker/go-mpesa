# Contributing to Go M-Pesa SDK

We welcome contributions to the Go M-Pesa SDK! This document outlines the process for contributing to this project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Code Style](#code-style)
- [Release Process](#release-process)

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct:

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Maintain professionalism in all interactions

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Git
- Basic understanding of M-Pesa API
- Familiarity with Go testing frameworks

### Types of Contributions

We welcome various types of contributions:

- 🐛 **Bug fixes**
- ✨ **New features**
- 📝 **Documentation improvements**
- 🧪 **Tests and test coverage**
- 🔧 **Performance improvements**
- 🎨 **Code style improvements**

## Development Setup

1. **Fork the repository**
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/YOUR_USERNAME/go-mpesa.git
   cd go-mpesa
   ```

2. **Set up the upstream remote**
   ```bash
   git remote add upstream https://github.com/venomous-maker/go-mpesa.git
   ```

3. **Install dependencies**
   ```bash
   go mod download
   ```

4. **Verify the setup**
   ```bash
   go test ./...
   ```

## Making Changes

### Branching Strategy

1. **Create a feature branch**
   ```bash
   git checkout -b feature/your-feature-name
   # or
   git checkout -b fix/issue-description
   ```

2. **Keep your branch updated**
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

### Commit Guidelines

Follow conventional commit format:

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```bash
git commit -m "feat(stk): add transaction status query"
git commit -m "fix(auth): handle token refresh edge case"
git commit -m "docs(readme): update installation instructions"
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestSTKPush ./Tests/

# Run tests with verbose output
go test -v ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Writing Tests

1. **Test file naming**: `*_test.go`
2. **Test function naming**: `TestFunctionName`
3. **Use table-driven tests** for multiple scenarios
4. **Mock external dependencies**

Example test structure:

```go
func TestSTKPush(t *testing.T) {
    tests := []struct {
        name           string
        amount         string
        phoneNumber    string
        expectedError  bool
        expectedResult map[string]any
    }{
        {
            name:          "successful transaction",
            amount:        "100",
            phoneNumber:   "254712345678",
            expectedError: false,
            expectedResult: map[string]any{
                "ResponseCode": "0",
                "CheckoutRequestID": "ws_CO_123456789",
            },
        },
        {
            name:          "invalid phone number",
            amount:        "100",
            phoneNumber:   "invalid",
            expectedError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

### Test Coverage Requirements

- **Minimum coverage**: 80%
- **New features**: Must include tests
- **Bug fixes**: Must include regression tests
- **Critical paths**: 100% coverage required

## Pull Request Process

### Before Submitting

1. **Ensure tests pass**
   ```bash
   go test ./...
   ```

2. **Run linting**
   ```bash
   golangci-lint run
   ```

3. **Format code**
   ```bash
   go fmt ./...
   ```

4. **Update documentation** if needed

### Pull Request Template

Use this template for your PR description:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests pass locally
- [ ] New tests added for new functionality
- [ ] Coverage requirements met

## Checklist
- [ ] Code follows project style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No breaking changes (or clearly documented)

## Related Issues
Fixes #123
```

### Review Process

1. **Automated checks** must pass
2. **At least one approval** from maintainers
3. **All conversations resolved**
4. **Up-to-date with main branch**

## Code Style

### Go Style Guidelines

Follow standard Go conventions:

1. **Formatting**
   ```bash
   go fmt ./...
   goimports -w .
   ```

2. **Naming conventions**
   - Use `camelCase` for unexported functions/variables
   - Use `PascalCase` for exported functions/variables
   - Use descriptive names

3. **Documentation**
   ```go
   // ServiceMethod performs a specific operation.
   // It takes parameter and returns result or error.
   //
   // Example:
   //   result, err := service.ServiceMethod("param")
   //   if err != nil {
   //       log.Fatal(err)
   //   }
   func ServiceMethod(param string) (result string, err error) {
       // implementation
   }
   ```

4. **Error handling**
   ```go
   // Wrap errors with context
   if err != nil {
       return fmt.Errorf("failed to process request: %w", err)
   }
   ```

### Project-Specific Guidelines

1. **Service methods** should return `(map[string]any, error)`
2. **Setter methods** should return the service instance for chaining
3. **Use meaningful variable names** in test files
4. **Add comments for exported functions**

## Release Process

### Versioning

We follow [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. **Update version** in relevant files
2. **Update CHANGELOG.md**
3. **Create release tag**
   ```bash
   git tag -a v1.2.3 -m "Release v1.2.3"
   git push origin v1.2.3
   ```
4. **Create GitHub release** with release notes

## Getting Help

- 📖 **Documentation**: Check existing docs first
- 💬 **Discussions**: Use GitHub Discussions for questions
- 🐛 **Issues**: Report bugs or request features
- 📧 **Email**: contact@venomous-maker.com for sensitive issues

## Recognition

Contributors will be recognized in:

- README.md contributors section
- Release notes
- Special thanks in major releases

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Go M-Pesa SDK! 🚀
