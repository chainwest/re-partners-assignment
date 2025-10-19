package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"sync/atomic"
	"time"

	"github.com/evgenijurbanovskij/re-partners-assignment/internal/domain"
	"github.com/redis/go-redis/v9"
)

const (
	// DefaultTTL - default cache TTL (24 hours)
	DefaultTTL = 24 * time.Hour

	// CacheKeyPrefix - prefix for cache keys
	CacheKeyPrefix = "solver:"
)

// CachedSolver wraps Solver with Redis caching
type CachedSolver struct {
	solver domain.Solver
	client *redis.Client
	ttl    time.Duration

	// Metrics
	cacheHits   atomic.Uint64
	cacheMisses atomic.Uint64
}

// NewCachedSolver creates a new instance of CachedSolver
func NewCachedSolver(solver domain.Solver, client *redis.Client, ttl time.Duration) *CachedSolver {
	if ttl == 0 {
		ttl = DefaultTTL
	}

	return &CachedSolver{
		solver: solver,
		client: client,
		ttl:    ttl,
	}
}

// Solve implements the domain.Solver interface with caching
func (cs *CachedSolver) Solve(ctx context.Context, sizes []int, amount int) (*domain.Solution, error) {
	// Generate cache key
	cacheKey := cs.generateCacheKey(sizes, amount)

	// Try to get from cache
	solution, err := cs.getFromCache(ctx, cacheKey)
	if err == nil && solution != nil {
		// Cache hit
		cs.cacheHits.Add(1)
		return solution, nil
	}

	// Cache miss
	cs.cacheMisses.Add(1)

	// Call the original solver
	solution, err = cs.solver.Solve(ctx, sizes, amount)
	if err != nil {
		// Don't cache errors
		return nil, err
	}

	// Save to cache (asynchronously to not block the response)
	// Use a separate context with timeout for cache write
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := cs.saveToCache(cacheCtx, cacheKey, solution); err != nil {
			// Log error, but don't return it to the user
			// In production, this should use a proper logger
			_ = err
		}
	}()

	return solution, nil
}

// generateCacheKey generates a cache key: sha256(sorted sizes) + ":" + amount
func (cs *CachedSolver) generateCacheKey(sizes []int, amount int) string {
	// Copy and sort sizes for consistency
	sortedSizes := make([]int, len(sizes))
	copy(sortedSizes, sizes)
	sort.Ints(sortedSizes)

	// Create string representation of sorted sizes
	sizesStr := fmt.Sprintf("%v", sortedSizes)

	// Calculate SHA256 hash
	hash := sha256.Sum256([]byte(sizesStr))
	hashStr := hex.EncodeToString(hash[:])

	// Form key: prefix + hash + ":" + amount
	return fmt.Sprintf("%s%s:%d", CacheKeyPrefix, hashStr, amount)
}

// getFromCache retrieves a solution from cache
func (cs *CachedSolver) getFromCache(ctx context.Context, key string) (*domain.Solution, error) {
	data, err := cs.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("cache miss")
		}
		return nil, fmt.Errorf("redis get error: %w", err)
	}

	var solution domain.Solution
	if err := json.Unmarshal(data, &solution); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &solution, nil
}

// saveToCache saves a solution to cache
func (cs *CachedSolver) saveToCache(ctx context.Context, key string, solution *domain.Solution) error {
	data, err := json.Marshal(solution)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	if err := cs.client.Set(ctx, key, data, cs.ttl).Err(); err != nil {
		return fmt.Errorf("redis set error: %w", err)
	}

	return nil
}

// GetMetrics возвращает метрики кэша
func (cs *CachedSolver) GetMetrics() (hits, misses uint64) {
	return cs.cacheHits.Load(), cs.cacheMisses.Load()
}

// ResetMetrics сбрасывает метрики кэша
func (cs *CachedSolver) ResetMetrics() {
	cs.cacheHits.Store(0)
	cs.cacheMisses.Store(0)
}

// ClearCache очищает весь кэш решений
func (cs *CachedSolver) ClearCache(ctx context.Context) error {
	// Используем SCAN для поиска всех ключей с префиксом
	iter := cs.client.Scan(ctx, 0, CacheKeyPrefix+"*", 0).Iterator()

	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	if len(keys) > 0 {
		if err := cs.client.Del(ctx, keys...).Err(); err != nil {
			return fmt.Errorf("delete error: %w", err)
		}
	}

	return nil
}

// Ensure CachedSolver implements domain.Solver interface
var _ domain.Solver = (*CachedSolver)(nil)
