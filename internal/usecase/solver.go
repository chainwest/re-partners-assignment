package usecase

import (
	"context"
	"sort"

	"github.com/evgenijurbanovskij/re-partners-assignment/internal/domain"
)

// DPSolver implements the domain.Solver interface using dynamic programming
// Algorithm: one-dimensional DP over sum 0..W with ancestor reconstruction
// Complexity: O(W * N) time, O(W) space
// Priority: minimize overage, then minimize number of packs
type DPSolver struct{}

// NewDPSolver creates a new instance of the DP solver
func NewDPSolver() *DPSolver {
	return &DPSolver{}
}

// dpState represents the DP state for a specific sum
// We use int32 to save memory where it's safe
type dpState struct {
	packs  int32 // Minimum number of packs to reach this sum
	parent int32 // Index of the pack size that led to this state (-1 if unreachable)
}

// Solve finds the optimal solution using dynamic programming
func (s *DPSolver) Solve(ctx context.Context, sizes []int, amount int) (*domain.Solution, error) {
	// Check context
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Validate input data
	if err := domain.ValidateSolverInput(sizes, amount); err != nil {
		return nil, err
	}

	// Normalize input sizes: remove duplicates and sort
	normalizedSizes := normalizeSizes(sizes)
	if len(normalizedSizes) == 0 {
		return nil, domain.NewSolverError(sizes, amount, "no valid sizes after normalization", domain.ErrInvalidInput)
	}

	// Early exit: check for exact match with a single pack
	for _, size := range normalizedSizes {
		if size == amount {
			breakdown := map[int]int{size: 1}
			return domain.NewSolution(breakdown, amount), nil
		}
	}

	// Determine the maximum sum for the DP table
	// We need to cover amount, but may have overage
	// Limit the search to a reasonable bound
	maxSum := calculateMaxSum(amount, normalizedSizes)

	// Initialize DP table
	// dp[i] = state for sum i
	dp := make([]dpState, maxSum+1)
	for i := range dp {
		dp[i] = dpState{packs: -1, parent: -1}
	}
	dp[0] = dpState{packs: 0, parent: -1}

	// Fill DP table
	bestSum := -1
	bestPacks := int32(-1)

	for sum := 0; sum <= maxSum; sum++ {
		// Check context periodically
		if sum%10000 == 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}
		}

		// If current sum is unreachable, skip
		if dp[sum].packs == -1 {
			continue
		}

		// Try adding each pack
		for idx, size := range normalizedSizes {
			newSum := sum + size
			if newSum > maxSum {
				continue
			}

			newPacks := dp[sum].packs + 1

			// Update state if this is the first reach or better by pack count
			if dp[newSum].packs == -1 || newPacks < dp[newSum].packs {
				dp[newSum].packs = newPacks
				dp[newSum].parent = int32(idx)
			}

			// If we reached or exceeded amount, update best solution
			if newSum >= amount {
				overage := newSum - amount

				// Update best solution by criteria:
				// 1. Less overage
				// 2. With equal overage - fewer packs
				if bestSum == -1 {
					bestSum = newSum
					bestPacks = newPacks
				} else {
					currentOverage := bestSum - amount
					if overage < currentOverage || (overage == currentOverage && newPacks < bestPacks) {
						bestSum = newSum
						bestPacks = newPacks
					}
				}

				// Early exit: if we found an exact solution with minimum pack count
				if overage == 0 {
					// Check if we can find better (fewer packs)
					// But for that we need to continue the search
					// However, if this is the first exact match, remember it
				}
			}
		}
	}

	// If no solution was found
	if bestSum == -1 {
		return nil, domain.NewSolverError(normalizedSizes, amount, "no solution found", domain.ErrNoSolution)
	}

	// Reconstruct solution
	breakdown := reconstructSolution(dp, normalizedSizes, bestSum)
	solution := domain.NewSolution(breakdown, amount)

	return solution, nil
}

// normalizeSizes removes duplicates and sorts sizes in ascending order
func normalizeSizes(sizes []int) []int {
	if len(sizes) == 0 {
		return nil
	}

	// Use map to remove duplicates
	uniqueSizes := make(map[int]bool)
	for _, size := range sizes {
		if size > 0 {
			uniqueSizes[size] = true
		}
	}

	// Convert to slice
	result := make([]int, 0, len(uniqueSizes))
	for size := range uniqueSizes {
		result = append(result, size)
	}

	// Sort in ascending order
	sort.Ints(result)

	return result
}

// calculateMaxSum calculates the maximum sum for the DP table
// Limit the search to a reasonable bound to avoid excessive memory usage
func calculateMaxSum(amount int, sizes []int) int {
	if len(sizes) == 0 {
		return amount
	}

	minSize := sizes[0] // sizes are sorted

	// Maximum overage = minimum size - 1
	// This guarantees we will find the optimal solution
	maxOverage := minSize - 1

	// Limit the maximum sum
	maxSum := amount + maxOverage

	// Additional check for reasonable memory limit
	const maxDPSize = 10_000_000 // 10M elements
	if maxSum > maxDPSize {
		maxSum = maxDPSize
	}

	return maxSum
}

// reconstructSolution reconstructs the solution from the DP table
func reconstructSolution(dp []dpState, sizes []int, targetSum int) map[int]int {
	breakdown := make(map[int]int)

	currentSum := targetSum
	for currentSum > 0 {
		if dp[currentSum].parent == -1 {
			break
		}

		sizeIdx := int(dp[currentSum].parent)
		size := sizes[sizeIdx]

		breakdown[size]++
		currentSum -= size
	}

	return breakdown
}

// Ensure DPSolver implements domain.Solver interface
var _ domain.Solver = (*DPSolver)(nil)
