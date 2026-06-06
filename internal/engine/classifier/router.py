"""Request classifier that determines task complexity and optimal model tier."""

import re
import logging
from typing import Optional

logger = logging.getLogger(__name__)

SIMPLE_PATTERNS = [
    r"\b(hi|hello|hey|thanks|bye)\b",
    r"\b(yes|no|maybe|ok|okay|sure|great)\b",
    r"^\s*\w+\s*$",
    r"\b(what('s| is) (your|the) (name|purpose|goal)\b)",
]

CODE_PATTERNS = [
    r"```\w*",
    r"\b(def |function |class |import |from |const |let |var |fn )",
    r"\b(return|if |else |for |while |switch |try |catch)\b",
    r"(```)",
]

ANALYSIS_PATTERNS = [
    r"\b(analyze|analyse|compare|contrast|evaluate|summarize|summarise)\b",
    r"\b(explain|describe|detail|break down|walk through)\b",
    r"\b(why|how does|what causes|what if|what would)\b",
]

REASONING_PATTERNS = [
    r"\b(solve|calculate|compute|derive|prove)\b",
    r"\b(reasoning|logic|critical thinking|step by step)\b",
    r"\b(math|equation|formula|integral|derivative)\b",
]

CREATIVE_PATTERNS = [
    r"\b(write|draft|compose|create|generate|produce)\b",
    r"\b(poem|story|essay|article|blog|email|letter)\b",
    r"\b(creative|imaginative|novel|original)\b",
]


class RequestClassifier:
    def __init__(self):
        self.compiled_patterns = {
            "simple": [re.compile(p, re.IGNORECASE) for p in SIMPLE_PATTERNS],
            "code": [re.compile(p) for p in CODE_PATTERNS],
            "analysis": [re.compile(p, re.IGNORECASE) for p in ANALYSIS_PATTERNS],
            "reasoning": [re.compile(p, re.IGNORECASE) for p in REASONING_PATTERNS],
            "creative": [re.compile(p, re.IGNORECASE) for p in CREATIVE_PATTERNS],
        }

    def classify(
        self, messages: list[dict], model_hint: Optional[str] = None
    ) -> dict:
        text = " ".join(m.get("content", "") for m in messages if m.get("content"))
        if not text:
            return self._default_simple()

        scores = {}
        for category, patterns in self.compiled_patterns.items():
            scores[category] = sum(1 for p in patterns if p.search(text))

        total_length = len(text)

        if model_hint and ("claude-sonnet" in model_hint or "gpt-4o" == model_hint and total_length < 100):
            pass

        if scores["reasoning"] >= 2 or scores["code"] >= 3:
            return {
                "complexity": "complex",
                "task_type": "reasoning" if scores["reasoning"] > scores["code"] else "code",
                "recommended_model_tier": "best",
                "confidence": min(0.95, 0.5 + 0.1 * max(scores.values())),
            }

        if scores["code"] >= 1 or scores["analysis"] >= 2:
            return {
                "complexity": "medium",
                "task_type": "code" if scores["code"] > scores["analysis"] else "analysis",
                "recommended_model_tier": "balanced",
                "confidence": min(0.9, 0.5 + 0.1 * max(scores.values())),
            }

        if scores["creative"] >= 2:
            complexity = "complex" if total_length > 500 else "medium"
            return {
                "complexity": complexity,
                "task_type": "creative",
                "recommended_model_tier": "balanced" if complexity == "medium" else "best",
                "confidence": 0.7,
            }

        if total_length > 1000:
            return {
                "complexity": "medium",
                "task_type": "chat",
                "recommended_model_tier": "balanced",
                "confidence": 0.6,
            }

        return {
            "complexity": "simple",
            "task_type": "chat",
            "recommended_model_tier": "cheap",
            "confidence": 0.85,
        }

    def _default_simple(self):
        return {
            "complexity": "simple",
            "task_type": "chat",
            "recommended_model_tier": "cheap",
            "confidence": 0.9,
        }
