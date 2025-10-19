package postgres

import (
	"context"
	"fmt"

	"github.com/evgenijurbanovskij/re-partners-assignment/internal/domain"
)

// RepositoryAdapter adapts the PostgreSQL repository for the HTTP handler
type RepositoryAdapter struct {
	repo *Repository
}

// NewRepositoryAdapter creates a new adapter
func NewRepositoryAdapter(repo *Repository) *RepositoryAdapter {
	return &RepositoryAdapter{repo: repo}
}

// SaveCalculation saves a calculation (implements interface for HTTP handler)
func (a *RepositoryAdapter) SaveCalculation(ctx context.Context, record interface{}) (int64, error) {
	// Convert generic record to typed structure
	recordMap, ok := record.(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid record type")
	}

	// Extract data
	packSizes, ok := recordMap["pack_sizes"].([]int)
	if !ok {
		return 0, fmt.Errorf("invalid pack_sizes type")
	}

	amount, ok := recordMap["amount"].(int)
	if !ok {
		return 0, fmt.Errorf("invalid amount type")
	}

	solution, ok := recordMap["solution"].(*domain.Solution)
	if !ok {
		return 0, fmt.Errorf("invalid solution type")
	}

	// Create record for saving
	calcRecord := &CalculationRecord{
		PackSetID: nil, // No link to pack_set yet
		PackSizes: packSizes,
		Amount:    amount,
		Solution:  solution,
	}

	// Save to database
	return a.repo.SaveCalculation(ctx, calcRecord)
}
