# Changelog

All notable changes to the Go M-Pesa SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Planned
- Advanced retry mechanisms with exponential backoff
- Request/response logging middleware
- Performance optimizations
- Additional webhook signature verification

## [1.0.1] - 2024-08-12

### Added
- Comprehensive documentation suite with API reference
- Complete examples directory with practical implementations
- STK Push example with web server and callback handling
- B2C payment example with result and timeout callbacks
- Contributing guidelines with development setup instructions
- MIT License for open source distribution
- Professional README with table of contents and badges
- API documentation with complete method signatures
- Environment configuration templates (.env.example files)
- Webhook handling patterns and best practices
- Error handling examples and troubleshooting guides

### Enhanced
- README.md with comprehensive service documentation
- Code examples for all M-Pesa services (STK Push, B2C, C2B, Account Balance, Transaction Status, Reversal)
- Documentation structure with docs/ directory
- Project organization with examples/ directory

### Documentation
- Added complete API reference documentation
- Created practical examples for real-world usage
- Added contributing guidelines for open source contributions
- Enhanced code comments and docstrings throughout the codebase
- Added changelog for version tracking

## [1.0.0] - 2024-08-12

### Added
- Initial release of the Go M-Pesa SDK
- STK Push (Lipa na M-Pesa Online) service
- B2C (Business to Customer) service  
- C2B (Customer to Business) service
- Account Balance service
- Transaction Status service
- Transaction Reversal service
- Comprehensive test suite with mocks
- Automatic token management
- Support for both sandbox and production environments
- Phone number validation and formatting
- Detailed error handling
- Type-safe API responses

### Security
- Encrypted credential handling
- Secure token storage and refresh
- HTTPS-only communication

---

## Release Template

When making a new release, copy this template:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Enhanced
- Improvements to existing features

### Deprecated
- Soon-to-be removed features

### Removed
- Now removed features

### Fixed
- Any bug fixes

### Security
- Security improvements

### Documentation
- Documentation updates
```

## Versioning Guidelines

- **MAJOR** version when you make incompatible API changes
- **MINOR** version when you add functionality in a backwards compatible manner
- **PATCH** version when you make backwards compatible bug fixes

Additional labels for pre-release and build metadata are available as extensions to the MAJOR.MINOR.PATCH format.
