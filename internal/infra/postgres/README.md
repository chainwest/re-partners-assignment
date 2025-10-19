# PostgreSQL Repository

## Overview

PostgreSQL repository for storing pack size sets (`pack_sets`) and calculation history (`calculations`).

## ⚠️ Important: Core Independence from Database

**Key architectural principle**: The application core (`internal/usecase`) **DOES NOT DEPEND** on the database.

- ✅ **Usecase works autonomously** — all business logic is implemented without storage dependency
- ✅ **Database is optional** — application is fully functional without PostgreSQL
- ✅ **Storage for demonstration and audit** — Database is used only for:
  - Demonstrating integration capabilities
  - Saving calculation history for audit
  - Managing predefined pack size sets

### Architectural Advantages

1. **Testability** — usecase is tested without database
2. **Flexibility** — easy to replace PostgreSQL with another storage
3. **Performance** — calculations are performed in memory, without database calls
4. **Deployment simplicity** — can run without database setup

## Structure

```
internal/infra/postgres/
├── db.go              # Database connection
├── models.go          # Data models and converters
├── repository.go      # CRUD operations
├── integration_test.go # Integration tests
└── README.md          # This documentation

deployments/migrations/
├── 001_create_pack_sets.up.sql    # Create pack_sets table
├── 001_create_pack_sets.down.sql  # Rollback pack_sets migration
├── 002_create_calculations.up.sql # Create calculations table
└── 002_create_calculations.down.sql # Rollback calculations migration
```

## Database Schema

### Table `pack_sets`

Stores pack size sets:

```sql
CREATE TABLE pack_sets (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    sizes JSONB NOT NULL,  -- Array of sizes [250, 500, 1000, ...]
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes:**
- `idx_pack_sets_name` — fast search by name
- `idx_pack_sets_created_at` — sorting by creation date

### Table `calculations`

Stores calculation history for audit:

```sql
CREATE TABLE calculations (
    id BIGSERIAL PRIMARY KEY,
    pack_set_id BIGINT REFERENCES pack_sets(id) ON DELETE SET NULL,
    pack_sizes JSONB NOT NULL,      -- Sizes used in calculation
    amount INTEGER NOT NULL,         -- Required amount
    breakdown JSONB NOT NULL,        -- Breakdown: {size: count}
    total_packs INTEGER NOT NULL,    -- Total number of packs
    overage INTEGER NOT NULL,        -- Overage
    calculated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

**Indexes:**
- `idx_calculations_pack_set_id` — filtering by set
- `idx_calculations_calculated_at` — sorting by time
- `idx_calculations_amount` — analytics by amount
- `idx_calculations_pack_set_amount` — composite index for frequent queries

## Usage

### Database Connection

```go
import "github.com/re-partners/assignment/internal/infra/postgres"

// Configuration
cfg := postgres.Config{
    Host:            "localhost",
    Port:            "5432",
    User:            "postgres",
    Password:        "postgres",
    Database:        "re_partners",
    SSLMode:         "disable",
    MaxOpenConns:    25,
    MaxIdleConns:    25,
    ConnMaxLifetime: 5 * time.Minute,
}

// Connection
db, err := postgres.Connect(cfg)
if err != nil {
    log.Fatal(err)
}
defer postgres.Close(db)

// Creating repository
repo := postgres.NewRepository(db)
```

### Working with Pack Size Sets (Pack Sets)

```go
ctx := context.Background()

// Creating a set
packSet := &domain.PackSizeSet{
    Name:  strPtr("Standard"),
    Sizes: []int{250, 500, 1000, 2000, 5000},
}

created, err := repo.CreatePackSet(ctx, packSet)
if err != nil {
    log.Fatal(err)
}

// Get by ID
packSet, err = repo.GetPackSet(ctx, created.ID)

// Get by name
packSet, err = repo.GetPackSetByName(ctx, "Standard")

// List all sets
packSets, err := repo.ListPackSets(ctx, 10, 0) // limit, offset

// Update
packSet.Sizes = []int{250, 500, 1000, 2000, 5000, 10000}
err = repo.UpdatePackSet(ctx, packSet)

// Delete
err = repo.DeletePackSet(ctx, *packSet.ID)
```

### Saving calculations (Calculations)

```go
// Performing calculation (без БД!)
solver := usecase.NewPackSolver()
solution, err := solver.Solve([]int{250, 500, 1000}, 1263)

// Optional saving result for audit
record := &postgres.CalculationRecord{
    PackSetID: packSetID, // optional
    PackSizes: []int{250, 500, 1000},
    Amount:    1263,
    Solution:  solution,
}

calcID, err := repo.SaveCalculation(ctx, record)

// Getting calculation
calc, err := repo.GetCalculation(ctx, calcID)

// List of calculations
calculations, err := repo.ListCalculations(ctx, nil, 10, 0)

// List of calculations для конкретного набора
calculations, err = repo.ListCalculations(ctx, packSetID, 10, 0)

// Statistics
stats, err := repo.GetCalculationStats(ctx)
// stats: {
//   "total_calculations": 1000,
//   "avg_packs": 5.2,
//   "avg_overage": 12.5,
//   "first_calculation": "2025-01-01T00:00:00Z",
//   "last_calculation": "2025-10-19T12:00:00Z"
// }
```

## Migrations

### Applying Migrations

Migrations are located in `deployments/migrations/`:

```bash
# Using psql
psql -U postgres -d re_partners -f deployments/migrations/001_create_pack_sets.up.sql
psql -U postgres -d re_partners -f deployments/migrations/002_create_calculations.up.sql

# Rollback migrations
psql -U postgres -d re_partners -f deployments/migrations/002_create_calculations.down.sql
psql -U postgres -d re_partners -f deployments/migrations/001_create_pack_sets.down.sql
```

### Usage migrate CLI

```bash
# Installing migrate
brew install golang-migrate

# Applying migrations
migrate -path deployments/migrations -database "postgresql://postgres:postgres@localhost:5432/re_partners?sslmode=disable" up

# Rollback migrations
migrate -path deployments/migrations -database "postgresql://postgres:postgres@localhost:5432/re_partners?sslmode=disable" down
```

## API Integration

Repository connects to API as **optional component**:

```go
// В main.go
var repo *postgres.Repository

// Attempting database connection (optional)
if dbEnabled := os.Getenv("DB_ENABLED"); dbEnabled == "true" {
    db, err := postgres.Connect(cfg.Database)
    if err != nil {
        log.Printf("Warning: failed to connect to database: %v", err)
        log.Printf("Running without database (calculations will not be persisted)")
    } else {
        repo = postgres.NewRepository(db)
        defer postgres.Close(db)
        log.Println("Database connected successfully")
    }
}

// Handler works with or without repository
handler := http.NewHandler(solver, cache, repo) // repo can be nil
```

## Testing

### Integration Tests

Tests check:
- ✅ Применение и откат миграций
- ✅ CRUD операции для pack_sets
- ✅ CRUD операции для calculations
- ✅ Валидацию данных
- ✅ Индексы и производительность

```bash
# Running integration tests (требуется PostgreSQL)
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=postgres
export TEST_DB_NAME=re_partners_test

go test -v ./internal/infra/postgres/... -tags=integration
```

### Unit Tests

Models and converters are tested without database:

```bash
go test -v ./internal/infra/postgres/...
```

## Environment Variables

```bash
# Main settings
DB_ENABLED=true              # Enable/disable database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=re_partners
DB_SSLMODE=disable

# Connection pool settings
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_CONN_MAX_LIFETIME=5m
```

## Docker Compose

```yaml
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: re_partners
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./deployments/migrations:/docker-entrypoint-initdb.d

volumes:
  postgres_data:
```

## Performance

### Optimizations

1. **JSONB индексы** — fast search in JSON fields
2. **Составные индексы** — optimization of frequent queries
3. **Пул соединений** — connection reuse
4. **Batch операции** — bulk data insertion

### Monitoring

```sql
-- Table sizes
SELECT 
    schemaname,
    tablename,
    pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename)) AS size
FROM pg_tables
WHERE schemaname = 'public'
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;

-- Index usage
SELECT 
    schemaname,
    tablename,
    indexname,
    idx_scan,
    idx_tup_read,
    idx_tup_fetch
FROM pg_stat_user_indexes
ORDER BY idx_scan DESC;
```

## Usage Examples

### Preloading Standard Sets

```go
// Creating standard sets on application startup
standardSets := []struct {
    name  string
    sizes []int
}{
    {"Standard", []int{250, 500, 1000, 2000, 5000}},
    {"Small", []int{100, 250, 500}},
    {"Large", []int{1000, 2000, 5000, 10000}},
}

for _, set := range standardSets {
    ps := &domain.PackSizeSet{
        Name:  &set.name,
        Sizes: set.sizes,
    }
    _, err := repo.CreatePackSet(ctx, ps)
    if err != nil {
        log.Printf("Warning: failed to create pack set %s: %v", set.name, err)
    }
}
```

### Calculation audit

```go
// Getting recent calculations
recent, err := repo.ListCalculations(ctx, nil, 10, 0)
for _, calc := range recent {
    log.Printf("Calculation #%d: amount=%d, packs=%d, overage=%d",
        calc.ID, calc.Amount, calc.TotalPacks, calc.Overage)
}

// Statistics
stats, err := repo.GetCalculationStats(ctx)
log.Printf("Total calculations: %d", stats["total_calculations"])
log.Printf("Average packs: %.2f", stats["avg_packs"])
```

## Troubleshooting

### Problem: Cannot connect to database

```bash
# Checking availability PostgreSQL
psql -U postgres -h localhost -p 5432

# Checking environment variables
echo $DB_HOST $DB_PORT $DB_USER
```

### Problem: Migration error

```sql
-- Check migration status
SELECT * FROM schema_migrations;

-- Manual rollback
DROP TABLE IF EXISTS calculations;
DROP TABLE IF EXISTS pack_sets;
```

### Problem: Slow queries

```sql
-- Enabling slow query logging
ALTER DATABASE re_partners SET log_min_duration_statement = 100;

-- Query plan analysis
EXPLAIN ANALYZE SELECT * FROM calculations WHERE pack_set_id = 1;
```

## Additional Resources

- [PostgreSQL JSONB Documentation](https://www.postgresql.org/docs/current/datatype-json.html)
- [golang-migrate](https://github.com/golang-migrate/migrate)
- [lib/pq Driver](https://github.com/lib/pq)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
