# Pack Calculator - Optimal Packing Solution

A production-ready microservice that solves the **bin packing problem** using dynamic programming. Given a set of pack sizes and a required amount, it finds the optimal combination that minimizes overage and pack count.

## 🌐 Live Demo

**🚀 Public URL:** [https://re-partners-api-production-f47f.up.railway.app](https://re-partners-api-production-f47f.up.railway.app)

**Quick Test:**
```bash
# Health check
curl https://re-partners-api-production-f47f.up.railway.app/healthz

# Web UI
open https://re-partners-api-production-f47f.up.railway.app

# API Example
curl -X POST https://re-partners-api-production-f47f.up.railway.app/packs/solve \
  -H "Content-Type: application/json" \
  -d '{"sizes": [250, 500, 1000, 2000, 5000], "amount": 12001}'
```

## 📋 Overview

This service implements an **optimal packing algorithm** to determine the minimum number of packs needed to fulfill an order, with the following priorities:

1. **Minimize overage** (excess items delivered)
2. **Minimize pack count** (when overage is equal)

### Why Dynamic Programming?

The greedy approach (always picking the largest pack) **fails** for many cases:

```
Packs: [3, 5]
Amount: 9

Greedy: 5 + 3 + 3 = 11 (overage=2, packs=3) ❌
DP:     3 + 3 + 3 = 9  (overage=0, packs=3) ✓
```

## 🚀 API

### Endpoints

- `GET /` - Web UI (interactive calculator)
- `GET /healthz` - Health check
- `POST /packs/solve` - Solve packing problem

### Request Example

```bash
curl -X POST http://localhost:8080/packs/solve \
  -H "Content-Type: application/json" \
  -d '{
    "sizes": [250, 500, 1000, 2000, 5000],
    "amount": 12001
  }'
```

### Response Example

```json
{
  "solution": {
    "250": 1,
    "2000": 1,
    "5000": 2
  },
  "overage": 249,
  "packs": 4
}
```

## ⚡ Quickstart

### Requirements

- Go 1.23 or higher
- Make (optional)

### Run Locally

```bash
# Clone the repository
git clone https://github.com/chainwest/re-partners-assignment.git
cd re-partners-assignment

# Run the application
make run
# Or: go run cmd/api/main.go
```

The service will start on `http://localhost:8080`

### Run Tests

```bash
make test
# Or: go test -v ./...
```

## 🐳 Docker

### Quick Start with Docker Compose

```bash
# Start all services (API, PostgreSQL, Redis)
docker-compose -f deployments/docker-compose.yaml up

# Stop services
docker-compose -f deployments/docker-compose.yaml down
```

### Standalone Docker

```bash
# Build image
docker build -f deployments/Dockerfile -t pack-calculator .

# Run container
docker run -p 8080:8080 pack-calculator
```

## 📊 Performance

**Platform:** Apple M1 Pro

| Scenario | Time | Memory |
|----------|------|--------|
| 500k items with [23,31,53] | 3.69ms | 4.0MB |
| 500k items with [100,250,500,1000,2500] | 0.93ms | 4.0MB |
| Edge case: 500k with 10 prime numbers | 9.55ms | 4.4MB |

✅ **52x faster** than required 500ms

## 🏗️ Architecture

Clean Architecture with strict separation of concerns:

```
cmd/api/              # Application entry point
internal/
  ├── domain/         # Business entities and interfaces
  ├── usecase/        # Core algorithm (DP solver)
  ├── adapters/http/  # HTTP handlers
  └── infra/          # Infrastructure (DB, cache, config)
web/                  # Static files (UI)
deployments/          # Docker, K8s configs
```

**Key Principle:** Usecase is independent of infrastructure - all calculations in memory, no DB dependencies.

## 📚 Documentation

- **[API.md](API.md)** - Complete API documentation
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Architecture details

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Built with ❤️ using Go, Clean Architecture, and Dynamic Programming**
