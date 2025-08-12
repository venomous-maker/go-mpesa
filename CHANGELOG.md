# Changelog

All notable changes to the Go M-Pesa SDK will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of Go M-Pesa SDK
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

### Changed
- N/A (Initial release)

### Deprecated
- N/A (Initial release)

### Removed
- N/A (Initial release)

### Fixed
- N/A (Initial release)

### Security
- Encrypted credential handling
- Secure token storage and refresh
- HTTPS-only communication

## [1.0.0] - 2024-08-12

### Added
- Initial release of the Go M-Pesa SDK
- Complete M-Pesa API integration
- Full documentation and examples
- Comprehensive test coverage

---

## Release Template

When making a new release, copy this template:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Now removed features

### Fixed
- Any bug fixes

### Security
- Security improvements
```

## Versioning Guidelines

- **MAJOR** version when you make incompatible API changes
- **MINOR** version when you add functionality in a backwards compatible manner
- **PATCH** version when you make backwards compatible bug fixes

Additional labels for pre-release and build metadata are available as extensions to the MAJOR.MINOR.PATCH format.
