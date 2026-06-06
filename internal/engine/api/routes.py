"""Engine API routes for classification, embeddings, and caching."""

from typing import Optional
from fastapi import APIRouter, Request
from pydantic import BaseModel

router = APIRouter()


class ClassifyRequest(BaseModel):
    messages: list[dict]
    model_hint: Optional[str] = None


class ClassifyResponse(BaseModel):
    complexity: str  # "simple" | "medium" | "complex"
    task_type: str   # "chat" | "code" | "reasoning" | "analysis" | "creative"
    recommended_model_tier: str  # "cheap" | "balanced" | "best"
    confidence: float


@router.post("/classify", response_model=ClassifyResponse)
async def classify(req: ClassifyRequest, request: Request):
    classifier = request.app.state.classifier
    result = classifier.classify(req.messages, req.model_hint)
    return result


class EmbedRequest(BaseModel):
    text: str


class EmbedResponse(BaseModel):
    embedding: list[float]
    model: str
    dimensions: int


@router.post("/embeddings", response_model=EmbedResponse)
async def embed(req: EmbedRequest):
    from internal.engine.embeddings.generator import generate_embedding
    emb = await generate_embedding(req.text)
    return EmbedResponse(
        embedding=emb,
        model="sentence-transformers/all-MiniLM-L6-v2",
        dimensions=len(emb),
    )


class CacheLookupRequest(BaseModel):
    text: str
    threshold: float = 0.92
    namespace: str = "default"


class CacheLookupResponse(BaseModel):
    found: bool
    response: Optional[str] = None
    similarity: float = 0.0


@router.post("/cache/lookup", response_model=CacheLookupResponse)
async def cache_lookup(req: CacheLookupRequest, request: Request):
    cache = request.app.state.cache
    result = await cache.lookup(req.text, req.threshold, req.namespace)
    return result


class CacheStoreRequest(BaseModel):
    text: str
    response: str
    namespace: str = "default"


@router.post("/cache/store")
async def cache_store(req: CacheStoreRequest, request: Request):
    cache = request.app.state.cache
    await cache.store(req.text, req.response, req.namespace)
    return {"status": "stored"}


class PrivacyScanRequest(BaseModel):
    text: str


class PrivacyScanResponse(BaseModel):
    has_pii: bool
    redacted_text: str
    findings: list[dict]


@router.post("/privacy/scan", response_model=PrivacyScanResponse)
async def privacy_scan(req: PrivacyScanRequest):
    from internal.engine.privacy.scanner import scan_pii
    result = await scan_pii(req.text)
    return result
