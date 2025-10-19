package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/evgenijurbanovskij/re-partners-assignment/internal/domain"
)

// SolveRequest represents a request to solve the packing problem
type SolveRequest struct {
	Sizes  []int `json:"sizes"`
	Amount int   `json:"amount"`
}

// SolveResponse represents a response with the packing solution
type SolveResponse struct {
	Solution map[int]int `json:"solution"` // size → count
	Overage  int         `json:"overage"`
	Packs    int         `json:"packs"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string                 `json:"error"`
	Message string                 `json:"message,omitempty"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Repository interface for optional storage
type Repository interface {
	SaveCalculation(ctx context.Context, record interface{}) (int64, error)
}

// PackHandler handles HTTP requests for solving the packing problem
type PackHandler struct {
	solver     domain.Solver
	logger     Logger
	repository Repository // Optional repository for audit
}

// NewPackHandler creates a new handler
func NewPackHandler(solver domain.Solver, logger Logger) *PackHandler {
	return &PackHandler{
		solver:     solver,
		logger:     logger,
		repository: nil, // No repository by default
	}
}

// WithRepository adds an optional repository
func (h *PackHandler) WithRepository(repo Repository) *PackHandler {
	h.repository = repo
	return h
}

// SolvePacks handles POST /packs/solve
func (h *PackHandler) SolvePacks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Check method (for compatibility with tests that call handler directly)
	if r.Method != http.MethodPost {
		h.respondError(w, r, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	// Check Content-Type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" && contentType != "" {
		h.respondError(w, r, http.StatusUnsupportedMediaType, "content type must be application/json", nil)
		return
	}

	// Decode request
	var req SolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, r, http.StatusBadRequest, "invalid JSON", map[string]interface{}{
			"parse_error": err.Error(),
		})
		return
	}

	// Validate request
	if err := h.validateRequest(&req); err != nil {
		var validationErr *domain.ValidationError
		if errors.As(err, &validationErr) {
			h.respondError(w, r, http.StatusUnprocessableEntity, "validation failed", map[string]interface{}{
				"field":   validationErr.Field,
				"value":   validationErr.Value,
				"message": validationErr.Message,
			})
			return
		}

		// General validation error
		if errors.Is(err, domain.ErrInvalidInput) {
			h.respondError(w, r, http.StatusUnprocessableEntity, err.Error(), nil)
			return
		}

		// Unexpected error
		h.logger.Error(ctx, "unexpected validation error", map[string]interface{}{
			"error": err.Error(),
		})
		h.respondError(w, r, http.StatusInternalServerError, "internal server error", nil)
		return
	}

	// Call solver
	solution, err := h.solver.Solve(ctx, req.Sizes, req.Amount)
	if err != nil {
		h.handleSolverError(w, r, err)
		return
	}

	// Optional save to DB for audit
	if h.repository != nil {
		// Create record for saving
		record := map[string]interface{}{
			"pack_sizes": req.Sizes,
			"amount":     req.Amount,
			"solution":   solution,
		}

		// Async save (don't block response)
		go func() {
			saveCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if _, err := h.repository.SaveCalculation(saveCtx, record); err != nil {
				h.logger.Error(saveCtx, "failed to save calculation", map[string]interface{}{
					"error": err.Error(),
				})
			}
		}()
	}

	// Build response
	response := SolveResponse{
		Solution: solution.Breakdown,
		Overage:  solution.Overage,
		Packs:    solution.Packs,
	}

	h.respondJSON(w, r, http.StatusOK, response)
}

// validateRequest validates the request
func (h *PackHandler) validateRequest(req *SolveRequest) error {
	// Validate sizes
	if len(req.Sizes) == 0 {
		return domain.NewValidationError("sizes", req.Sizes, "must not be empty")
	}

	// Validate through domain
	if err := domain.ValidateSolverInput(req.Sizes, req.Amount); err != nil {
		return err
	}

	return nil
}

// handleSolverError handles solver errors
func (h *PackHandler) handleSolverError(w http.ResponseWriter, r *http.Request, err error) {
	ctx := r.Context()

	// Validation errors
	if errors.Is(err, domain.ErrInvalidInput) {
		h.respondError(w, r, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	// No solution errors
	if errors.Is(err, domain.ErrNoSolution) || errors.Is(err, domain.ErrNoSolutionStrict) {
		h.respondError(w, r, http.StatusUnprocessableEntity, err.Error(), nil)
		return
	}

	// Context errors
	if errors.Is(err, context.Canceled) {
		h.respondError(w, r, http.StatusRequestTimeout, "request canceled", nil)
		return
	}

	if errors.Is(err, context.DeadlineExceeded) {
		h.respondError(w, r, http.StatusRequestTimeout, "request timeout", nil)
		return
	}

	// Unexpected error
	h.logger.Error(ctx, "solver error", map[string]interface{}{
		"error": err.Error(),
	})
	h.respondError(w, r, http.StatusInternalServerError, "internal server error", nil)
}

// respondJSON sends JSON response
func (h *PackHandler) respondJSON(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error(r.Context(), "failed to encode response", map[string]interface{}{
			"error": err.Error(),
		})
	}
}

// respondError sends error response
func (h *PackHandler) respondError(w http.ResponseWriter, r *http.Request, status int, message string, details map[string]interface{}) {
	response := ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
		Details: details,
	}

	h.respondJSON(w, r, status, response)
}
