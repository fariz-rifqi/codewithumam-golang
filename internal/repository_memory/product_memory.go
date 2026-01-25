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

type ProductRepo struct {
	mu       sync.RWMutex
	nextID   int
	products map[int]domain.Product
}

func NewProductRepo() *ProductRepo {
	return &ProductRepo{
		nextID:   1,
		products: make(map[int]domain.Product),
	}
}

func (r *ProductRepo) Seed(items []domain.Product) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products = make(map[int]domain.Product, len(items))

	maxID := 0
	for _, p := range items {
		r.products[p.ID] = p
		if p.ID > maxID {
			maxID = p.ID
		}
	}
	r.nextID = maxID + 1
}
func (r *ProductRepo) Create(ctx context.Context, p domain.Product) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now().UTC()

	p.ID = r.nextID
	r.nextID++

	p.Name = strings.TrimSpace(p.Name)
	p.CreatedAt = now
	p.UpdatedAt = now

	r.products[p.ID] = p
	return p, nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id int) (domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.products[id]
	if !ok {
		return domain.Product{}, errors.New("not found")
	}
	return p, nil
}

func (r *ProductRepo) List(ctx context.Context, lp repository.ListParams) ([]domain.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]int, 0, len(r.products))
	for id := range r.products {
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
		return []domain.Product{}, nil
	}

	end := offset + limit
	if end > len(ids) {
		end = len(ids)
	}

	out := make([]domain.Product, 0, end-offset)
	for _, id := range ids[offset:end] {
		out = append(out, r.products[id])
	}
	return out, nil
}

func (r *ProductRepo) Update(ctx context.Context, id int, patch domain.Product) (domain.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, ok := r.products[id]
	if !ok {
		return domain.Product{}, errors.New("not found")
	}

	existing.Name = strings.TrimSpace(patch.Name)
	existing.Price = patch.Price
	existing.Quantity = patch.Quantity
	existing.UpdatedAt = time.Now().UTC()

	r.products[id] = existing
	return existing, nil
}

func (r *ProductRepo) Delete(ctx context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.products[id]; !ok {
		return errors.New("not found")
	}
	delete(r.products, id)
	return nil
}
