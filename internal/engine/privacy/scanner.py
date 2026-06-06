"""PII detection and redaction using regex patterns."""

import re
import logging

logger = logging.getLogger(__name__)

PII_PATTERNS = [
    ("email", re.compile(r"\b[\w.+-]+@[\w-]+\.[\w.-]+\b")),
    ("phone", re.compile(r"\b\+?\d[\d\s\-().]{7,}\d\b")),
    ("ssn", re.compile(r"\b\d{3}-\d{2}-\d{4}\b")),
    ("credit_card", re.compile(r"\b(?:\d[ -]*?){13,16}\b")),
    ("ip_address", re.compile(r"\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b")),
    ("api_key", re.compile(r"\b(sk-[a-zA-Z0-9]{20,}|sk-ant-[a-zA-Z0-9]{20,})\b")),
]

REDACT_TOKEN = "[REDACTED]"


async def scan_pii(text: str) -> dict:
    findings = []
    redacted = text

    for pii_type, pattern in PII_PATTERNS:
        for match in pattern.finditer(text):
            findings.append({
                "type": pii_type,
                "start": match.start(),
                "end": match.end(),
                "value": match.group(),
            })
            redacted = redacted.replace(match.group(), REDACT_TOKEN)

    has_pii = len(findings) > 0
    if has_pii:
        logger.info("Found %d PII items: %s", len(findings), [f["type"] for f in findings])

    return {
        "has_pii": has_pii,
        "redacted_text": redacted,
        "findings": findings,
    }
