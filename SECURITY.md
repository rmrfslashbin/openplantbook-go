# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Currently supported versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.x.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

We take the security of the OpenPlantbook Go SDK seriously. If you believe you have found a security vulnerability, please report it to us as described below.

### Please DO NOT:
- Open a public GitHub issue for security vulnerabilities
- Discuss the vulnerability publicly until it has been addressed

### Please DO:
1. **Email** security details to: [INSERT SECURITY EMAIL]
2. **Include** as much information as possible:
   - Type of vulnerability
   - Full paths of source file(s) related to the manifestation of the vulnerability
   - Location of the affected source code (tag/branch/commit or direct URL)
   - Step-by-step instructions to reproduce the vulnerability
   - Proof-of-concept or exploit code (if possible)
   - Impact of the vulnerability, including how an attacker might exploit it

### What to Expect

- **Acknowledgment**: We will acknowledge receipt of your vulnerability report within 48 hours
- **Updates**: We will send you regular updates about our progress
- **Timeline**: We aim to address critical vulnerabilities within 7 days
- **Disclosure**: Once the vulnerability is fixed, we will work with you on disclosure timing

## Security Best Practices

When using this SDK:

### API Keys and Secrets
- **Never commit** API keys, OAuth2 credentials, or secrets to version control
- **Use environment variables** or secure secret management systems
- **Rotate credentials** regularly
- **Use .env files** locally (already gitignored)

### OAuth2 Security
- **Protect client secrets** - treat them like passwords
- **Use HTTPS** for all API communications (enforced by SDK)
- **Implement token refresh** logic if storing tokens
- **Never log** or expose tokens in error messages

### Rate Limiting
- **Respect rate limits** - the SDK includes built-in rate limiting
- **Handle 429 errors** gracefully in your applications
- **Don't disable** rate limiting in production unless necessary

### Input Validation
- **Validate user input** before passing to SDK methods
- **Sanitize search queries** to prevent injection attacks
- **Use context timeouts** to prevent resource exhaustion

### Dependencies
- **Keep dependencies updated** - run `go get -u ./...` regularly
- **Monitor security advisories** for Go and dependencies
- **Use tools** like `govulncheck` to scan for vulnerabilities

## Security Features

The SDK includes several built-in security features:

- **TLS/HTTPS enforced** - all API communications use HTTPS
- **Context support** - proper request cancellation and timeouts
- **Rate limiting** - prevents accidental API abuse
- **Error sanitization** - sensitive data not exposed in errors
- **Input validation** - validates required parameters

## Vulnerability Disclosure Policy

When a security vulnerability is fixed:

1. We will create a security advisory on GitHub
2. We will release a patch version
3. We will update the CHANGELOG with security fix details
4. We will credit the reporter (unless they prefer to remain anonymous)

## Compliance

This SDK:
- Does not store or transmit personal data beyond API requests
- Uses standard OAuth2 flows for authentication
- Follows OWASP security best practices
- Maintains audit logs through standard Go logging interfaces

## Contact

For security concerns, please contact: [INSERT SECURITY EMAIL]

For general questions, use GitHub Discussions or Issues.

---

**Thank you for helping keep OpenPlantbook Go SDK and its users safe!**
