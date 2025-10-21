# Security Policy

## Supported Versions

We take security seriously and provide security updates for the following versions of ezBookkeeping:

| Version | Supported          | End of Support |
| ------- | ------------------ | -------------- |
| 1.1.x   | :white_check_mark: | Current        |
| 1.0.x   | :white_check_mark: | TBD            |
| < 1.0   | :x:                | Ended          |

**Note:** We recommend always running the latest stable version to ensure you have the latest security patches and features.

## Security Features

ezBookkeeping includes several built-in security features to protect your financial data:

- **Two-Factor Authentication (2FA)**: Enhanced account security with TOTP support
- **Application Lock**: PIN code or WebAuthn-based app locking
- **Login Rate Limiting**: Protection against brute-force attacks
- **Encrypted Data Storage**: Secure storage of sensitive financial information
- **Self-Hosted Architecture**: Your data stays on your own infrastructure

## Reporting a Vulnerability

We appreciate the security research community's efforts in responsibly disclosing vulnerabilities. If you believe you've found a security issue in ezBookkeeping, please follow these guidelines:

### Where to Report

**DO NOT** report security vulnerabilities through public GitHub issues.

Instead, please report security vulnerabilities by:

1. **Email**: Send an email to the project maintainer via the contact information in the [GitHub profile](https://github.com/mayswind)
2. **GitHub Security Advisory**: Use GitHub's private security advisory feature at https://github.com/mayswind/ezbookkeeping/security/advisories/new

### What to Include

To help us understand and resolve the issue quickly, please include as much of the following information as possible:

- **Type of vulnerability** (e.g., SQL injection, XSS, authentication bypass, etc.)
- **Full paths of source file(s)** related to the vulnerability
- **Location of the affected source code** (tag/branch/commit or direct URL)
- **Step-by-step instructions** to reproduce the issue
- **Proof-of-concept or exploit code** (if possible)
- **Impact of the vulnerability** and potential attack scenarios
- **Any suggested fixes** or mitigations you may have

### Response Timeline

We are committed to working with security researchers and will respond according to the following timeline:

- **Initial Response**: Within 48-72 hours of receiving your report
- **Confirmation**: Within 7 days - we will confirm whether the issue is valid and its severity
- **Status Updates**: Every 7-14 days until the issue is resolved
- **Resolution**: Timeline depends on severity:
  - **Critical**: 7-14 days
  - **High**: 14-30 days
  - **Medium**: 30-60 days
  - **Low**: 60-90 days

### Disclosure Policy

- We request that you give us reasonable time to investigate and fix the issue before public disclosure
- We will acknowledge your responsible disclosure in our security advisories and release notes (unless you prefer to remain anonymous)
- We will work with you to understand and resolve the issue promptly
- Once the vulnerability is fixed, we will coordinate with you on the disclosure timeline

### What to Expect

When we receive a security report:

- ✅ **Accepted**: If the vulnerability is confirmed, we will:
  - Work on a fix as a priority
  - Keep you informed of our progress
  - Credit you in the security advisory (if desired)
  - Release a security patch and advisory

- ❌ **Declined**: If we determine the report is not a security vulnerability, we will:
  - Provide a detailed explanation of why it was declined
  - Suggest alternative reporting channels if appropriate (e.g., bug report for non-security issues)

## Security Best Practices for Deployment

To ensure your ezBookkeeping installation remains secure, we recommend:

### General Recommendations

1. **Keep Updated**: Always run the latest stable version
2. **Strong Passwords**: Enforce strong password policies for all users
3. **Enable 2FA**: Require two-factor authentication for all accounts
4. **HTTPS Only**: Always use HTTPS in production environments
5. **Regular Backups**: Maintain encrypted backups of your database

### Self-Hosting Security

If you're self-hosting ezBookkeeping:

1. **Firewall Configuration**: Restrict access to only necessary ports
2. **Database Security**: 
   - Use strong database passwords
   - Limit database access to the application only
   - Keep database software updated
3. **Reverse Proxy**: Use a reverse proxy (e.g., Nginx, Caddy) with SSL/TLS
4. **Container Security**: If using Docker, follow Docker security best practices
5. **Network Isolation**: Run ezBookkeeping in an isolated network segment when possible

### Configuration Security

1. **Review Default Settings**: Change all default credentials and keys
2. **Disable Unnecessary Features**: Only enable features you actively use
3. **Monitor Logs**: Regularly review application logs for suspicious activity
4. **Rate Limiting**: Configure appropriate rate limits to prevent abuse
5. **Session Security**: Configure secure session timeouts

## Security Advisories

Security advisories will be published in the following locations:

- **GitHub Security Advisories**: https://github.com/mayswind/ezbookkeeping/security/advisories
- **Release Notes**: https://github.com/mayswind/ezbookkeeping/releases
- **Documentation**: https://ezbookkeeping.mayswind.net

## Scope

### In Scope

The following are considered in scope for security reports:

- The core ezBookkeeping application (backend and frontend)
- Authentication and authorization mechanisms
- Data storage and encryption
- API endpoints
- Docker images and deployment configurations
- Dependencies with known vulnerabilities

### Out of Scope

The following are generally not eligible for security reports:

- Issues in third-party dependencies already reported upstream
- Social engineering attacks
- Issues requiring physical access to servers
- Denial of Service (DoS) attacks against the demo site
- Issues specific to outdated/unsupported versions
- Theoretical vulnerabilities without proof of exploitability

## Safe Harbor

We support safe harbor for security researchers who:

- Make a good faith effort to avoid privacy violations, data destruction, and service interruption
- Only interact with accounts you own or with explicit permission from the account holder
- Do not exploit a security issue beyond what is necessary to demonstrate the vulnerability
- Report vulnerabilities promptly and keep details confidential until we've resolved the issue

If you comply with these guidelines, we will not pursue legal action or ask law enforcement to investigate you.

## Recognition

We appreciate the security research community and, with your permission, will recognize researchers who report valid security issues in:

- Our security advisories
- Release notes
- A potential security researchers acknowledgment page (if we create one in the future)

You may choose to remain anonymous if you prefer.

## Additional Resources

- **Project Homepage**: https://github.com/mayswind/ezbookkeeping
- **Documentation**: https://ezbookkeeping.mayswind.net
- **Live Demo**: https://ezbookkeeping-demo.mayswind.net (for testing purposes only - do not test vulnerabilities here)
- **Docker Hub**: https://hub.docker.com/r/mayswind/ezbookkeeping

## Questions?

If you have any questions about this security policy or the security of ezBookkeeping, please reach out through the reporting channels mentioned above.

Thank you for helping keep ezBookkeeping and our users safe!
