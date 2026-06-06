# Security Policy

## Supported Versions

| Version | Supported |
|---|---|
| 0.x (alpha) | ✅ |

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please report it privately:

- **GitHub**: Open a private security advisory at https://github.com/muhammad6535/conduit/security/advisories
- **Do not** open a public GitHub issue.

We will acknowledge receipt within 48 hours and provide a timeline for a fix.

## Security Features

- API keys via environment variables only (never in config files)
- PII redaction before data leaves your network (Phase 2)
- Rate limiting to prevent abuse
- No telemetry or external calls from the gateway
- Configurable request logging (no sensitive content by default)

## Responsible Disclosure

We ask that researchers:
1. Report privately first
2. Give us reasonable time to fix before disclosure
3. Do not exfiltrate data during testing
