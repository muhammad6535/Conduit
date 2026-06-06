# Contributing to Conduit

We welcome contributions! Here's how to get started.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/your-username/gateway.git`
3. Create a feature branch: `git checkout -b feat/my-feature`
4. Make your changes
5. Run tests: `go test ./...`
6. Submit a PR

## Development Setup

### Go Gateway

```bash
go build
./gateway --config config.test.yaml
```

### Python Engine

```bash
cd cmd/engine
pip install -e ../..
uvicorn main:app --reload --port 9091
```

## Code Style

- **Go**: Follow standard `gofmt` formatting. Run `golint` and `go vet` before committing.
- **Python**: Follow PEP 8. Run `ruff check` before committing.
- Keep functions focused and small.
- Write tests for new features.
- Update documentation for API changes.

## Pull Request Process

1. Ensure your code compiles and passes tests
2. Update README.md if needed
3. Add a clear description of what your PR does
4. Reference any related issues

## Feature Requests & Bugs

Open an issue with:
- Clear description of the problem
- Steps to reproduce (for bugs)
- Expected behavior
- Environment details

## License

By contributing, you agree that your contributions will be licensed under Apache 2.0.
