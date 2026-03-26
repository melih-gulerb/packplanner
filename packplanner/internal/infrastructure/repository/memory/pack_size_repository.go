package memory

import (
	"context"
	"sync"

	"packplanner/internal/domain/pack"
)

// PackSizeRepository keeps pack sizes in memory for simple local usage and tests.
type PackSizeRepository struct {
	mu        sync.RWMutex
	packSizes []int
}

// NewPackSizeRepository creates an in-memory repository with normalized initial data.
func NewPackSizeRepository(initialPackSizes []int) (*PackSizeRepository, error) {
	normalized, err := pack.NormalizePackSizes(initialPackSizes)
	if err != nil {
		return nil, err
	}

	return &PackSizeRepository{
		packSizes: normalized,
	}, nil
}

func (r *PackSizeRepository) List(_ context.Context) ([]int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return a copy so callers cannot mutate repository state by accident.
	return append([]int(nil), r.packSizes...), nil
}

func (r *PackSizeRepository) Replace(_ context.Context, packSizes []int) error {
	normalized, err := pack.NormalizePackSizes(packSizes)
	if err != nil {
		return err
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Store a copy to keep ownership of the underlying slice inside the repository.
	r.packSizes = append([]int(nil), normalized...)
	return nil
}
