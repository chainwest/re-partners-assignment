# Pack Calculator - Optimal Packing Solution

[![CI](https://github.com/evgenijurbanovskij/re-partners-assignment/actions/workflows/ci.yml/badge.svg)](https://github.com/evgenijurbanovskij/re-partners-assignment/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/evgenijurbanovskij/re-partners-assignment/branch/main/graph/badge.svg)](https://codecov.io/gh/evgenijurbanovskij/re-partners-assignment)
[![Go Report Card](https://goreportcard.com/badge/github.com/evgenijurbanovskij/re-partners-assignment)](https://goreportcard.com/report/github.com/evgenijurbanovskij/re-partners-assignment)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A production-ready microservice that solves the **bin packing problem** using dynamic programming. Given a set of pack sizes and a required amount, it finds the optimal combination that minimizes overage and pack count.

## 🌐 Live Demo

**🚀 Public URL:** [https://re-partners-api.onrender.com](https://re-partners-api.onrender.com)

> ⚠️ **Note:** Free tier sleeps after 15 minutes of inactivity. First request may take 30-60 seconds to "wake up" the service.

**Quick Test:**
```bash
# Health check
curl https://re-partners-api.onrender.com/healthz

# Web UI
open https://re-partners-api.onrender.com

# API Example
curl -X POST https://re-partners-api.onrender.com/packs/solve \
  -H "Content-Type: application/json" \
  -d '{"sizes": [250, 500, 1000, 2000, 5000], "amount": 12001}'
```

## 📋 Overview

This service implements an **optimal packing algorithm** to determine the minimum number of packs needed to fulfill an order, with the following priorities:

1. **Minimize overage** (excess items delivered)
2. **Minimize pack count** (when overage is equal)

### Why Dynamic Programming over Greedy?

The greedy approach (always picking the largest pack) **fails** for many cases. Here's a classic counter-example:

**Change-Making Problem:**
```
Packs: [3, 5]
Amount: 9

Greedy: 5 + 3 + 3 = 11 (overage=2, packs=3) ❌
DP:     3 + 3 + 3 = 9  (overage=0, packs=3) ✓
```

**Another Example:**
```
Packs: [3, 5]
Amount: 7

Greedy: 5 + 3 = 8 (overage=1, packs=2) ✓
DP:     5 + 3 = 8 (overage=1, packs=2) ✓ (same result)

But for amount=6:
Greedy: 5 + 3 = 8 (overage=2, packs=2) ❌
DP:     3 + 3 = 6 (overage=0, packs=2) ✓
```

**Dynamic Programming guarantees:**
- ✅ **Correctness:** Always finds the optimal solution
- ✅ **Completeness:** Explores all possible combinations
- ✅ **Optimality:** Minimizes overage first, then pack count
- ✅ **Performance:** O(W × N) time complexity with memoization

### Key Features

- ✅ **Optimal algorithm** - Dynamic programming with proven correctness
- ✅ **Clean architecture** - Domain-driven design with clear separation of concerns
- ✅ **Production-ready** - Comprehensive testing, CI/CD, monitoring
- ✅ **High performance** - ~9.5ms for 500k items (52x faster than required 500ms)
- ✅ **Redis caching** - 6-18x speedup for repeated queries
- ✅ **PostgreSQL audit** - Optional calculation history
- ✅ **Web UI** - Interactive calculator interface
- ✅ **Docker support** - Full containerization with docker-compose
- ✅ **Graceful shutdown** - Proper signal handling
- ✅ **Structured logging** - JSON logs with correlation IDs
- ✅ **Metrics** - Request counts and latencies

## 🎯 Rules (Priority Order)

The solver follows these rules in strict priority order:

1. **Minimum Overage:** Minimize excess items delivered
   - Example: `amount=251, sizes=[250,500]` → `{500:1}` (overage=249)
   - Not `{250:2}` (overage=249) because next rule applies

2. **Minimum Packs:** When overage is equal, minimize pack count
   - Example: `amount=1250, sizes=[250,500,1000]` → `{250:1, 1000:1}` (2 packs)
   - Not `{250:5}` (5 packs) even though overage is same (0)

3. **Always Cover:** Solution must cover at least the requested amount
   - Never deliver less than requested
   - Overage is acceptable, shortage is not

### Examples

| Amount | Sizes | Solution | Packs | Overage | Explanation |
|--------|-------|----------|-------|---------|-------------|
| 250 | [250,500,1000] | {250:1} | 1 | 0 | Exact match |
| 251 | [250,500,1000] | {500:1} | 1 | 249 | Min overage |
| 1250 | [250,500,1000] | {250:1, 1000:1} | 2 | 0 | Exact, min packs |
| 12001 | [250,500,1000,2000,5000] | {250:1, 2000:1, 5000:2} | 4 | 249 | Min overage, then min packs |

## 🚀 API

### Endpoints

- `GET /` - Web UI (interactive calculator)
- `GET /healthz` - Health check
- `GET /version` - Version information
- `POST /packs/solve` - Solve packing problem
- `GET /metrics` - Application metrics

### Request Example

```bash
curl -X POST http://localhost:8080/packs/solve \
  -H "Content-Type: application/json" \
  -H "X-Correlation-ID: my-request-123" \
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

### Validation Rules

- `sizes`: non-empty array, all values > 0 and ≤ 1,000,000
- `amount`: > 0 and ≤ 1,000,000,000
- No duplicate sizes allowed

**Full API documentation:** See [API.md](API.md) for complete details including error handling, correlation IDs, and monitoring.

## ⚡ Quickstart

### Requirements

- Go 1.21 or higher
- Make (optional but recommended)

### Run Locally

```bash
# Clone the repository
git clone https://github.com/evgenijurbanovskij/re-partners-assignment.git
cd re-partners-assignment

# Run the application
make run

# Or directly with go
go run cmd/api/main.go
```

The service will start on `http://localhost:8080`

### Test the API

```bash
# Health check
curl http://localhost:8080/healthz

# Web UI
open http://localhost:8080

# Solve a packing problem
curl -X POST http://localhost:8080/packs/solve \
  -H "Content-Type: application/json" \
  -d '{"sizes": [250, 500, 1000], "amount": 1250}'
```

### Run Tests

```bash
# All tests
make test

# With coverage
go test -v -cover ./...

# Benchmarks
go test -bench=. ./internal/usecase/... -benchmem
```

## 🐳 Docker

### Quick Start with Docker Compose

```bash
# Start all services (API, PostgreSQL, Redis)
make compose-up

# Check status
make compose-ps

# View logs
make compose-logs

# Run smoke tests
make smoke

# Stop all services
make compose-down
```

### Standalone Docker

```bash
# Build image
make docker-build

# Run container
make docker-run
```

### Architecture

**Multi-stage Dockerfile:**
- Build stage: Compile Go application in Alpine Linux
- Final stage: Minimal runtime image (~20MB)
- Optimizations: `-ldflags="-s -w"`, `-trimpath`, `CGO_ENABLED=0`
- Security: non-root user, minimal dependencies

**Docker Compose:**
- **API**: port 8080, depends on PostgreSQL and Redis
- **PostgreSQL**: port 5432, automatic migrations on startup
- **Redis**: port 6379, configured LRU cache (256MB)
- Health checks for all services
- Persistent volumes for data

## 🔄 CI/CD

### GitHub Actions Workflow

The project includes a comprehensive CI/CD pipeline:

**Jobs:**
1. **Lint** - Code formatting and static analysis (gofmt, go vet, staticcheck)
2. **Test** - Unit tests with coverage (93.6%+) and race detection
3. **Benchmark** - Performance tests with results reporting
4. **Docker** - Build and push images to GitHub Container Registry
5. **Post-Deploy** - Health checks and smoke tests
6. **Summary** - Consolidated build report

**Triggers:**
- Push to `main` or `develop` branches
- Pull requests
- Manual workflow dispatch

**Services:**
- PostgreSQL 15 (for integration tests)
- Redis 7 (for cache tests)

**Artifacts:**
- Coverage reports (30 days retention)
- Benchmark results (90 days retention)
- Docker images (tagged with branch and SHA)

**View CI:** [GitHub Actions](https://github.com/evgenijurbanovskij/re-partners-assignment/actions)

## 📊 Benchmarks

### Performance Results

**Platform:** Apple M1 Pro (darwin/arm64)

| Scenario | Time | Memory | Allocations |
|----------|------|--------|-------------|
| 500k with small denominations [23,31,53] | 3.69ms | 4.0MB | 5 |
| 500k with medium denominations [100,250,500,1000,2500] | 0.93ms | 4.0MB | 5 |
| 500k with large denominations [1000,2500,5000,10000,25000] | 0.93ms | 4.0MB | 5 |
| 500k with mixed denominations (7 sizes) | 6.13ms | 4.0MB | 5 |
| 500k with 10 denominations | 8.19ms | 4.0MB | 11 |
| 450k edge case [19,29,41,59,83] | 4.54ms | 3.6MB | 5 |
| 550k edge case (10 primes) | 9.55ms | 4.4MB | 11 |

**Key Metrics:**
- ✅ **Maximum time:** ~9.5ms
- ✅ **Requirement:** ≤500ms
- ✅ **Performance margin:** 52x faster than required
- ✅ **Memory usage:** ~4MB per operation
- ✅ **Allocations:** 5-11 (minimal)

### Standard Benchmarks

```
BenchmarkDPSolver_SmallAmount    47112    25055 ns/op    98576 B/op    5 allocs/op
BenchmarkDPSolver_MediumAmount    5581   204388 ns/op   803089 B/op    5 allocs/op
BenchmarkDPSolver_LargeAmount     1314   885194 ns/op  4014338 B/op    5 allocs/op
```

### Run Benchmarks

```bash
# All benchmarks
go test -bench=. ./internal/usecase/... -benchmem

# Large-scale benchmarks only
go test -bench=BenchmarkSolveLarge ./internal/usecase/... -benchmem

# With timeout for long tests
go test -bench=. ./internal/usecase/... -benchmem -timeout 10m
```

## 🎯 Edge Cases

The solver handles various edge cases correctly:

### 1. Exact Match
```
Input:  sizes=[250,500,1000], amount=250
Output: {250:1}, packs=1, overage=0
```

### 2. Minimum Overage
```
Input:  sizes=[250,500,1000], amount=251
Output: {500:1}, packs=1, overage=249
```

### 3. Non-Multiple Amounts
```
Input:  sizes=[250,500,1000], amount=333
Output: {500:1}, packs=1, overage=167
```

### 4. Large Amount with Small Denominations
```
Input:  sizes=[23,31,53], amount=500000
Output: {23:2, 31:7, 53:9429}, packs=9438, overage=0
Verification: 23×2 + 31×7 + 53×9429 = 46 + 217 + 499737 = 500000 ✓
Time: ~3.69ms
```

### 5. Tie-Breaking (Equal Overage)
```
Input:  sizes=[3,5], amount=7
Output: {3:1, 5:1}, packs=2, overage=1
Not:    {3:3}, packs=3, overage=2 (more overage)
```

### 6. Unordered Sizes
```
Input:  sizes=[5000,250,1000,500,2000], amount=12001
Output: {250:1, 2000:1, 5000:2}, packs=4, overage=249
(Solver automatically normalizes and sorts)
```

### 7. Single Pack Size
```
Input:  sizes=[500], amount=1000
Output: {500:2}, packs=2, overage=0
```

### 8. Prime Numbers (Worst Case)
```
Input:  sizes=[11,13,17,19,23,29,31,37,41,43], amount=550000
Output: Optimal solution in ~9.55ms
(Coprime numbers are hardest for DP)
```

## 🎨 Web UI

The project includes an interactive web interface:

**Features:**
- ✅ Dynamic pack size management (add/remove)
- ✅ Client-side validation
- ✅ Results displayed in table format
- ✅ Auto-fill edge case example (263 items)
- ✅ Friendly error messages
- ✅ Responsive design

**Access:** Open `http://localhost:8080` in your browser

**Demo:** [Live Web UI](https://re-partners-api.onrender.com)

<!-- TODO: Add GIF demo here -->

## 🏗️ Architecture

The project follows **clean architecture** principles with strict separation of concerns:

```
┌─────────────────────────────────────────┐
│           HTTP Handler                  │
│      (adapters/http)                    │
└──────────────┬──────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────┐
│         Domain Interfaces               │
│      domain.Solver                      │
└──────────────┬──────────────────────────┘
               │
      ┌────────┴────────┐
      ▼                 ▼
┌──────────┐    ┌──────────────┐
│ Usecase  │    │ Redis Cache  │ (optional)
│ (core)   │    │ (wrapper)    │
└──────────┘    └──────────────┘
      │
      │ NO DATABASE DEPENDENCY!
      │
      ▼
   Pure Go
   Algorithms
```

### Layers

- **cmd/api** - Entry point, server initialization
- **internal/domain** - Business logic, independent of external dependencies
- **internal/usecase** - Use cases, business logic orchestration (**NO DB DEPENDENCY**)
- **internal/adapters** - Adapters for external interfaces (HTTP, gRPC)
- **internal/infra** - Infrastructure layer (DB, cache, config, logs)

### Key Principles

1. **Usecase is independent of infrastructure**
   - All calculations performed in memory
   - No dependencies on DB, cache, or external services
   - Easy to test without mocks

2. **Infrastructure is optional**
   - PostgreSQL - for audit and demonstration
   - Redis - for accelerating repeated requests
   - Application works without them

3. **Dependency Inversion**
   - Domain defines interfaces (`domain.Solver`)
   - Infrastructure implements interfaces
   - Dependencies point inward (toward domain)

## 🚧 Limitations & Future Work

### Current Limitations

1. **Memory Constraints**
   - Maximum DP table size: 10M elements (~40MB)
   - Very large amounts (>10M) may hit memory limits
   - Mitigation: Early exit optimizations, memory-efficient int32

2. **Single-threaded**
   - Each request processed sequentially
   - No parallel DP computation
   - Acceptable for current performance (9.5ms max)

3. **No Pack Inventory**
   - Assumes unlimited pack availability
   - Doesn't check stock levels
   - Real-world systems need inventory integration

4. **Fixed Pack Sizes**
   - Pack sizes must be predefined
   - No dynamic size generation
   - No custom pack creation

### Future Improvements

**Performance:**
- [ ] Parallel DP for multi-core utilization
- [ ] Approximate algorithms for very large amounts (>10M)
- [ ] Incremental computation for similar queries
- [ ] GPU acceleration for massive parallelism

**Features:**
- [ ] Pack inventory management
- [ ] Cost optimization (minimize cost, not just packs)
- [ ] Multi-objective optimization (cost + delivery time + carbon footprint)
- [ ] Batch processing for multiple orders
- [ ] Real-time pack size recommendations

**Infrastructure:**
- [ ] gRPC API for high-performance clients
- [ ] GraphQL for flexible queries
- [ ] Distributed caching (Redis Cluster)
- [ ] Read replicas for PostgreSQL
- [ ] Kubernetes deployment with auto-scaling

**Monitoring:**
- [ ] Prometheus metrics export
- [ ] Grafana dashboards
- [ ] Distributed tracing (OpenTelemetry)
- [ ] Alerting on performance degradation

**Testing:**
- [ ] Property-based testing (go-fuzz)
- [ ] Chaos engineering tests
- [ ] Load testing (k6, Locust)
- [ ] Mutation testing

## 📚 Documentation

- **[API.md](API.md)** - Complete API documentation
- **[ARCHITECTURE.md](ARCHITECTURE.md)** - Architecture overview and design decisions
- **[internal/usecase/README.md](internal/usecase/README.md)** - Solver algorithm details
- **[internal/infra/postgres/README.md](internal/infra/postgres/README.md)** - PostgreSQL integration
- **[internal/infra/redis/README.md](internal/infra/redis/README.md)** - Redis caching guide
- **[web/README.md](web/README.md)** - Web UI documentation

## 🛠️ Development

### Available Commands

```bash
# Show all commands
make help

# Build binary
make build

# Run tests
make test

# Format code
make fmt

# Vet code
make vet

# Lint code
make lint

# Update dependencies
make tidy

# Clean artifacts
make clean

# Format + vet + run
make dev
```

### Project Structure

```
.
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── domain/           # Business logic and models
│   ├── usecase/          # Use cases (DP solver)
│   ├── adapters/
│   │   └── http/         # HTTP handlers
│   └── infra/
│       ├── postgres/     # PostgreSQL client
│       ├── redis/        # Redis client
│       ├── config/       # Configuration
│       └── logger/       # Logging
├── web/                  # Static files (frontend)
├── configs/              # Configuration files
├── deployments/          # Docker, Kubernetes manifests
├── scripts/              # Utilities and scripts
└── build/                # Compiled binaries
```

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📧 Contact

For questions or feedback, please open an issue on GitHub.

---

**Built with ❤️ using Go, Clean Architecture, and Dynamic Programming**
