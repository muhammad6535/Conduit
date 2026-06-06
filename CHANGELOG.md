# Changelog

## [0.2.0] - 2026-06-06

### Added
- GitHub CI workflow: lint, build, test, cross-compile, Docker
- GitHub Release workflow with GoReleaser
- Issue templates (bug report, feature request, config help)
- Pull request template
- Dependabot config for Go, Python, GitHub Actions deps
- `AGENTS.md` — AI assistant guide for the codebase
- `CLAUDE.md` — Claude-specific development guide
- `.golangci.yml` — Go lint configuration
- `.goreleaser.yaml` — Automated release pipeline
- GitHub community files (FUNDING.yml, Scorecard)
- Docker build in CI with BuildKit caching

### Changed
- Updated `.gitignore` with build/ and dist/ patterns

## [0.1.0] - 2026-06-05

### Added
- Go gateway: HTTP server, config, logging
- OpenAI provider adapter (passthrough)
- Anthropic provider adapter (format conversion)
- Cost tracker with 15+ model pricing entries
- Rule-based model router
- Streaming SSE support (OpenAI + Anthropic)
- CLI flags: `--version`, `--init`, `--config`
- Python ML engine scaffold: FastAPI service
- ML classifier (simple/medium/complex task detection)
- Semantic cache (in-memory vector similarity)
- PII scanner (6 pattern types)
- Prompt optimizer (token estimation, compression)
- Docker multi-stage build + docker-compose.yml
- Configuration profiles: cost-optimized, balanced, privacy-first
- Documentation: README, ARCHITECTURE, COST_SAVINGS, SECURITY, CONTRIBUTING
- Quickstart and benchmark scripts
- SVG logo and brand assets
- Apache 2.0 license
