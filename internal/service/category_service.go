package service

import (
	"context"
	"pos-api/internal/domain"
	"pos-api/internal/repository"
	"strings"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(r repository.CategoryRepository) *CategoryService {
	return &CategoryService{repo: r}
}

func (s *CategoryService) Create(ctx context.Context, in domain.Category) (domain.Category, error) {
	created, err := s.repo.Create(ctx, domain.Category{
		Name:        strings.TrimSpace(in.Name),
		Description: in.Description,
	})
	if err != nil {
		return domain.Category{}, err
	}
	return created, nil
}

func (s *CategoryService) Get(ctx context.Context, id int) (domain.Category, error) {

	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.Category{}, err
	}
	return p, nil
}

func (s *CategoryService) List(ctx context.Context, limit, offset int) ([]domain.Category, error) {
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

func (s *CategoryService) Update(ctx context.Context, id int, in domain.Category) (domain.Category, error) {
	updated, err := s.repo.Update(ctx, id, domain.Category{
		ID:          id,
		Name:        strings.TrimSpace(in.Name),
		Description: in.Description,
	})
	if err != nil {
		return domain.Category{}, err
	}
	return updated, nil
}

func (s *CategoryService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
