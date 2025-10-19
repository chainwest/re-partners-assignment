package domain

import (
	"errors"
	"testing"
)

func TestValidatePackSizes(t *testing.T) {
	tests := []struct {
		name    string
		sizes   []int
		wantErr bool
		errType error
	}{
		{
			name:    "valid sizes",
			sizes:   []int{250, 500, 1000, 2000, 5000},
			wantErr: false,
		},
		{
			name:    "empty sizes",
			sizes:   []int{},
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name:    "duplicate sizes",
			sizes:   []int{250, 500, 250},
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name:    "zero size",
			sizes:   []int{0, 500},
			wantErr: true,
			errType: ErrInvalidInput,
		},
		{
			name:    "negative size",
			sizes:   []int{-100, 500},
			wantErr: true,
			errType: ErrInvalidInput,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePackSizes(tt.sizes)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePackSizes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("ValidatePackSizes() error type = %v, want %v", err, tt.errType)
			}
		})
	}
}

func TestNewSolution(t *testing.T) {
	tests := []struct {
		name          string
		breakdown     map[int]int
		amount        int
		wantPacks     int
		wantOverage   int
		wantTotalItem int
	}{
		{
			name:          "exact solution",
			breakdown:     map[int]int{250: 1, 500: 1, 1000: 1},
			amount:        1750,
			wantPacks:     3,
			wantOverage:   0,
			wantTotalItem: 1750,
		},
		{
			name:          "solution with overage",
			breakdown:     map[int]int{500: 1},
			amount:        251,
			wantPacks:     1,
			wantOverage:   249,
			wantTotalItem: 500,
		},
		{
			name:          "empty solution",
			breakdown:     map[int]int{},
			amount:        100,
			wantPacks:     0,
			wantOverage:   0,
			wantTotalItem: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			solution := NewSolution(tt.breakdown, tt.amount)

			if solution.Packs != tt.wantPacks {
				t.Errorf("Packs = %d, want %d", solution.Packs, tt.wantPacks)
			}

			if solution.Overage != tt.wantOverage {
				t.Errorf("Overage = %d, want %d", solution.Overage, tt.wantOverage)
			}

			if solution.TotalItems() != tt.wantTotalItem {
				t.Errorf("TotalItems() = %d, want %d", solution.TotalItems(), tt.wantTotalItem)
			}
		})
	}
}

func TestSolutionIsValid(t *testing.T) {
	tests := []struct {
		name     string
		solution *Solution
		want     bool
	}{
		{
			name: "valid solution",
			solution: &Solution{
				Breakdown: map[int]int{250: 1, 500: 1},
				Packs:     2,
				Overage:   0,
				Amount:    750,
			},
			want: true,
		},
		{
			name: "invalid - nil breakdown",
			solution: &Solution{
				Breakdown: nil,
				Packs:     0,
				Overage:   0,
				Amount:    100,
			},
			want: false,
		},
		{
			name: "invalid - insufficient items",
			solution: &Solution{
				Breakdown: map[int]int{250: 1},
				Packs:     1,
				Overage:   0,
				Amount:    500,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.solution.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateAmount(t *testing.T) {
	tests := []struct {
		name    string
		amount  int
		wantErr bool
	}{
		{
			name:    "valid amount",
			amount:  1000,
			wantErr: false,
		},
		{
			name:    "minimum valid amount",
			amount:  1,
			wantErr: false,
		},
		{
			name:    "zero amount",
			amount:  0,
			wantErr: true,
		},
		{
			name:    "negative amount",
			amount:  -100,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateAmount(tt.amount)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAmount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCompareSolutions(t *testing.T) {
	tests := []struct {
		name string
		s1   *Solution
		s2   *Solution
		want *Solution
	}{
		{
			name: "s1 has less overage",
			s1:   &Solution{Overage: 10, Packs: 2},
			s2:   &Solution{Overage: 20, Packs: 1},
			want: &Solution{Overage: 10, Packs: 2},
		},
		{
			name: "same overage, s1 has fewer packs",
			s1:   &Solution{Overage: 10, Packs: 2},
			s2:   &Solution{Overage: 10, Packs: 3},
			want: &Solution{Overage: 10, Packs: 2},
		},
		{
			name: "s1 is nil",
			s1:   nil,
			s2:   &Solution{Overage: 10, Packs: 2},
			want: &Solution{Overage: 10, Packs: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CompareSolutions(tt.s1, tt.s2)
			if got.Overage != tt.want.Overage || got.Packs != tt.want.Packs {
				t.Errorf("CompareSolutions() = %v, want %v", got, tt.want)
			}
		})
	}
}
