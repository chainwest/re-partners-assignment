# Architecture

Project follows **Clean Architecture** principles with clear layer separation.

## Layer Structure

```
HTTP Handler (adapters/http)
         ↓
   Use Case (usecase)
         ↓
     Domain (domain)
         ↑
Infrastructure (infra)
```

## Layers

### Domain (`/internal/domain`)
- Business logic and models
- No dependencies
- Repository interfaces

### Use Case (`/internal/usecase`)
- Application use cases
- Depends only on Domain
- **No DB dependency** - all calculations in memory

### Adapters (`/internal/adapters`)
- HTTP handlers and middleware
- Coordinates layer interactions

### Infrastructure (`/internal/infra`)
- PostgreSQL (optional, for audit)
- Redis (optional, for cache)
- Config, Logger

## Key Principles

1. **Use Case independence from infrastructure**
   - All calculations in memory
   - Easy to test without mocks

2. **Optional infrastructure**
   - Application works without DB and cache
   - PostgreSQL and Redis for demonstration

3. **Dependency Inversion**
   - Domain defines interfaces
   - Infrastructure implements interfaces

## Adding Features

1. Define entities in **Domain**
2. Create use case in **Use Case**
3. Implement handler in **Adapters**
4. Add infrastructure in **Infrastructure** (if needed)
