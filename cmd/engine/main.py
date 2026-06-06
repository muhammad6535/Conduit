"""Conduit ML Engine - sidecar service for intelligent routing & caching."""

import os
import logging
from contextlib import asynccontextmanager

import uvicorn
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from internal.engine.api.routes import router
from internal.engine.cache.store import SemanticCache
from internal.engine.classifier.router import RequestClassifier

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("conduit-engine")


@asynccontextmanager
async def lifespan(app: FastAPI):
    logger.info("Conduit Engine starting...")
    app.state.cache = SemanticCache()
    app.state.classifier = RequestClassifier()
    yield
    logger.info("Conduit Engine stopped.")


app = FastAPI(
    title="Conduit Engine",
    version="0.1.0",
    lifespan=lifespan,
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_methods=["*"],
    allow_headers=["*"],
)

app.include_router(router, prefix="/v1")


@app.get("/health")
async def health():
    return {"status": "ok", "service": "conduit-engine"}


if __name__ == "__main__":
    port = int(os.getenv("ENGINE_PORT", "9091"))
    uvicorn.run("main:app", host="0.0.0.0", port=port, reload=True)
