package domain

import (
	"errors"
	"fmt"
)

// PackSizeSet represents a set of pack sizes
type PackSizeSet struct {
	ID    *int64  // Optional identifier
	Name  *string // Optional set name
	Sizes []int   // Pack sizes
}

// Solution represents a packing problem solution
type Solution struct {
	Breakdown map[int]int // Breakdown: pack size -> quantity
	Packs     int         // Total number of packs
	Overage   int         // Overage (how much more than required)
	Amount    int         // Required amount
}

// NewPackSizeSet creates a new pack size set with validation
func NewPackSizeSet(sizes []int, id *int64, name *string) (*PackSizeSet, error) {
	if err := ValidatePackSizes(sizes); err != nil {
		return nil, err
	}

	return &PackSizeSet{
		ID:    id,
		Name:  name,
		Sizes: sizes,
	}, nil
}

// ValidatePackSizes checks pack size validation policy:
// - sizes must be unique
// - sizes must be greater than 0
// - sizes must not exceed 1e6
func ValidatePackSizes(sizes []int) error {
	if len(sizes) == 0 {
		return fmt.Errorf("%w: sizes cannot be empty", ErrInvalidInput)
	}

	const maxSize = 1_000_000

	seen := make(map[int]bool)
	for _, size := range sizes {
		// Check for positive value
		if size <= 0 {
			return fmt.Errorf("%w: size must be greater than 0, got %d", ErrInvalidInput, size)
		}

		// Check for maximum size
		if size > maxSize {
			return fmt.Errorf("%w: size must not exceed %d, got %d", ErrInvalidInput, maxSize, size)
		}

		// Check for uniqueness
		if seen[size] {
			return fmt.Errorf("%w: duplicate size %d", ErrInvalidInput, size)
		}
		seen[size] = true
	}

	return nil
}

// Validate checks the validity of the size set
func (p *PackSizeSet) Validate() error {
	return ValidatePackSizes(p.Sizes)
}

// NewSolution creates a new solution
func NewSolution(breakdown map[int]int, amount int) *Solution {
	totalPacks := 0
	totalItems := 0

	for size, count := range breakdown {
		totalPacks += count
		totalItems += size * count
	}

	overage := 0
	if totalItems > amount {
		overage = totalItems - amount
	}

	return &Solution{
		Breakdown: breakdown,
		Packs:     totalPacks,
		Overage:   overage,
		Amount:    amount,
	}
}

// IsValid checks if the solution is correct
func (s *Solution) IsValid() bool {
	if s.Breakdown == nil {
		return false
	}

	totalItems := 0
	for size, count := range s.Breakdown {
		if size <= 0 || count < 0 {
			return false
		}
		totalItems += size * count
	}

	// Solution must cover required amount
	return totalItems >= s.Amount
}

// TotalItems returns the total number of items in the solution
func (s *Solution) TotalItems() int {
	total := 0
	for size, count := range s.Breakdown {
		total += size * count
	}
	return total
}

// ValidateAmount checks the validity of the required amount
func ValidateAmount(amount int) error {
	if amount <= 0 {
		return fmt.Errorf("%w: amount must be greater than 0, got %d", ErrInvalidInput, amount)
	}

	const maxAmount = 1_000_000_000 // Reasonable maximum for amount
	if amount > maxAmount {
		return fmt.Errorf("%w: amount must not exceed %d, got %d", ErrInvalidInput, maxAmount, amount)
	}

	return nil
}

// ValidateSolverInput validates input data for the solver
func ValidateSolverInput(sizes []int, amount int) error {
	if err := ValidatePackSizes(sizes); err != nil {
		return err
	}

	if err := ValidateAmount(amount); err != nil {
		return err
	}

	return nil
}

// IsSolutionStrict checks if the solution is exact (without overage)
func IsSolutionStrict(solution *Solution) bool {
	return solution != nil && solution.Overage == 0
}

// CompareSolutions compares two solutions and returns the better one
// Criteria (by priority):
// 1. Less overage
// 2. Fewer packs
func CompareSolutions(s1, s2 *Solution) *Solution {
	if s1 == nil {
		return s2
	}
	if s2 == nil {
		return s1
	}

	// Priority 1: less overage
	if s1.Overage != s2.Overage {
		if s1.Overage < s2.Overage {
			return s1
		}
		return s2
	}

	// Priority 2: fewer packs
	if s1.Packs < s2.Packs {
		return s1
	}

	return s2
}

// Copy creates a copy of the solution
func (s *Solution) Copy() *Solution {
	if s == nil {
		return nil
	}

	breakdown := make(map[int]int, len(s.Breakdown))
	for k, v := range s.Breakdown {
		breakdown[k] = v
	}

	return &Solution{
		Breakdown: breakdown,
		Packs:     s.Packs,
		Overage:   s.Overage,
		Amount:    s.Amount,
	}
}

// EmptySolution creates an empty solution for the given amount
func EmptySolution(amount int) *Solution {
	return &Solution{
		Breakdown: make(map[int]int),
		Packs:     0,
		Overage:   0,
		Amount:    amount,
	}
}

// Validate checks the correctness of the solution
func (s *Solution) Validate() error {
	if s == nil {
		return errors.New("solution is nil")
	}

	if s.Breakdown == nil {
		return errors.New("breakdown is nil")
	}

	if s.Amount <= 0 {
		return fmt.Errorf("invalid amount: %d", s.Amount)
	}

	totalPacks := 0
	totalItems := 0

	for size, count := range s.Breakdown {
		if size <= 0 {
			return fmt.Errorf("invalid pack size: %d", size)
		}
		if count < 0 {
			return fmt.Errorf("invalid pack count: %d for size %d", count, size)
		}
		totalPacks += count
		totalItems += size * count
	}

	if totalPacks != s.Packs {
		return fmt.Errorf("packs mismatch: expected %d, got %d", s.Packs, totalPacks)
	}

	expectedOverage := 0
	if totalItems > s.Amount {
		expectedOverage = totalItems - s.Amount
	}

	if expectedOverage != s.Overage {
		return fmt.Errorf("overage mismatch: expected %d, got %d", s.Overage, expectedOverage)
	}

	if totalItems < s.Amount {
		return fmt.Errorf("solution does not cover required amount: %d < %d", totalItems, s.Amount)
	}

	return nil
}
