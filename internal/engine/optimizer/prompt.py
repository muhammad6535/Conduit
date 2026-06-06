"""Prompt optimization utilities for token reduction."""

import logging
from typing import Optional

logger = logging.getLogger(__name__)


def estimate_tokens(text: str) -> int:
    """Rough estimate: ~4 chars per token for English text."""
    return len(text) // 4


def compress_messages(messages: list[dict], max_tokens: int = 4096) -> list[dict]:
    """Compress messages to fit within token budget."""
    total_tokens = sum(estimate_tokens(m.get("content", "")) for m in messages)

    if total_tokens <= max_tokens:
        return messages

    logger.info("Compressing %d tokens to %d", total_tokens, max_tokens)

    compressed = list(messages)
    while compressed and estimate_tokens(
        " ".join(m.get("content", "") for m in compressed)
    ) > max_tokens:
        oldest_content = compressed[0].get("content", "")
        if len(oldest_content) > 100:
            compressed[0]["content"] = oldest_content[:len(oldest_content) // 2] + "..."
        elif len(compressed) > 1:
            compressed.pop(0)
        else:
            break

    return compressed


def trim_history(
    messages: list[dict], max_context: int = 4
) -> list[dict]:
    """Keep only system prompt and last N exchanges."""
    if len(messages) <= max_context + 1:
        return messages

    system_msgs = [m for m in messages if m.get("role") == "system"]
    recent = messages[-(max_context):]

    return system_msgs + recent
