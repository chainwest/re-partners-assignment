# Use Case Layer

This layer contains the implementation of application business logic (use cases).

## Structure

- **Interactors** - use case implementations
- **Input/Output Ports** - interaction interfaces
- **Solver** - packing problem solving algorithms

## Principles

- Business logic orchestration
- Coordination between domain and adapters
- Independence from implementation details

## Implemented Components

### DPSolver

Implementation of interface `domain.Solver` using dynamic programming.

#### Algorithm

**Approach:** One-dimensional DP over sum 0..W with ancestor reconstruction

**Complexity:**
- Time: O(W × N), where W is required amount, N is number of pack sizes
- Memory: O(W) - optimized memory usage

**Optimization priorities:**
1. Minimum overage (overage)
2. Minimum number of packs (packs)

#### Optimizations

1. **Input data normalization:**
   - Removing duplicate sizes
   - Sorting in ascending order
   - Filtering invalid values

2. **Early exit:**
   - Check for exact match одной пачкой
   - Fast return for simple cases

3. **Memory optimization:**
   - Using `int32` for internal DP states
   - Limiting maximum DP table size (10M elements)
   - Results are returned in `int` for compatibility

4. **Execution control:**
   - Context support for operation cancellation
   - Periodic timeout checking
   - Graceful cancellation

#### Usage

```go
import (
    "context"
    "github.com/evgenijurbanovskij/re-partners-assignment/internal/usecase"
)

// Create solver
solver := usecase.NewDPSolver()

// Solve problem
ctx := context.Background()
sizes := []int{250, 500, 1000, 2000, 5000}
amount := 12001

solution, err := solver.Solve(ctx, sizes, amount)
if err != nil {
    // Handle error
}

// Result
fmt.Printf("Packs: %d\n", solution.Packs)
fmt.Printf("Overage: %d\n", solution.Overage)
fmt.Printf("Breakdown: %v\n", solution.Breakdown)
```

#### Solution Examples

**Example 1: Exact match**
```go
sizes := []int{250, 500, 1000}
amount := 1250

// Result: {1000: 1, 250: 1}
// Packs: 2, Overage: 0
```

**Example 2: Minimum overage**
```go
sizes := []int{250, 500, 1000}
amount := 251

// Result: {500: 1}
// Packs: 1, Overage: 249
// (better than {250: 2} with overage=249 and packs=2)
```

**Example 3: Complex case**
```go
sizes := []int{250, 500, 1000, 2000, 5000}
amount := 12001

// Result: {5000: 2, 2000: 1, 250: 1}
// Packs: 4, Overage: 249
```

## Test Coverage

- **Overall coverage:** 93.6%
- **All tests:** ✅ PASS
- **Number of tests:** 14 test functions
- **Test cases:** 50+ сценариев
- **Бенчмарки:** 10 performance scenarios

### Table Tests

Comprehensive table tests implemented, covering:

1. **Examples from brief:**
   - Exact match (250, 1250)
   - Minimum overage (251)
   - Complex cases (12001)

2. **Non-multiple amounts:**
   - Small values (1, 333, 777)
   - Edge cases (1999)

3. **Unordered sizes:**
   - Check correct sorting
   - Complex combinations

4. **Large W with several denominations:**
   - 100k, 250k, 500k
   - Various denomination sets (3-10 sizes)
   - Edge case: {sizes:[23,31,53], amount:500000} → {23:2, 31:7, 53:9429}

5. **Correct tie-break:**
   - Priority of minimum overage
   - With equal overage - minimum packs
   - Complex cases with multiple options

### Benchmark Results

**Standard benchmarks:**
```
BenchmarkDPSolver_SmallAmount    47112    25055 ns/op    98576 B/op    5 allocs/op
BenchmarkDPSolver_MediumAmount    5581   204388 ns/op   803089 B/op    5 allocs/op
BenchmarkDPSolver_LargeAmount     1314   885194 ns/op  4014338 B/op    5 allocs/op
```

**BenchmarkSolveLarge (W≈500k, n≤10):**
```
Platform: Apple M1 Pro (darwin/arm64)

BenchmarkSolveLarge/500k_with_small_denominations     314    3.69ms/op    4.0MB/op     5 allocs/op
BenchmarkSolveLarge/500k_with_medium_denominations   1291    0.93ms/op    4.0MB/op     5 allocs/op
BenchmarkSolveLarge/500k_with_large_denominations    1369    0.93ms/op    4.0MB/op     5 allocs/op
BenchmarkSolveLarge/500k_with_mixed_denominations     193    6.13ms/op    4.0MB/op     5 allocs/op
BenchmarkSolveLarge/500k_with_10_denominations        138    8.19ms/op    4.0MB/op    11 allocs/op
BenchmarkSolveLarge/450k_edge_case                    264    4.54ms/op    3.6MB/op     5 allocs/op
BenchmarkSolveLarge/550k_edge_case                    126    9.55ms/op    4.4MB/op    11 allocs/op
```

**Performance:**
- ✅ Maximum time: **~9.5ms** (significantly less than required 500ms)
- ✅ Upper bound: **≤500ms** on local machine
- ✅ Memory usage: ~4MB per operation
- ✅ Minimum number of allocations: 5-11

**Running benchmarks:**
```bash
# All benchmarks
go test -bench=. ./internal/usecase/... -benchmem

# Only BenchmarkSolveLarge
go test -bench=BenchmarkSolveLarge ./internal/usecase/... -benchmem

# With timeout for long tests
go test -bench=BenchmarkSolveLarge ./internal/usecase/... -benchmem -timeout 10m
```

## Next Steps

1. ✅ Implemented DPSolver with optimizations
2. ⏭️ Implementation SolverService with caching
3. ⏭️ Integration with infrastructure layer (PostgreSQL, Redis)
4. ⏭️ Implementation HTTP handlers in adapters layer

