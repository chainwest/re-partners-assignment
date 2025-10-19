package usecase

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestDemoEdgeCase demonstrates the edge case solution from the assignment
func TestDemoEdgeCase(t *testing.T) {
	solver := NewDPSolver()
	ctx := context.Background()

	// Edge case from assignment
	sizes := []int{23, 31, 53}
	amount := 500_000

	fmt.Println("\n=== Edge case: {sizes:[23,31,53], amount:500000} ===")
	fmt.Printf("Input data:\n")
	fmt.Printf("  Pack sizes: %v\n", sizes)
	fmt.Printf("  Required amount: %d\n\n", amount)

	// Measure execution time
	start := time.Now()
	solution, err := solver.Solve(ctx, sizes, amount)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Solution error: %v", err)
	}

	fmt.Printf("Solution found in %v\n\n", elapsed)
	fmt.Printf("Result:\n")
	fmt.Printf("  Breakdown: %v\n", solution.Breakdown)
	fmt.Printf("  Number of packs: %d\n", solution.Packs)
	fmt.Printf("  Overage: %d\n", solution.Overage)
	fmt.Printf("  Total items: %d\n\n", solution.TotalItems())

	// Check correctness
	expectedBreakdown := map[int]int{23: 2, 31: 7, 53: 9429}
	expectedPacks := 9438
	expectedOverage := 0

	if !equalBreakdown(solution.Breakdown, expectedBreakdown) {
		t.Errorf("Breakdown = %v, expected %v", solution.Breakdown, expectedBreakdown)
	}

	if solution.Packs != expectedPacks {
		t.Errorf("Packs = %d, expected %d", solution.Packs, expectedPacks)
	}

	if solution.Overage != expectedOverage {
		t.Errorf("Overage = %d, expected %d", solution.Overage, expectedOverage)
	}

	// Check math
	total := 23*solution.Breakdown[23] + 31*solution.Breakdown[31] + 53*solution.Breakdown[53]
	fmt.Printf("Verification: 23×%d + 31×%d + 53×%d = %d\n",
		solution.Breakdown[23], solution.Breakdown[31], solution.Breakdown[53], total)

	if total != amount {
		t.Errorf("Sum doesn't match: %d ≠ %d", total, amount)
	}

	// Check performance
	if elapsed > 500*time.Millisecond {
		t.Errorf("Execution time %v exceeds required 500ms", elapsed)
	} else {
		fmt.Printf("\n✅ Performance: %v (requirement ≤500ms met)\n", elapsed)
	}

	fmt.Println("\n=== Test passed successfully ===")
}
