package domain

import "context"

// Solver defines the interface for solving the packing problem
// This interface represents a Port in Clean Architecture terms
type Solver interface {
	// Solve finds the optimal solution for packing a given amount of items
	// using available pack sizes.
	//
	// Parameters:
	//   - ctx: context for execution time management and operation cancellation
	//   - sizes: available pack sizes (must be unique, >0, ≤1e6)
	//   - amount: required amount of items to pack (must be >0)
	//
	// Returns:
	//   - Solution: optimal solution with minimum overage and number of packs
	//   - error: error if solution not found or input data is invalid
	//
	// Errors:
	//   - ErrInvalidInput: if input data fails validation
	//   - ErrNoSolutionStrict: if no exact solution found (without overage)
	//   - context.Canceled: if operation was canceled via context
	//   - context.DeadlineExceeded: if operation timeout exceeded
	Solve(ctx context.Context, sizes []int, amount int) (*Solution, error)
}

// PackSizeRepository defines the interface for working with pack size sets
// This interface represents a Port for the repository
type PackSizeRepository interface {
	// Create creates a new pack size set
	Create(ctx context.Context, packSizeSet *PackSizeSet) (*PackSizeSet, error)

	// GetByID gets a pack size set by identifier
	GetByID(ctx context.Context, id int64) (*PackSizeSet, error)

	// GetByName gets a pack size set by name
	GetByName(ctx context.Context, name string) (*PackSizeSet, error)

	// List returns a list of all pack size sets
	List(ctx context.Context) ([]*PackSizeSet, error)

	// Update updates an existing pack size set
	Update(ctx context.Context, packSizeSet *PackSizeSet) error

	// Delete deletes a pack size set by identifier
	Delete(ctx context.Context, id int64) error
}

// SolutionCache defines the interface for caching solutions
// This interface represents a Port for the cache
type SolutionCache interface {
	// Get retrieves a cached solution
	Get(ctx context.Context, key string) (*Solution, error)

	// Set saves a solution to cache
	Set(ctx context.Context, key string, solution *Solution) error

	// Delete removes a solution from cache
	Delete(ctx context.Context, key string) error

	// Clear clears the entire cache
	Clear(ctx context.Context) error
}

// SolverService defines the interface for the packing problem solving service
// with additional business logic (caching, logging, etc.)
type SolverService interface {
	// SolveWithCache solves the problem using cache
	SolveWithCache(ctx context.Context, sizes []int, amount int) (*Solution, error)

	// SolveWithPackSizeSet solves the problem using a saved pack size set
	SolveWithPackSizeSet(ctx context.Context, packSizeSetID int64, amount int) (*Solution, error)

	// ValidateInput validates input data before solving
	ValidateInput(sizes []int, amount int) error
}
