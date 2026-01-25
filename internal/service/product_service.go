package service

import (
	"context"
	"pos-api/internal/domain"
	"pos-api/internal/repository"
	"strings"
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(r repository.ProductRepository) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) Create(ctx context.Context, in domain.Product) (domain.Product, error) {
	created, err := s.repo.Create(ctx, domain.Product{
		Name:     strings.TrimSpace(in.Name),
		Price:    in.Price,
		Quantity: in.Quantity,
	})
	if err != nil {
		return domain.Product{}, err
	}
	return created, nil
}

func (s *ProductService) Get(ctx context.Context, id int) (domain.Product, error) {

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Product{}, err
	}
	return p, nil
}

func (s *ProductService) List(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	items, err := s.repo.List(ctx, repository.ListParams{Limit: limit, Offset: offset})
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *ProductService) Update(ctx context.Context, id int, in domain.Product) (domain.Product, error) {
	updated, err := s.repo.Update(ctx, id, domain.Product{
		ID:       id,
		Name:     strings.TrimSpace(in.Name),
		Price:    in.Price,
		Quantity: in.Quantity,
	})
	if err != nil {
		return domain.Product{}, err
	}
	return updated, nil
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
