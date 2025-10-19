# Domain Layer

This layer contains business logic and domain models of the project.

## Structure

### Entities (entities.go)

#### PackSizeSet
Represents a pack size set:
- `ID` - optional identifier
- `Name` - optional set name
- `Sizes` - array of pack sizes

#### Solution
Represents a packing problem solution:
- `Breakdown` - breakdown: pack size → quantity
- `Packs` - total number of packs
- `Overage` - overage (how much more than required)
- `Amount` - required amount of elements

### Validation Policies

#### ValidatePackSizes
Checks pack sizes against business rules:
- Sizes must be unique
- Sizes must be greater than 0
- Sizes must not exceed 1,000,000

#### ValidateAmount
Checks required amount:
- Amount must be greater than 0
- Amount must not exceed 1,000,000,000

### Port Interfaces (ports.go)

#### Solver
Main interface for solving packing problem:
```go
type Solver interface {
    Solve(ctx context.Context, sizes []int, amount int) (*Solution, error)
}
```

#### PackSizeRepository
Interface for working with pack size sets:
- `Create` - creating new set
- `GetByID` - get by identifier
- `GetByName` - get by name
- `List` - list all sets
- `Update` - update set
- `Delete` - delete set

#### SolutionCache
Interface for caching solutions:
- `Get` - get cached solution
- `Set` - save solution to cache
- `Delete` - delete solution from cache
- `Clear` - clear entire cache

#### SolverService
Interface for solver service with additional logic:
- `SolveWithCache` - solving with cache
- `SolveWithPackSizeSet` - solving with saved set
- `ValidateInput` - input data validation

### Domain Errors (errors.go)

#### Main Errors
- `ErrInvalidInput` - invalid input data
- `ErrNoSolutionStrict` - exact solution not found (without overage)
- `ErrNoSolution` - solution not found
- `ErrPackSizeSetNotFound` - pack size set not found
- `ErrPackSizeSetAlreadyExists` - set with this name already exists
- `ErrSolutionNotFound` - solution not found in cache
- `ErrCacheUnavailable` - cache unavailable

#### Specialized Errors
- `ValidationError` - validation error with context
- `SolverError` - error while solving problem

### Helper Functions

#### Working with Solutions
- `NewSolution` - creating new solution
- `EmptySolution` - creating empty solution
- `CompareSolutions` - comparing two solutions
- `IsSolutionStrict` - checking for exact solution

#### Working with pack size sets
- `NewPackSizeSet` - creating new set with validation

#### Error Type Checking
- `IsValidationError` - checking for validation error
- `IsNotFoundError` - checking for "not found" error
- `IsNoSolutionError` - checking for "solution not found" error

## Principles

1. **Independence from external dependencies** - layer does not depend on frameworks, databases or external libraries
2. **Pure business logic** - only domain rules and policies
3. **Domain-level validation** - all business rules are checked here
4. **Ports for external world** - interfaces for repositories, cache and services

## Testing

All components are covered by unit tests:
- `entities_test.go` - entity and validation tests
- `errors_test.go` - domain error tests

Running tests:
```bash
go test -v ./internal/domain/...
```

## Usage

### Creating a pack size set
```go
sizes := []int{250, 500, 1000, 2000, 5000}
id := int64(1)
name := "Standard Sizes"

packSizeSet, err := domain.NewPackSizeSet(sizes, &id, &name)
if err != nil {
    // Handling validation error
}
```

### Input Data Validation
```go
sizes := []int{250, 500, 1000}
amount := 1001

if err := domain.ValidateSolverInput(sizes, amount); err != nil {
    // Handling validation error
}
```

### Creating Solution
```go
breakdown := map[int]int{
    250: 1,
    500: 2,
}
solution := domain.NewSolution(breakdown, 1001)

// Validity check
if err := solution.Validate(); err != nil {
    // Solution is invalid
}
```

### Using Solver
```go
var solver domain.Solver // Implementation will be in usecase

solution, err := solver.Solve(ctx, sizes, amount)
if err != nil {
    if errors.Is(err, domain.ErrInvalidInput) {
        // Invalid input data
    } else if errors.Is(err, domain.ErrNoSolutionStrict) {
        // Exact solution not found
    }
}
```
