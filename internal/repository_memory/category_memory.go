package repository_memory

import (
	"context"
	"errors"
	"pos-api/internal/domain"
	"pos-api/internal/repository"
	"sort"
	"strings"
	"sync"
	"time"
)

type CategoryRepo struct {
	mu         sync.RWMutex
	nextID     int
	categories map[int]domain.Category
}

func NewCategoryRepo() *CategoryRepo {
	return &CategoryRepo{
		nextID:     1,
		categories: make(map[int]domain.Category),
	}
}

func (r *CategoryRepo) Seed(items []domain.Category) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.categories = make(map[int]domain.Category, len(items))

	maxID := 0
	for _, p := range items {
		r.categories[p.ID] = p
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	r.nextID = maxID + 1
}
func (r *CategoryRepo) Create(ctx context.Context, p domain.Category) (domain.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()

	p.ID = r.nextID
	r.nextID++

	p.Name = strings.TrimSpace(p.Name)
	p.CreatedAt = now
	p.UpdatedAt = now

	r.categories[p.ID] = p
	return p, nil
}

func (r *CategoryRepo) GetByID(ctx context.Context, id int) (domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.categories[id]
	if !ok {
		return domain.Category{}, errors.New("not found")
	}
	return p, nil
}

func (r *CategoryRepo) List(ctx context.Context, lp repository.ListParams) ([]domain.Category, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]int, 0, len(r.categories))
	for id := range r.categories {
		ids = append(ids, id)
	}

	sort.Slice(ids, func(i, j int) bool { return ids[i] > ids[j] })

	limit := lp.Limit
	offset := lp.Offset
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	if offset >= len(ids) {
		return []domain.Category{}, nil
	}

	end := offset + limit
	if end > len(ids) {
		end = len(ids)
	}

	out := make([]domain.Category, 0, end-offset)
	for _, id := range ids[offset:end] {
		out = append(out, r.categories[id])
	}
	return out, nil
}

func (r *CategoryRepo) Update(ctx context.Context, id int, patch domain.Category) (domain.Category, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.categories[id]
	if !ok {
		return domain.Category{}, errors.New("not found")
	}

	existing.Name = strings.TrimSpace(patch.Name)
	existing.Description = strings.TrimSpace(patch.Description)
	existing.UpdatedAt = time.Now().UTC()

	r.categories[id] = existing
	return existing, nil
}

func (r *CategoryRepo) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.categories[id]; !ok {
		return errors.New("not found")
	}
	delete(r.categories, id)
	return nil
}
