package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/evgenijurbanovskij/re-partners-assignment/internal/domain"
)

func TestDPSolver_Solve(t *testing.T) {
	solver := NewDPSolver()
	ctx := context.Background()

	tests := []struct {
		name          string
		sizes         []int
		amount        int
		wantBreakdown map[int]int
		wantPacks     int
		wantOverage   int
		wantErr       bool
	}{
		// Examples from brief (required)
		{
			name:          "Brief example 1 - exact match 250",
			sizes:         []int{250, 500, 1000},
			amount:        250,
			wantBreakdown: map[int]int{250: 1},
			wantPacks:     1,
			wantOverage:   0,
			wantErr:       false,
		},
		{
			name:          "Brief example 2 - minimal overage 251",
			sizes:         []int{250, 500, 1000},
			amount:        251,
			wantBreakdown: map[int]int{500: 1},
			wantPacks:     1,
			wantOverage:   249,
			wantErr:       false,
		},
		{
			name:          "Brief example 3 - combination 1250",
			sizes:         []int{250, 500, 1000},
			amount:        1250,
			wantBreakdown: map[int]int{250: 1, 1000: 1},
			wantPacks:     2,
			wantOverage:   0,
			wantErr:       false,
		},
		{
			name:          "Brief example 4 - complex 12001",
			sizes:         []int{250, 500, 1000, 2000, 5000},
			amount:        12001,
			wantBreakdown: map[int]int{5000: 2, 2000: 1, 250: 1},
			wantPacks:     4,
			wantOverage:   249,
			wantErr:       false,
		},

		// Edge case from assignment (REQUIRED)
		{
			name:          "Edge case - 500k with small denominations",
			sizes:         []int{23, 31, 53},
			amount:        500000,
			wantBreakdown: map[int]int{23: 2, 31: 7, 53: 9429},
			wantPacks:     9438,
			wantOverage:   0,
			wantErr:       false,
		},

		// Critical edge cases
		{
			name:          "Minimal overage priority",
			sizes:         []int{3, 5},
			amount:        7,
			wantBreakdown: map[int]int{5: 1, 3: 1},
			wantPacks:     2,
			wantOverage:   1,
			wantErr:       false,
		},
		{
			name:          "Minimal packs priority",
			sizes:         []int{250, 500, 1000},
			amount:        1001,
			wantBreakdown: map[int]int{250: 1, 1000: 1},
			wantPacks:     2,
			wantOverage:   249,
			wantErr:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solution, err := solver.Solve(ctx, tt.sizes, tt.amount)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if solution == nil {
				t.Fatal("solution is nil")
			}

			// Check number of packs
			if solution.Packs != tt.wantPacks {
				t.Errorf("Packs = %d, want %d", solution.Packs, tt.wantPacks)
			}

			// Check overage
			if solution.Overage != tt.wantOverage {
				t.Errorf("Overage = %d, want %d", solution.Overage, tt.wantOverage)
			}

			// Check breakdown
			if !equalBreakdown(solution.Breakdown, tt.wantBreakdown) {
				t.Errorf("Breakdown = %v, want %v", solution.Breakdown, tt.wantBreakdown)
			}

			// Check solution validity
			if err := solution.Validate(); err != nil {
				t.Errorf("solution validation failed: %v", err)
			}

			// Check that solution covers required amount
			if solution.TotalItems() < tt.amount {
				t.Errorf("solution does not cover amount: %d < %d", solution.TotalItems(), tt.amount)
			}
		})
	}
}

func TestDPSolver_InvalidInput(t *testing.T) {
	solver := NewDPSolver()
	ctx := context.Background()

	tests := []struct {
		name   string
		sizes  []int
		amount int
	}{
		{
			name:   "Empty sizes",
			sizes:  []int{},
			amount: 100,
		},
		{
			name:   "Zero amount",
			sizes:  []int{250, 500},
			amount: 0,
		},
		{
			name:   "Negative amount",
			sizes:  []int{250, 500},
			amount: -100,
		},
		{
			name:   "Negative size",
			sizes:  []int{-250, 500},
			amount: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := solver.Solve(ctx, tt.sizes, tt.amount)
			if err == nil {
				t.Error("expected error, got nil")
			}
			if !domain.IsValidationError(err) {
				t.Errorf("expected validation error, got: %v", err)
			}
		})
	}
}

// Helper functions

func equalBreakdown(a, b map[int]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

// Benchmark tests

func BenchmarkDPSolver_SmallAmount(b *testing.B) {
	solver := NewDPSolver()
	ctx := context.Background()
	sizes := []int{250, 500, 1000, 2000, 5000}
	amount := 12001

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = solver.Solve(ctx, sizes, amount)
	}
}

func BenchmarkDPSolver_MediumAmount(b *testing.B) {
	solver := NewDPSolver()
	ctx := context.Background()
	sizes := []int{250, 500, 1000, 2000, 5000}
	amount := 100_000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = solver.Solve(ctx, sizes, amount)
	}
}

func BenchmarkDPSolver_LargeAmount(b *testing.B) {
	solver := NewDPSolver()
	ctx := context.Background()
	sizes := []int{1000, 2000, 5000, 10000}
	amount := 500_000

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = solver.Solve(ctx, sizes, amount)
	}
}

// BenchmarkSolveLarge tests performance for large W≈500k
func BenchmarkSolveLarge_EdgeCase(b *testing.B) {
	solver := NewDPSolver()
	ctx := context.Background()
	sizes := []int{23, 31, 53}
	amount := 500_000

	// Check that solution is found correctly
	solution, err := solver.Solve(ctx, sizes, amount)
	if err != nil {
		b.Fatalf("Failed to solve: %v", err)
	}
	if solution.TotalItems() < amount {
		b.Fatalf("Solution doesn't cover amount: %d < %d", solution.TotalItems(), amount)
	}

	b.ResetTimer()
	start := time.Now()
	for i := 0; i < b.N; i++ {
		_, _ = solver.Solve(ctx, sizes, amount)
	}
	elapsed := time.Since(start)

	// Check upper bound time for single operation
	avgTime := elapsed / time.Duration(b.N)
	if b.N == 1 && avgTime > 500*time.Millisecond {
		b.Logf("WARNING: Single operation took %v, expected ≤500ms", avgTime)
	}
}
