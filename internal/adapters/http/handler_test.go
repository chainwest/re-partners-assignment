package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/evgenijurbanovskij/re-partners-assignment/internal/domain"
)

// Mock solver for tests
type mockSolver struct {
	solution *domain.Solution
	err      error
}

func (m *mockSolver) Solve(ctx context.Context, sizes []int, amount int) (*domain.Solution, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.solution, nil
}

// Mock logger for tests
type mockLogger struct{}

func (m *mockLogger) Info(ctx context.Context, msg string, fields map[string]interface{})  {}
func (m *mockLogger) Error(ctx context.Context, msg string, fields map[string]interface{}) {}
func (m *mockLogger) Debug(ctx context.Context, msg string, fields map[string]interface{}) {}
func (m *mockLogger) Warn(ctx context.Context, msg string, fields map[string]interface{})  {}

func TestPackHandler_SolvePacks_Success(t *testing.T) {
	// Arrange
	mockSol := &mockSolver{
		solution: &domain.Solution{
			Breakdown: map[int]int{250: 1, 500: 1},
			Packs:     2,
			Overage:   0,
			Amount:    750,
		},
	}
	handler := NewPackHandler(mockSol, &mockLogger{})

	reqBody := SolveRequest{
		Sizes:  []int{250, 500, 1000},
		Amount: 750,
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPost, "/packs/solve", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	handler.SolvePacks(w, req)

	// Assert
	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp SolveResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Packs != 2 {
		t.Errorf("expected 2 packs, got %d", resp.Packs)
	}
	if resp.Overage != 0 {
		t.Errorf("expected 0 overage, got %d", resp.Overage)
	}
}

func TestPackHandler_SolvePacks_ValidationError(t *testing.T) {
	tests := []struct {
		name       string
		reqBody    SolveRequest
		wantStatus int
	}{
		{
			name: "empty sizes",
			reqBody: SolveRequest{
				Sizes:  []int{},
				Amount: 100,
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "negative amount",
			reqBody: SolveRequest{
				Sizes:  []int{250, 500},
				Amount: -100,
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "zero amount",
			reqBody: SolveRequest{
				Sizes:  []int{250, 500},
				Amount: 0,
			},
			wantStatus: http.StatusUnprocessableEntity,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSol := &mockSolver{}
			handler := NewPackHandler(mockSol, &mockLogger{})

			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest(http.MethodPost, "/packs/solve", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.SolvePacks(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, w.Code)
			}
		})
	}
}

func TestPackHandler_SolvePacks_InvalidJSON(t *testing.T) {
	mockSol := &mockSolver{}
	handler := NewPackHandler(mockSol, &mockLogger{})

	req := httptest.NewRequest(http.MethodPost, "/packs/solve", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.SolvePacks(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestPackHandler_SolvePacks_MethodNotAllowed(t *testing.T) {
	mockSol := &mockSolver{}
	handler := NewPackHandler(mockSol, &mockLogger{})

	req := httptest.NewRequest(http.MethodGet, "/packs/solve", nil)
	w := httptest.NewRecorder()

	handler.SolvePacks(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status 405, got %d", w.Code)
	}
}
