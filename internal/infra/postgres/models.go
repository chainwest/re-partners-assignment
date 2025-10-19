package postgres

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/evgenijurbanovskij/re-partners-assignment/internal/domain"
)

// PackSetModel represents the pack size set model in the database
type PackSetModel struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Sizes     IntArray  `db:"sizes"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// CalculationModel represents the calculation model in the database
type CalculationModel struct {
	ID           int64        `db:"id"`
	PackSetID    *int64       `db:"pack_set_id"`
	PackSizes    IntArray     `db:"pack_sizes"`
	Amount       int          `db:"amount"`
	Breakdown    BreakdownMap `db:"breakdown"`
	TotalPacks   int          `db:"total_packs"`
	Overage      int          `db:"overage"`
	CalculatedAt time.Time    `db:"calculated_at"`
}

// IntArray represents an array of integers for JSONB
type IntArray []int

// Value implements driver.Valuer for IntArray
func (a IntArray) Value() (driver.Value, error) {
	if a == nil {
		return json.Marshal([]int{})
	}
	return json.Marshal(a)
}

// Scan implements sql.Scanner for IntArray
func (a *IntArray) Scan(value interface{}) error {
	if value == nil {
		*a = []int{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal IntArray value: %v", value)
	}

	var arr []int
	if err := json.Unmarshal(bytes, &arr); err != nil {
		return fmt.Errorf("failed to unmarshal IntArray: %w", err)
	}

	*a = arr
	return nil
}

// BreakdownMap represents map[int]int for JSONB
type BreakdownMap map[int]int

// Value implements driver.Valuer for BreakdownMap
func (m BreakdownMap) Value() (driver.Value, error) {
	if m == nil {
		return json.Marshal(map[int]int{})
	}
	return json.Marshal(m)
}

// Scan implements sql.Scanner for BreakdownMap
func (m *BreakdownMap) Scan(value interface{}) error {
	if value == nil {
		*m = make(map[int]int)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal BreakdownMap value: %v", value)
	}

	result := make(map[int]int)
	if err := json.Unmarshal(bytes, &result); err != nil {
		return fmt.Errorf("failed to unmarshal BreakdownMap: %w", err)
	}

	*m = result
	return nil
}

// ToPackSizeSet converts PackSetModel to domain.PackSizeSet
func (m *PackSetModel) ToPackSizeSet() *domain.PackSizeSet {
	id := m.ID
	name := m.Name
	return &domain.PackSizeSet{
		ID:    &id,
		Name:  &name,
		Sizes: m.Sizes,
	}
}

// FromPackSizeSet creates PackSetModel from domain.PackSizeSet
func FromPackSizeSet(ps *domain.PackSizeSet) *PackSetModel {
	model := &PackSetModel{
		Sizes: ps.Sizes,
	}

	if ps.ID != nil {
		model.ID = *ps.ID
	}

	if ps.Name != nil {
		model.Name = *ps.Name
	}

	return model
}

// CalculationRecord represents a calculation record for saving
type CalculationRecord struct {
	PackSetID *int64
	PackSizes []int
	Amount    int
	Solution  *domain.Solution
}

// ToCalculationModel converts CalculationRecord to CalculationModel
func (r *CalculationRecord) ToCalculationModel() *CalculationModel {
	return &CalculationModel{
		PackSetID:  r.PackSetID,
		PackSizes:  IntArray(r.PackSizes),
		Amount:     r.Amount,
		Breakdown:  BreakdownMap(r.Solution.Breakdown),
		TotalPacks: r.Solution.Packs,
		Overage:    r.Solution.Overage,
	}
}

// ToSolution converts CalculationModel to domain.Solution
func (m *CalculationModel) ToSolution() *domain.Solution {
	return &domain.Solution{
		Breakdown: map[int]int(m.Breakdown),
		Packs:     m.TotalPacks,
		Overage:   m.Overage,
		Amount:    m.Amount,
	}
}
