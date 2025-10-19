package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/evgenijurbanovskij/re-partners-assignment/internal/domain"
	"github.com/jmoiron/sqlx"
)

// Repository represents the PostgreSQL repository for pack_sets and calculations
type Repository struct {
	db *sqlx.DB
}

// NewRepository creates a new instance of the PostgreSQL repository
func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

// PackSet operations

// CreatePackSet creates a new pack size set
func (r *Repository) CreatePackSet(ctx context.Context, ps *domain.PackSizeSet) (*domain.PackSizeSet, error) {
	if err := ps.Validate(); err != nil {
		return nil, fmt.Errorf("invalid pack size set: %w", err)
	}

	model := FromPackSizeSet(ps)
	now := time.Now()
	model.CreatedAt = now
	model.UpdatedAt = now

	query := `
		INSERT INTO pack_sets (name, sizes, created_at, updated_at)
		VALUES (:name, :sizes, :created_at, :updated_at)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.GetContext(ctx, &model.ID, model)
	if err != nil {
		return nil, fmt.Errorf("failed to create pack set: %w", err)
	}

	return model.ToPackSizeSet(), nil
}

// GetPackSet получает набор размеров по ID
func (r *Repository) GetPackSet(ctx context.Context, id int64) (*domain.PackSizeSet, error) {
	query := `
		SELECT id, name, sizes, created_at, updated_at
		FROM pack_sets
		WHERE id = $1
	`

	var model PackSetModel
	err := r.db.GetContext(ctx, &model, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pack set not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pack set: %w", err)
	}

	return model.ToPackSizeSet(), nil
}

// GetPackSetByName получает набор размеров по имени
func (r *Repository) GetPackSetByName(ctx context.Context, name string) (*domain.PackSizeSet, error) {
	query := `
		SELECT id, name, sizes, created_at, updated_at
		FROM pack_sets
		WHERE name = $1
	`

	var model PackSetModel
	err := r.db.GetContext(ctx, &model, query, name)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pack set not found: %s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get pack set by name: %w", err)
	}

	return model.ToPackSizeSet(), nil
}

// ListPackSets получает список всех наборов размеров
func (r *Repository) ListPackSets(ctx context.Context, limit, offset int) ([]*domain.PackSizeSet, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, name, sizes, created_at, updated_at
		FROM pack_sets
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	var models []PackSetModel
	err := r.db.SelectContext(ctx, &models, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list pack sets: %w", err)
	}

	packSets := make([]*domain.PackSizeSet, 0, len(models))
	for i := range models {
		packSets = append(packSets, models[i].ToPackSizeSet())
	}

	return packSets, nil
}

// UpdatePackSet обновляет набор размеров
func (r *Repository) UpdatePackSet(ctx context.Context, ps *domain.PackSizeSet) error {
	if ps.ID == nil {
		return fmt.Errorf("pack set ID is required for update")
	}

	if err := ps.Validate(); err != nil {
		return fmt.Errorf("invalid pack size set: %w", err)
	}

	model := FromPackSizeSet(ps)
	model.UpdatedAt = time.Now()

	query := `
		UPDATE pack_sets
		SET name = :name, sizes = :sizes, updated_at = :updated_at
		WHERE id = :id
	`

	result, err := r.db.NamedExecContext(ctx, query, model)
	if err != nil {
		return fmt.Errorf("failed to update pack set: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pack set not found: %d", *ps.ID)
	}

	return nil
}

// DeletePackSet удаляет набор размеров
func (r *Repository) DeletePackSet(ctx context.Context, id int64) error {
	query := `DELETE FROM pack_sets WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete pack set: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("pack set not found: %d", id)
	}

	return nil
}

// Calculation операции

// SaveCalculation сохраняет результат расчёта
func (r *Repository) SaveCalculation(ctx context.Context, record *CalculationRecord) (int64, error) {
	if record.Solution == nil {
		return 0, fmt.Errorf("solution is required")
	}

	if err := record.Solution.Validate(); err != nil {
		return 0, fmt.Errorf("invalid solution: %w", err)
	}

	model := record.ToCalculationModel()
	model.CalculatedAt = time.Now()

	query := `
		INSERT INTO calculations (pack_set_id, pack_sizes, amount, breakdown, total_packs, overage, calculated_at)
		VALUES (:pack_set_id, :pack_sizes, :amount, :breakdown, :total_packs, :overage, :calculated_at)
		RETURNING id
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.GetContext(ctx, &id, model)
	if err != nil {
		return 0, fmt.Errorf("failed to save calculation: %w", err)
	}

	return id, nil
}

// GetCalculation получает расчёт по ID
func (r *Repository) GetCalculation(ctx context.Context, id int64) (*CalculationModel, error) {
	query := `
		SELECT id, pack_set_id, pack_sizes, amount, breakdown, total_packs, overage, calculated_at
		FROM calculations
		WHERE id = $1
	`

	var model CalculationModel
	err := r.db.GetContext(ctx, &model, query, id)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("calculation not found: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get calculation: %w", err)
	}

	return &model, nil
}

// ListCalculations получает список расчётов с фильтрацией
func (r *Repository) ListCalculations(ctx context.Context, packSetID *int64, limit, offset int) ([]*CalculationModel, error) {
	if limit <= 0 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	query := `
		SELECT id, pack_set_id, pack_sizes, amount, breakdown, total_packs, overage, calculated_at
		FROM calculations
	`

	var args []interface{}
	argIndex := 1

	if packSetID != nil {
		query += fmt.Sprintf(" WHERE pack_set_id = $%d", argIndex)
		args = append(args, *packSetID)
		argIndex++
	}

	query += fmt.Sprintf(" ORDER BY calculated_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	var models []*CalculationModel
	err := r.db.SelectContext(ctx, &models, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list calculations: %w", err)
	}

	return models, nil
}

// DeleteCalculation удаляет расчёт
func (r *Repository) DeleteCalculation(ctx context.Context, id int64) error {
	query := `DELETE FROM calculations WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete calculation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("calculation not found: %d", id)
	}

	return nil
}

// GetCalculationStats получает статистику по расчётам
func (r *Repository) GetCalculationStats(ctx context.Context) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_calculations,
			AVG(total_packs) as avg_packs,
			AVG(overage) as avg_overage,
			MIN(calculated_at) as first_calculation,
			MAX(calculated_at) as last_calculation
		FROM calculations
	`

	var statsModel struct {
		TotalCalculations int64           `db:"total_calculations"`
		AvgPacks          sql.NullFloat64 `db:"avg_packs"`
		AvgOverage        sql.NullFloat64 `db:"avg_overage"`
		FirstCalculation  sql.NullTime    `db:"first_calculation"`
		LastCalculation   sql.NullTime    `db:"last_calculation"`
	}

	err := r.db.GetContext(ctx, &statsModel, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get calculation stats: %w", err)
	}

	stats := map[string]interface{}{
		"total_calculations": statsModel.TotalCalculations,
	}

	if statsModel.AvgPacks.Valid {
		stats["avg_packs"] = statsModel.AvgPacks.Float64
	}
	if statsModel.AvgOverage.Valid {
		stats["avg_overage"] = statsModel.AvgOverage.Float64
	}
	if statsModel.FirstCalculation.Valid {
		stats["first_calculation"] = statsModel.FirstCalculation.Time
	}
	if statsModel.LastCalculation.Valid {
		stats["last_calculation"] = statsModel.LastCalculation.Time
	}

	return stats, nil
}
