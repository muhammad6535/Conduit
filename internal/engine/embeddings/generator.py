"""Embedding generation for semantic caching."""

import logging
from typing import Optional

logger = logging.getLogger(__name__)

_embedding_model = None


def _get_model():
    global _embedding_model
    if _embedding_model is None:
        try:
            from sentence_transformers import SentenceTransformer
            _embedding_model = SentenceTransformer("all-MiniLM-L6-v2")
            logger.info("Loaded embedding model: all-MiniLM-L6-v2")
        except ImportError:
            logger.warning(
                "sentence-transformers not installed. "
                "Install with: pip install sentence-transformers"
            )
            return None
    return _embedding_model


async def generate_embedding(text: str) -> list[float]:
    model = _get_model()
    if model is None:
        return _fake_embedding(text)

    embedding = model.encode(text, normalize_embeddings=True)
    return embedding.tolist()


def _fake_embedding(text: str) -> list[float]:
    """Fallback deterministic embedding for when model isn't available."""
    import hashlib
    h = hashlib.md5(text.encode()).digest()
    return [b / 255.0 for b in h] + [0.0] * 384
