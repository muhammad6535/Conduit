"""Semantic cache using vector similarity."""

import logging
from typing import Optional

logger = logging.getLogger(__name__)


class SemanticCache:
    def __init__(self):
        self._stores: dict[str, list[dict]] = {}
        logger.info("Semantic cache initialized (in-memory)")

    async def lookup(
        self, text: str, threshold: float = 0.92, namespace: str = "default"
    ) -> dict:
        store = self._stores.get(namespace, [])
        if not store:
            return {"found": False, "response": None, "similarity": 0.0}

        query_emb = await self._embed(text)
        best = {"similarity": 0.0, "response": None}

        for entry in store:
            sim = self._cosine_similarity(query_emb, entry["embedding"])
            if sim > best["similarity"]:
                best = {"similarity": sim, "response": entry["response"]}

        if best["similarity"] >= threshold:
            logger.info(
                "Cache HIT (sim=%.4f, threshold=%.2f)", best["similarity"], threshold
            )
            return {"found": True, "response": best["response"], "similarity": best["similarity"]}

        logger.info(
            "Cache MISS (best=%.4f, threshold=%.2f)", best["similarity"], threshold
        )
        return {"found": False, "response": None, "similarity": best["similarity"]}

    async def store(self, text: str, response: str, namespace: str = "default"):
        embedding = await self._embed(text)
        if namespace not in self._stores:
            self._stores[namespace] = []
        self._stores[namespace].append({
            "text": text,
            "response": response,
            "embedding": embedding,
        })
        logger.info("Cached response (%s): %d entries", namespace, len(self._stores[namespace]))

    async def _embed(self, text: str) -> list[float]:
        from internal.engine.embeddings.generator import generate_embedding
        return await generate_embedding(text)

    def _cosine_similarity(self, a: list[float], b: list[float]) -> float:
        dot = sum(x * y for x, y in zip(a, b))
        norm_a = sum(x * x for x in a) ** 0.5
        norm_b = sum(x * x for x in b) ** 0.5
        if norm_a == 0 or norm_b == 0:
            return 0.0
        return dot / (norm_a * norm_b)
